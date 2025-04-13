// A very simple build system written in 100% golang.
package sbbs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

type (
	// The function that will be executed when a target is run. This function
	// will be given all of the leftover cmd line arguments that were supplied
	// after the target. Parsing of these arguments is up to the logic defined
	// be the targets stages.
	TargetFunc func(cmdLineArgs ...string)

	// The function that will be executed to perform an operation for a given
	// target. The supplied context is meant to be used to control the runtime
	// of the stage operation.
	StageFunc func(ctxt context.Context, cmdLineArgs ...string) error
)

const (
	// The identifier that will be printed on log lines that span for multiple
	// lines. The output will look like the following:
	//
	//	<log data> <log line 1>
	//	<log data>  |> <log line 2>
	//	<log data>  |> <log line 3>
	//	<log data>  ...
	multiLineIndent = " |> "

	// The color code to restore the consoles default colors.
	noColor = "\u001b[0m"
)

var (
	// The targets that are available to be called in the build system created
	// by the user. Targets are registered here through the [RegisterTarget]
	// function.
	targets = map[string]TargetFunc{}

	// An error that a stage can return to stop the target it is part of from
	// further execution. This is intended to be used when other error
	// information has been printed to the console.
	StopErr = errors.New("Generic stop error. See log above for error details.")
)

// Logs messages, splitting multi-line message into the following format:
//
//	<log data> <log line 1>
//	<log data>  |> <log line 2>
//	<log data>  |> <log line 3>
//	<log data>  ...
func multiLineLog(color string, fmtStr string, args ...any) {
	// This is a dumb hack to get arround any errors that look like the following:
	// bs/bs.go:20:17: non-constant format string in call to github.com/barbell-math/smoothbrain-bs.LogErr
	// See also: https://github.com/kubernetes/kubernetes/issues/127191
	_fmtStr := fmtStr

	str := fmt.Sprintf(_fmtStr, args...)
	lines := strings.Split(str, "\n")
	log.Printf(color + lines[0] + noColor)
	for i := 1; i < len(lines); i++ {
		log.Printf(multiLineIndent + color + lines[i] + noColor)
	}
}

// Logs info in cyan.
func LogInfo(fmt string, args ...any) {
	multiLineLog("\u001b[36m", fmt, args...)
}

// Logs quiet info in gray.
func LogQuietInfo(fmt string, args ...any) {
	multiLineLog("\u001b[90m", fmt, args...)
}

// Logs successes in green.
func LogSuccess(fmt string, args ...any) {
	multiLineLog("\u001b[32m", fmt, args...)
}

// Logs warnings in yellow.
func LogWarn(fmt string, args ...any) {
	multiLineLog("\u001b[33m", fmt, args...)
}

// Logs errors in red.
func LogErr(fmt string, args ...any) {
	multiLineLog("\u001b[31m", fmt, args...)
}

// Logs errors in bold red and exits.
func LogPanic(fmt string, args ...any) {
	multiLineLog("\u001b[1m\u001b[31m", fmt, args...)
	os.Exit(1)
}

// Runs the program with the specified `args` using the supplied context. The
// supplied pipe will be used to capture Stdout. Stderr will always be printed
// to the console.
func Run(ctxt context.Context, pipe io.Writer, prog string, args ...string) error {
	var cmd *exec.Cmd
	cmd = exec.CommandContext(ctxt, prog, args...)
	cmd.Stdout = pipe
	cmd.Stderr = os.Stderr

	LogQuietInfo("Running: '%s'", cmd.String())
	err := cmd.Run()
	if err != nil {
		return err
	}

	if cmd.ProcessState.ExitCode() != 0 {
		LogErr(
			"The process exited with a non-zero exit code: %d",
			cmd.ProcessState.ExitCode(),
		)
		return StopErr
	}

	return nil
}

