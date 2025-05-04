package sbbs

import (
	"bytes"
	"context"
	"fmt"
	"path"
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

				reset, err := TmpEnvVarSet("GOPROXY", "direct")
				if err != nil {
					return err
				}

				lines := strings.Split(packages.String(), "\n")
				// First line is the current package, skip it
				for i := 1; i < len(lines); i++ {
					iterPackage := strings.SplitN(lines[i], " ", 2)
					if !strings.Contains(iterPackage[0], "barbell-math") {
						continue
					}
					if strings.Count(iterPackage[0], "/") != 2 {
						continue
					}

					if err := RunStdout(
						ctxt, "go", "get", iterPackage[0]+"@latest",
					); err != nil {
						return err
					}
				}

				return reset()
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
		Stage(
			"Check if bs updated",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var buf bytes.Buffer
				if err := Run(
					ctxt, &buf, "git", "diff", "--unified=0", "go.mod",
				); err != nil {
					return err
				}
				if strings.Contains(
					strings.TrimSpace(buf.String()),
					"github.com/barbell-math/smoothbrain-bs",
				) {
					LogWarn("The build system package was upgraded!")
					LogWarn("It is recommended to rebuild your projects build system after this command completes.")
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

// Registers two targets:
//  1. The first target will run sqlc generate in the provided path, relative to
//     the repo root dir.
//  2. The second target will install sqlc using go intstall
func RegisterSqlcTargets(pathInRepo string) {
	RegisterTarget(
		context.Background(),
		"sqlc",
		Stage(
			fmt.Sprintf("Cd to %s", pathInRepo),
			func(ctxt context.Context, cmdLineArgs ...string) error {
				root, err := GitRevParse(ctxt)
				if err != nil {
					return err
				}
				finalPath := path.Join(root, pathInRepo)
				return Cd(finalPath)
			},
		),
		Stage(
			"Run sqlc generate",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				err := RunStdout(ctxt, "sqlc", "generate")
				if err != nil {
					LogQuietInfo("Consider running build system with sqlcInstall target if sqlc is not installed")
				}
				return err
			},
		),
	)
	RegisterTarget(
		context.Background(),
		"sqlcInstall",
		Stage(
			"Run sqlc install",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return RunStdout(
					ctxt, "go", "install",
					"github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
				)
			},
		),
	)
}

// Defines the available targets that can be added by
// [RegisterCommonGoCmdTargets].
type GoTargets struct {
	// When true a target will be added that runs `go test ./...`
	GenericTestTarget bool
	// When true a target will be added that runs `go test -bench=. ./...`
	GenericBenchTarget bool
	// When true a target will be added that runs `go fmt ./...`
	GenericFmtTarget bool
	// When true a target will be added that runs `go generate ./...`
	GenericGenerateTarget bool
}

// Registers some common go cmds as targets. See the [MergegateTargets] struct
// for details about the available targets that can be added.
func RegisterCommonGoCmdTargets(g GoTargets) {
	if g.GenericFmtTarget {
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
	}

	if g.GenericGenerateTarget {
		RegisterTarget(
			context.Background(),
			"generate",
			CdToRepoRoot(),
			Stage(
				"Run go generate",
				func(ctxt context.Context, cmdLineArgs ...string) error {
					return RunStdout(ctxt, "go", "generate", "./...")
				},
			),
		)
	}

	if g.GenericTestTarget {
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
	}

	if g.GenericBenchTarget {
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
	// When true a stage will run go generate and make sure that the generated
	// code matches what is commited to the repo.
	CheckGeneratedCode bool
}

// Registers a mergegate target that will perform the actions that are defined
// by the [MergegateTargets] struct. See the [MergegateTargets] struct for
// details about the available stages the mergegate target can run.
func RegisterMergegateTarget(a MergegateTargets) {
	// Generate a stage that runs `git diff` and returns an error if there are any
	// differences.
	stages := []StageFunc{}
	if a.CheckFmt {
		stages = append(
			stages,
			TargetAsStage("fmt"),
			GitDiffStage("Fix formatting to get a passing run!", "fmt"),
		)
	}
	if a.CheckReadmeGomarkdoc {
		stages = append(
			stages,
			TargetAsStage("gomarkdocInstall"),
			TargetAsStage("gomarkdocReadme"),
			GitDiffStage("Readme is out of date", "gomarkdocReadme"),
		)
	}
	if a.CheckDepsUpdated {
		stages = append(
			stages,
			TargetAsStage("updateDeps"),
			GitDiffStage("Out of date packages were detected", "updateDeps"),
		)
	}
	if a.CheckGeneratedCode {
		stages = append(
			stages,
			TargetAsStage("generate"),
			GitDiffStage("Out of sync generated code was detected", "generate"),
		)
	}
	if a.CheckUnitTests {
		stages = append(
			stages,
			TargetAsStage("test"),
		)
	}

	RegisterTarget(context.Background(), "mergegate", stages...)
}
