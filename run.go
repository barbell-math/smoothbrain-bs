package sbbs

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Runs the program with the specified `args` using the supplied context. The
// supplied pipe will be used to capture Stdout. Stderr will always be printed
// to the console.
func RunCwd(
	ctxt context.Context,
	pipe io.Writer,
	cwd string,
	prog string,
	args ...string,
) error {
	var cmd *exec.Cmd
	cmd = exec.CommandContext(ctxt, prog, args...)
	cmd.Dir = cwd
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

// Runs the program with the specified `args` using the supplied context in the
// current working directory. The supplied pipe will be used to capture Stdout.
// Stderr will always be printed to the console.
func Run(
	ctxt context.Context,
	pipe io.Writer,
	prog string,
	args ...string,
) error {
	return RunCwd(ctxt, pipe, "", prog, args...)
}

// Runs the program with the specified `args` using the supplied context. All
// output of the program will be printed to stdout. Equivalent to calling [Run]
// and providing [os.Stdout] for the `pipe` argument.
func RunCwdStdout(
	ctxt context.Context,
	cwd string,
	prog string,
	args ...string,
) error {
	return RunCwd(ctxt, os.Stdout, cwd, prog, args...)
}

// Runs the program with the specified `args` using the supplied context in the
// current working directory. All output of the program will be printed to
// stdout. Equivalent to calling [Run] and providing [os.Stdout] for the `pipe`
// argument.
func RunStdout(ctxt context.Context, prog string, args ...string) error {
	return RunCwd(ctxt, os.Stdout, "", prog, args...)
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

// A helpful utility function that runs `git rev-parse --show-toplevel` and
// returns the stdout. This is often useful when attempting to change the
// current working directory to a repositories root directory.
func GitRevParse(ctxt context.Context) (string, error) {
	var buf bytes.Buffer
	err := Run(ctxt, &buf, "git", "rev-parse", "--show-toplevel")
	return strings.TrimSpace(buf.String()), err
}