// Runs the program with the specified `args` using the supplied context. All
// output of the program will be printed to stdout. Equivalent to calling [Run]
// and providing [os.Stdout] for the `pipe` argument.
func RunStdout(ctxt context.Context, prog string, args ...string) error {
	return Run(ctxt, os.Stdout, prog, args...)
}

// Runs the supplied target, given that the supplied target is present in the
// build systems target list. Execution of all further targets/stages will stop
// if running the supplied target fails.
func RunTarget(ctxt context.Context, target string, cmdLineArgs ...string) {
	if _, ok := targets[target]; !ok {
		LogPanic("Unrecognized target: %s", target)
	}
	targets[target](cmdLineArgs...)
}

// Creates a stage that can be added to a build target. Stages define the
// operations that will take place when a build target is executing. The
// supplied context can be modified and passed to [Run] functions to
// deterministically control how long various operations take. This prevents
// builds from hanging forever.
func Stage(
	name string,
	op func(ctxt context.Context, cmdLineArgs ...string) error,
) StageFunc {
	return func(ctxt context.Context, cmdLineArgs ...string) error {
		start := time.Now()
		LogInfo("Starting '%s' stage...", name)

		doneCh := make(chan error)
		go func() {
			doneCh <- op(ctxt, cmdLineArgs...)
		}()

		select {
		case err := <-doneCh:
			if err == nil {
				LogSuccess("Stage '%s': Completed Successfully", name)
			} else {
				LogErr("Stage '%s': Encountered an error: %s", name, err)
			}
			LogQuietInfo(multiLineIndent+"Time Delta: %s", time.Now().Sub(start))
			return err
		case <-ctxt.Done():
			LogErr("Stage '%s': Encountered an error: %s", name, ctxt.Err())
			return ctxt.Err()
		}
	}
}

// Runs the supplied target as though it were a stage, given that the supplied
// target is preset in the build systems target list. Execution of all further
// targets/stages will stop if running the supplied target fails.
func TargetAsStage(target string) StageFunc {
	return Stage(
		fmt.Sprintf("target:%s", target),
		func(ctxt context.Context, cmdLineArgs ...string) error {
			RunTarget(ctxt, target, cmdLineArgs...)
			return nil
		},
	)
}

// Registers a new build target to the build system. When run, the new target
// will sequentially run all provided stages, stopping if an error is
// encountered.
func RegisterTarget(ctxt context.Context, name string, stages ...StageFunc) {
	if _, ok := targets[name]; ok {
		LogPanic("Duplicate target name: %s", name)
	}
	targets[name] = func(cmdLineArgs ...string) {
		for i := range stages {
			if err := stages[i](ctxt, cmdLineArgs...); err != nil {
				LogPanic("An error was encountered, exiting.")
			}
		}
	}
}

func logUsage(progName string) {
	LogInfo("Usage:")
	LogInfo("\t./%s [target | -h | --help] [target specific args...]", progName)
	LogInfo("\tValid targets: %v", slices.Collect(maps.Keys(targets)))
}

// The main function that runs the build system. This is intended to be called
// by the `main` function of any code that uses this library.
func Main(progName string) {
	log.SetPrefix("smoothbrain-bs | ")
	if len(os.Args) == 2 && slices.Contains([]string{"-h", "--help"}, strings.ToLower(os.Args[1])) {
		logUsage(progName)
		LogQuietInfo("Consider: Re-runing with a target")
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		LogErr("Expected target to be provided.")
		logUsage(progName)
		LogQuietInfo("Consider: Re-runing with a target")
		os.Exit(1)
	}

	availableTargets := slices.Collect(maps.Keys(targets))
	if !slices.Contains(availableTargets, os.Args[1]) {
		LogErr("An invalid target was provided")
		logUsage(progName)
		LogQuietInfo("Consider: Re-runing with a valid target")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		targets[os.Args[1]](os.Args[2:]...)
	} else {
		targets[os.Args[1]]()
	}
}
