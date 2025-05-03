package sbbs

import (
	"context"
	"fmt"
	"time"
)

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

// Changes the current working directory to the repositories root directory if
// the current working directory is inside a repo. Results in an error if the
// current working directory is not inside a repo.
func CdToRepoRoot() StageFunc {
	return Stage(
		"cd to repo root",
		func(ctxt context.Context, cmdLineArgs ...string) error {
			root, err := GitRevParse(ctxt)
			if err != nil {
				return err
			}

			return Cd(root)
		},
	)
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
