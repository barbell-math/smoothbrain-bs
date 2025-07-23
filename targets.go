package sbbs

import (
	"bytes"
	"context"
	"fmt"
	"os/user"
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
					ctxt, "gomarkdoc",
					"-vv", "--embed",
					"--repository.default-branch", "main",
					"--output", "README.md", ".",
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

// Registers one target:
//  1. The first target will run install go-enum in ~/go/bin
func RegisterGoEnumTargets() {
	RegisterTarget(
		context.Background(),
		"goenumInstall",
		Stage(
			"Install go-enum",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var unameS bytes.Buffer
				var unameM bytes.Buffer

				if err := Run(ctxt, &unameS, "uname", "-s"); err != nil {
					return err
				}
				if err := Run(ctxt, &unameM, "uname", "-m"); err != nil {
					return err
				}

				usr, err := user.Current()
				if err != nil {
					return err
				}
				finalPath := path.Join(usr.HomeDir, "go", "bin", "go-enum")

				if err := RunStdout(
					ctxt,
					"curl", "-fsSL",
					fmt.Sprintf(
						"https://github.com/abice/go-enum/releases/download/v0.6.1/go-enum_%s_%s",
						strings.TrimSpace(unameS.String()),
						strings.TrimSpace(unameM.String()),
					),
					"-o", finalPath,
				); err != nil {
					return err
				}

				return RunStdout(ctxt, "chmod", "+x", finalPath)
			},
		),
	)
}

// Defines the available targets that can be added by
// [RegisterCommonGoCmdTargets].
type goTargets struct {
	TestTargetName     string
	TestArgs           []string
	BenchTargetName    string
	BenchArgs          []string
	FmtTargetName      string
	FmtArgs            []string
	GenerateTargetName string
	GenerateArgs       []string
}

func AllGoTargets() *goTargets {
	return (&goTargets{}).
		DefaultFmtTarget().
		DefaultGenerateTarget().
		DefaultTestTarget().
		DefaultBenchTarget()
}
func NewGoTargets() *goTargets {
	return &goTargets{}
}

const DefaultFmtTargetName = "fmt"
const DefaultGenerateTargetName = "generate"
const DefaultTestTargetName = "test"
const DefaultBenchTargetName = "bench"

func (g *goTargets) DefaultFmtTarget() *goTargets {
	g.FmtTargetName = DefaultFmtTargetName
	g.FmtArgs = []string{"./..."}
	return g
}
func (g *goTargets) SetFmtTarget(name string, args ...string) *goTargets {
	g.FmtTargetName = name
	g.FmtArgs = args
	return g
}

func (g *goTargets) DefaultGenerateTarget() *goTargets {
	g.GenerateTargetName = DefaultGenerateTargetName
	g.GenerateArgs = []string{"./..."}
	return g
}
func (g *goTargets) SetGenerateTarget(name string, args ...string) *goTargets {
	g.GenerateTargetName = name
	g.GenerateArgs = args
	return g
}

func (g *goTargets) DefaultTestTarget() *goTargets {
	g.TestTargetName = DefaultTestTargetName
	g.TestArgs = []string{"-v", "./..."}
	return g
}
func (g *goTargets) SetTestTarget(name string, args ...string) *goTargets {
	g.TestTargetName = name
	g.TestArgs = args
	return g
}

func (g *goTargets) DefaultBenchTarget() *goTargets {
	g.BenchTargetName = DefaultBenchTargetName
	g.BenchArgs = []string{"-bench=.", "-v", "./..."}
	return g
}
func (g *goTargets) SetBenchTarget(name string, args ...string) *goTargets {
	g.BenchTargetName = name
	g.BenchArgs = args
	return g
}

// Registers some common go cmds as targets. See the [MergegateTargets] struct
// for details about the available targets that can be added.
func RegisterCommonGoCmdTargets(g *goTargets) {
	if len(g.FmtArgs) > 0 && len(g.FmtTargetName) > 0 {
		args := []string{"fmt"}
		args = append(args, g.FmtArgs...)
		RegisterTarget(
			context.Background(),
			g.FmtTargetName,
			CdToRepoRoot(),
			Stage(
				"Run go fmt",
				func(ctxt context.Context, cmdLineArgs ...string) error {
					return RunStdout(ctxt, "go", args...)
				},
			),
		)
	}

	if len(g.GenerateArgs) > 0 && len(g.GenerateTargetName) > 0 {
		args := []string{"generate"}
		args = append(args, g.GenerateArgs...)
		RegisterTarget(
			context.Background(),
			g.GenerateTargetName,
			CdToRepoRoot(),
			Stage(
				"Run go generate",
				func(ctxt context.Context, cmdLineArgs ...string) error {
					return RunStdout(ctxt, "go", args...)
				},
			),
		)
	}

	if len(g.TestArgs) > 0 && len(g.TestTargetName) > 0 {
		args := []string{"test"}
		args = append(args, g.TestArgs...)
		RegisterTarget(
			context.Background(),
			g.TestTargetName,
			CdToRepoRoot(),
			Stage(
				"Run go test",
				func(ctxt context.Context, cmdLineArgs ...string) error {
					return RunStdout(ctxt, "go", args...)
				},
			),
		)
	}

	if len(g.BenchArgs) > 0 && len(g.BenchTargetName) > 0 {
		args := []string{"test"}
		args = append(args, g.BenchArgs...)
		RegisterTarget(
			context.Background(),
			g.BenchTargetName,
			CdToRepoRoot(),
			Stage(
				"Run go test",
				func(ctxt context.Context, cmdLineArgs ...string) error {
					return RunStdout(ctxt, "go", args...)
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
	// `gomarkdocReadme` target, and run a diff to make sure that the committed
	// readme is up to date.
	CheckReadmeGomarkdoc bool
	// When supplied, the given target will be expected to format the code. A
	// diff will then be run to make sure that the commited code is properly
	// formated.
	FmtTarget string
	// When supplied, the given target will be expected to test the code to make
	// sure the commited code passes all unit tests.
	TestTarget string
	// When supplied, the given target will be expected to generate the code
	// required for the project. A diff will then be run to make sure that the
	// commited code is properly formated.
	GenerateTarget string
	// Any stages that should be run prior to all other mergegate stages as
	// defined by the other flags in this struct. Useful for installing
	// dependencies that the other stages might rely upon.
	PreStages []StageFunc
	// Any stages that should be run after all other mergegate stages as defined
	// by the other flags in this struct. Useful for adding additional mergegate
	// checks.
	PostStages []StageFunc
}

// Registers a mergegate target that will perform the actions that are defined
// by the [MergegateTargets] struct. See the [MergegateTargets] struct for
// details about the available stages the mergegate target can run.
func RegisterMergegateTarget(a MergegateTargets) {
	// Generate a stage that runs `git diff` and returns an error if there are any
	// differences.
	stages := []StageFunc{}
	stages = append(stages, a.PreStages...)
	if len(a.FmtTarget) > 0 {
		stages = append(
			stages,
			TargetAsStage(a.FmtTarget),
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
	if len(a.GenerateTarget) > 0 {
		stages = append(
			stages,
			TargetAsStage(a.GenerateTarget),
			GitDiffStage("Out of sync generated code was detected", "generate"),
		)
	}
	if len(a.TestTarget) > 0 {
		stages = append(
			stages,
			TargetAsStage(a.TestTarget),
		)
	}
	stages = append(stages, a.PostStages...)

	RegisterTarget(context.Background(), "mergegate", stages...)
}
