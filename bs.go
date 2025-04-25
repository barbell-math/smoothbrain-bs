// A very simple build system written in 100% golang to avoid the need to have
// cmake as a dependency.
package sbbs

import (
	"context"
	"errors"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
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
				// Note that the error was already printed out by the stage, it
				// does not need to be printed out here. It is meerly returned
				// to indicate if execution of the target should stop.
				LogPanic("An error was encountered, exiting.")
			}
		}
	}
}

func logUsage(progName string, availableTargets []string) {
	slices.Sort(availableTargets)
	LogInfo("Usage:")
	LogInfo("\t%s [target | -h | --help] [target specific args...]", progName)
	LogInfo("\tValid targets: %v", availableTargets)
}

// The main function that runs the build system. This is intended to be called
// by the `main` function of any code that uses this library.
func Main(progName string) {
	log.SetPrefix("smoothbrain-bs | ")
	availableTargets := slices.Collect(maps.Keys(targets))

	if len(os.Args) == 2 && slices.Contains([]string{"-h", "--help"}, strings.ToLower(os.Args[1])) {
		logUsage(progName, availableTargets)
		LogQuietInfo("Consider: Re-runing with a target")
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		LogErr("Expected target to be provided.")
		logUsage(progName, availableTargets)
		LogQuietInfo("Consider: Re-runing with a target")
		os.Exit(1)
	}

	if !slices.Contains(availableTargets, os.Args[1]) {
		LogErr("An invalid target was provided")
		logUsage(progName, availableTargets)
		LogQuietInfo("Consider: Re-runing with a valid target")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		targets[os.Args[1]](os.Args[2:]...)
	} else {
		targets[os.Args[1]]()
	}
}
