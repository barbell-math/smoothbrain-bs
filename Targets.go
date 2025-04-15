package sbbs

import (
	"bytes"
	"context"
	"strings"
)

// Registers a target that rebuilds the build system. This is often useful when
// changes are made to the build system of a project.
func RegisterBsBuildTarget() {
	RegisterTarget(
		context.Background(),
		"buildbs",
		Stage(
			"Run go build",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(ctxt, "go", "build", "-o", "./bs/bs", "./bs")
			},
		),
	)
}

// Registers a target that updates all dependences. Dependencies that are in
// the `barbell-math` repo will always be pinned at latest and all other
// dependencies will be updated to the latest version.
func RegisterUpdateDepsTarget() {
	RegisterTarget(
		context.Background(),
		"updateDeps",
		CdToRepoRoot(),
		Stage(
			"barbell-math package updates",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var packages bytes.Buffer
				if err := Run(
					ctxt, &packages, "go", "list", "-m", "-u", "all",
				); err != nil {
					return err
				}

				lines := strings.Split(packages.String(), "\n")
				// First line is the current package, skip it
				for i := 1; i < len(lines); i++ {
					iterPackage := strings.SplitN(lines[i], " ", 2)
					if !strings.Contains(iterPackage[0], "barbell-math") {
						continue
					}

					if err := RunStdout(
						ctxt, "go", "get", iterPackage[0]+"@latest",
					); err != nil {
						return err
					}
				}
				return nil
			},
		),
		Stage(
			"Non barbell-math package updates",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := RunStdout(ctxt, "go", "get", "-u", "./..."); err != nil {
					return err
				}
				if err := RunStdout(ctxt, "go", "mod", "tidy"); err != nil {
					return err
				}

				return nil
			},
		),
	)
}

// Registers two targets:
//  1. The first target will run gomarkdoc, embeding the results in README.md
//  2. The second target will install gomarkdoc using go intstall
func RegisterGoMarkDocTargets() {
	RegisterTarget(
		context.Background(),
		"gomarkdocInstall",
		CdToRepoRoot(),
		Stage(
			"Install gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(
					ctxt, "go",
					"install", "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest",
				)
			},
		),
	)

	RegisterTarget(
		context.Background(),
		"gomarkdocReadme",
		CdToRepoRoot(),
		Stage(
			"Run gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				err := RunStdout(
					ctxt, "gomarkdoc", "--embed", "--output", "README.md", ".",
				)
				if err != nil {
					LogQuietInfo("Consider running build system with gomarkdocInstall target if gomarkdoc is not installed")
				}
				return err
			},
		),
	)
}

// Registers three targets:
//  1. The first runs go fmt
//  2. The second runs go test without running any benchmarks
//  3. The third runs go test and runs all benchmarks
func RegisterCommonGoCmdTargets() {
	RegisterTarget(
		context.Background(),
		"fmt",
		CdToRepoRoot(),
		Stage(
			"Run go fmt",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(ctxt, "go", "fmt", "./...")
			},
		),
	)

	RegisterTarget(
		context.Background(),
		"test",
		CdToRepoRoot(),
		Stage(
			"Run go test",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(ctxt, "go", "test", "-v", "./...")
			},
		),
	)

	RegisterTarget(
		context.Background(),
		"bench",
		CdToRepoRoot(),
		Stage(
			"Run go test",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(ctxt, "go", "test", "-bench=.", "-v", "./...")
			},
		),
	)
}

// Defines all possible stages that can run in a mergegate target.
type MergegateTargets struct {
	// When true a stage will update all deps and run a diff to make sure that
	// the commited code is using all of the up to date dependencies.
	CheckDepsUpdated bool
	// When true a stage will install gomarkdoc, update the readme using the
	// `gomarkdocReadme` targer, and run a diff to make sure that the commited
	// readme is up to date.
	CheckReadmeGomarkdoc bool
	// When true a stage will run go fmt and then run a diff to make sure that
	// the commited code is properly formated.
	CheckFmt bool
	// When true a stage will run all unit tests in the repo to make sure that
	// the commited code passes all unit tests.
	CheckUnitTests bool
}

// Registers a mergegate target that will perform the actions that are defined
// by the [MergegateTargets] struct. See the [MergegateTargets] struct for
// details about the available stages the mergegate target can run.
func RegisterMergegateTarget(a MergegateTargets) {
	// Generate a stage that runs `git diff` and returns an error if there are any
	// differences.
	gitDiffStage := func(errMessage string, targetToRun string) StageFunc {
		return Stage(
			"Run Diff",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var buf bytes.Buffer
				if err := Run(ctxt, &buf, "git", "diff"); err != nil {
					return err
				}
				if buf.Len() > 0 {
					LogErr(errMessage)
					LogQuietInfo(buf.String())
					LogErr(
						"Run build system with %s and push any changes",
						targetToRun,
					)
					return StopErr
				}
				return nil
			},
		)
	}

	stages := []StageFunc{}
	if a.CheckDepsUpdated {
		stages = append(
			stages,
			TargetAsStage("updateDeps"),
			gitDiffStage("Out of date packages were detected", "updateDeps"),
		)
	}
	if a.CheckReadmeGomarkdoc {
		stages = append(
			stages,
			TargetAsStage("gomarkdocInstall"),
			TargetAsStage("gomarkdocReadme"),
			gitDiffStage("Readme is out of date", "gomarkdocReadme"),
		)
	}
	if a.CheckFmt {
		stages = append(
			stages,
			TargetAsStage("fmt"),
			gitDiffStage("Fix formatting to get a passing run!", "fmt"),
		)
	}
	if a.CheckUnitTests {
		stages = append(stages, TargetAsStage("unitTests"))
	}

	RegisterTarget(context.Background(), "mergegate", stages...)
}
