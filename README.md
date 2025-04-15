<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# sbbs

```go
import "github.com/barbell-math/smoothbrain-bs"
```

A very simple build system written in 100% golang to avoid the need to have cmake as a dependency.

## Index

- [Variables](<#variables>)
- [func Cd\(dir string\) error](<#Cd>)
- [func CreateFile\(name string\) \(\*os.File, error\)](<#CreateFile>)
- [func GitRevParse\(ctxt context.Context\) \(string, error\)](<#GitRevParse>)
- [func LogErr\(fmt string, args ...any\)](<#LogErr>)
- [func LogInfo\(fmt string, args ...any\)](<#LogInfo>)
- [func LogPanic\(fmt string, args ...any\)](<#LogPanic>)
- [func LogQuietInfo\(fmt string, args ...any\)](<#LogQuietInfo>)
- [func LogSuccess\(fmt string, args ...any\)](<#LogSuccess>)
- [func LogWarn\(fmt string, args ...any\)](<#LogWarn>)
- [func Main\(progName string\)](<#Main>)
- [func Mkdir\(path string\) error](<#Mkdir>)
- [func Open\(name string\) \(\*os.File, error\)](<#Open>)
- [func RegisterBsBuildTarget\(\)](<#RegisterBsBuildTarget>)
- [func RegisterCommonGoCmdTargets\(\)](<#RegisterCommonGoCmdTargets>)
- [func RegisterGoMarkDocTargets\(\)](<#RegisterGoMarkDocTargets>)
- [func RegisterMergegateTarget\(a MergegateTargets\)](<#RegisterMergegateTarget>)
- [func RegisterTarget\(ctxt context.Context, name string, stages ...StageFunc\)](<#RegisterTarget>)
- [func RegisterUpdateDepsTarget\(\)](<#RegisterUpdateDepsTarget>)
- [func Run\(ctxt context.Context, pipe io.Writer, prog string, args ...string\) error](<#Run>)
- [func RunCwd\(ctxt context.Context, pipe io.Writer, cwd string, prog string, args ...string\) error](<#RunCwd>)
- [func RunCwdStdout\(ctxt context.Context, cwd string, prog string, args ...string\) error](<#RunCwdStdout>)
- [func RunStdout\(ctxt context.Context, prog string, args ...string\) error](<#RunStdout>)
- [func RunTarget\(ctxt context.Context, target string, cmdLineArgs ...string\)](<#RunTarget>)
- [func TmpEnvVarSet\(name string, val string\) \(reset func\(\) error, err error\)](<#TmpEnvVarSet>)
- [func Touch\(name string\) error](<#Touch>)
- [type MergegateTargets](<#MergegateTargets>)
- [type StageFunc](<#StageFunc>)
  - [func CdToRepoRoot\(\) StageFunc](<#CdToRepoRoot>)
  - [func Stage\(name string, op func\(ctxt context.Context, cmdLineArgs ...string\) error\) StageFunc](<#Stage>)
  - [func TargetAsStage\(target string\) StageFunc](<#TargetAsStage>)
- [type TargetFunc](<#TargetFunc>)


## Variables

<a name="StopErr"></a>

```go
var (

    // An error that a stage can return to stop the target it is part of from
    // further execution. This is intended to be used when other error
    // information has been printed to the console.
    StopErr = errors.New("Generic stop error. See log above for error details.")
)
```

<a name="Cd"></a>
## func [Cd](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L37>)

```go
func Cd(dir string) error
```

A utility function that changes the programs current working directory and logs the old and new current working directories.

<a name="CreateFile"></a>
## func [CreateFile](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L8>)

```go
func CreateFile(name string) (*os.File, error)
```

A utility function that creates a file and logs the file's path.

<a name="GitRevParse"></a>
## func [GitRevParse](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L90>)

```go
func GitRevParse(ctxt context.Context) (string, error)
```

A helpful utility function that runs \`git rev\-parse \-\-show\-toplevel\` and returns the stdout. This is often useful when attempting to change the current working directory to a repositories root directory.

<a name="LogErr"></a>
## func [LogErr](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L65>)

```go
func LogErr(fmt string, args ...any)
```

Logs errors in red.

<a name="LogInfo"></a>
## func [LogInfo](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L45>)

```go
func LogInfo(fmt string, args ...any)
```

Logs info in cyan.

<a name="LogPanic"></a>
## func [LogPanic](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L70>)

```go
func LogPanic(fmt string, args ...any)
```

Logs errors in bold red and exits.

<a name="LogQuietInfo"></a>
## func [LogQuietInfo](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L50>)

```go
func LogQuietInfo(fmt string, args ...any)
```

Logs quiet info in gray.

<a name="LogSuccess"></a>
## func [LogSuccess](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L55>)

```go
func LogSuccess(fmt string, args ...any)
```

Logs successes in green.

<a name="LogWarn"></a>
## func [LogWarn](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Logs.go#L60>)

```go
func LogWarn(fmt string, args ...any)
```

Logs warnings in yellow.

<a name="Main"></a>
## func [Main](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L65>)

```go
func Main(progName string)
```

The main function that runs the build system. This is intended to be called by the \`main\` function of any code that uses this library.

<a name="Mkdir"></a>
## func [Mkdir](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L30>)

```go
func Mkdir(path string) error
```

A utility function that creates the supplied directory as well as all necessary parent directories.

<a name="Open"></a>
## func [Open](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L14>)

```go
func Open(name string) (*os.File, error)
```

A utility function that opens a file and logs the file's path.

<a name="RegisterBsBuildTarget"></a>
## func [RegisterBsBuildTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L11>)

```go
func RegisterBsBuildTarget()
```

Registers a target that rebuilds the build system. This is often useful when changes are made to the build system of a project.

<a name="RegisterCommonGoCmdTargets"></a>
## func [RegisterCommonGoCmdTargets](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L140>)

```go
func RegisterCommonGoCmdTargets()
```

Registers three targets:

1. The first runs go fmt
2. The second runs go test without running any benchmarks
3. The third runs go test and runs all benchmarks

<a name="RegisterGoMarkDocTargets"></a>
## func [RegisterGoMarkDocTargets](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L101>)

```go
func RegisterGoMarkDocTargets()
```

Registers two targets:

1. The first target will run gomarkdoc, embeding the results in README.md
2. The second target will install gomarkdoc using go intstall

<a name="RegisterMergegateTarget"></a>
## func [RegisterMergegateTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L198>)

```go
func RegisterMergegateTarget(a MergegateTargets)
```

Registers a mergegate target that will perform the actions that are defined by the [MergegateTargets](<#MergegateTargets>) struct. See the [MergegateTargets](<#MergegateTargets>) struct for details about the available stages the mergegate target can run.

<a name="RegisterTarget"></a>
## func [RegisterTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L43>)

```go
func RegisterTarget(ctxt context.Context, name string, stages ...StageFunc)
```

Registers a new build target to the build system. When run, the new target will sequentially run all provided stages, stopping if an error is encountered.

<a name="RegisterUpdateDepsTarget"></a>
## func [RegisterUpdateDepsTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L27>)

```go
func RegisterUpdateDepsTarget()
```

Registers a target that updates all dependences. Dependencies that are in the \`barbell\-math\` repo will always be pinned at latest and all other dependencies will be updated to the latest version.

<a name="Run"></a>
## func [Run](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L48-L53>)

```go
func Run(ctxt context.Context, pipe io.Writer, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context in the current working directory. The supplied pipe will be used to capture Stdout. Stderr will always be printed to the console.

<a name="RunCwd"></a>
## func [RunCwd](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L15-L21>)

```go
func RunCwd(ctxt context.Context, pipe io.Writer, cwd string, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context. The supplied pipe will be used to capture Stdout. Stderr will always be printed to the console.

<a name="RunCwdStdout"></a>
## func [RunCwdStdout](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L60-L65>)

```go
func RunCwdStdout(ctxt context.Context, cwd string, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context. All output of the program will be printed to stdout. Equivalent to calling [Run](<#Run>) and providing [os.Stdout](<https://pkg.go.dev/os/#Stdout>) for the \`pipe\` argument.

<a name="RunStdout"></a>
## func [RunStdout](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L73>)

```go
func RunStdout(ctxt context.Context, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context in the current working directory. All output of the program will be printed to stdout. Equivalent to calling [Run](<#Run>) and providing [os.Stdout](<https://pkg.go.dev/os/#Stdout>) for the \`pipe\` argument.

<a name="RunTarget"></a>
## func [RunTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Run.go#L80>)

```go
func RunTarget(ctxt context.Context, target string, cmdLineArgs ...string)
```

Runs the supplied target, given that the supplied target is present in the build systems target list. Execution of all further targets/stages will stop if running the supplied target fails.

<a name="TmpEnvVarSet"></a>
## func [TmpEnvVarSet](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L57>)

```go
func TmpEnvVarSet(name string, val string) (reset func() error, err error)
```

A utility function that changes the supplied env variable to the supplied value, returning a closure that can be used to set the env variable back to it's original value. If the supplied env variable did not exist before calling this function then the returned closure will remove the env variable instead of reseting it to it's original value.

<a name="Touch"></a>
## func [Touch](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Utility.go#L21>)

```go
func Touch(name string) error
```

A utility function that creates but does not open a file and logs the file's path.

<a name="MergegateTargets"></a>
## type [MergegateTargets](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Targets.go#L179-L193>)

Defines all possible stages that can run in a mergegate target.

```go
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
```

<a name="StageFunc"></a>
## type [StageFunc](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L25>)

The function that will be executed to perform an operation for a given target. The supplied context is meant to be used to control the runtime of the stage operation.

```go
type StageFunc func(ctxt context.Context, cmdLineArgs ...string) error
```

<a name="CdToRepoRoot"></a>
### func [CdToRepoRoot](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Stages.go#L46>)

```go
func CdToRepoRoot() StageFunc
```

Changes the current working directory to the repositories root directory if the current working directory is inside a repo. Results in an error if the current working directory is not inside a repo.

<a name="Stage"></a>
### func [Stage](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Stages.go#L14-L17>)

```go
func Stage(name string, op func(ctxt context.Context, cmdLineArgs ...string) error) StageFunc
```

Creates a stage that can be added to a build target. Stages define the operations that will take place when a build target is executing. The supplied context can be modified and passed to [Run](<#Run>) functions to deterministically control how long various operations take. This prevents builds from hanging forever.

<a name="TargetAsStage"></a>
### func [TargetAsStage](<https://github.com/barbell-math/smoothbrain-bs/blob/main/Stages.go#L63>)

```go
func TargetAsStage(target string) StageFunc
```

Runs the supplied target as though it were a stage, given that the supplied target is preset in the build systems target list. Execution of all further targets/stages will stop if running the supplied target fails.

<a name="TargetFunc"></a>
## type [TargetFunc](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L20>)

The function that will be executed when a target is run. This function will be given all of the leftover cmd line arguments that were supplied after the target. Parsing of these arguments is up to the logic defined be the targets stages.

```go
type TargetFunc func(cmdLineArgs ...string)
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

## Examples

For examples of using this build system refer to the following repositories:

1. [smoothbrain-errs](https://github.com/barbell-math/smoothbrain-errs/tree/main/bs/bs.go)
1. [smoothbrain-test](https://github.com/barbell-math/smoothbrain-test/tree/main/bs/bs.go)
1. [smoothbrain-hashmap](https://github.com/barbell-math/smoothbrain-hashmap/tree/main/bs)
1. [smoothbrain-argparse](https://github.com/barbell-math/smoothbrain-argparse/tree/main/bs/bs.go)

## Helpful Developer Cmds

To build the build system the first time:

```
go build -o ./bs/bs ./bs
```

The build system can then be used as usual:

```
./bs/bs --help
./bs/bs buildBs # builds the build system!
```
