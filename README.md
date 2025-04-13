<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# sbbs

```go
import "github.com/barbell-math/smoothbrain-bs"
```

A very simple build system written in 100% golang to avoid the need to have cmake as a dependency.

## Index

- [Variables](<#variables>)
- [func LogErr\(fmt string, args ...any\)](<#LogErr>)
- [func LogInfo\(fmt string, args ...any\)](<#LogInfo>)
- [func LogPanic\(fmt string, args ...any\)](<#LogPanic>)
- [func LogQuietInfo\(fmt string, args ...any\)](<#LogQuietInfo>)
- [func LogSuccess\(fmt string, args ...any\)](<#LogSuccess>)
- [func LogWarn\(fmt string, args ...any\)](<#LogWarn>)
- [func Main\(progName string\)](<#Main>)
- [func RegisterTarget\(ctxt context.Context, name string, stages ...StageFunc\)](<#RegisterTarget>)
- [func Run\(ctxt context.Context, pipe io.Writer, prog string, args ...string\) error](<#Run>)
- [func RunStdout\(ctxt context.Context, prog string, args ...string\) error](<#RunStdout>)
- [func RunTarget\(ctxt context.Context, target string, cmdLineArgs ...string\)](<#RunTarget>)
- [type StageFunc](<#StageFunc>)
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

<a name="LogErr"></a>
## func [LogErr](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L99>)

```go
func LogErr(fmt string, args ...any)
```

Logs errors in red.

<a name="LogInfo"></a>
## func [LogInfo](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L79>)

```go
func LogInfo(fmt string, args ...any)
```

Logs info in cyan.

<a name="LogPanic"></a>
## func [LogPanic](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L104>)

```go
func LogPanic(fmt string, args ...any)
```

Logs errors in bold red and exits.

<a name="LogQuietInfo"></a>
## func [LogQuietInfo](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L84>)

```go
func LogQuietInfo(fmt string, args ...any)
```

Logs quiet info in gray.

<a name="LogSuccess"></a>
## func [LogSuccess](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L89>)

```go
func LogSuccess(fmt string, args ...any)
```

Logs successes in green.

<a name="LogWarn"></a>
## func [LogWarn](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L94>)

```go
func LogWarn(fmt string, args ...any)
```

Logs warnings in yellow.

<a name="Main"></a>
## func [Main](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L223>)

```go
func Main(progName string)
```

The main function that runs the build system. This is intended to be called by the \`main\` function of any code that uses this library.

<a name="RegisterTarget"></a>
## func [RegisterTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L202>)

```go
func RegisterTarget(ctxt context.Context, name string, stages ...StageFunc)
```

Registers a new build target to the build system. When run, the new target will sequentially run all provided stages, stopping if an error is encountered.

<a name="Run"></a>
## func [Run](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L112>)

```go
func Run(ctxt context.Context, pipe io.Writer, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context. The supplied pipe will be used to capture Stdout. Stderr will always be printed to the console.

<a name="RunStdout"></a>
## func [RunStdout](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L138>)

```go
func RunStdout(ctxt context.Context, prog string, args ...string) error
```

Runs the program with the specified \`args\` using the supplied context. All output of the program will be printed to stdout. Equivalent to calling [Run](<#Run>) and providing [os.Stdout](<https://pkg.go.dev/os/#Stdout>) for the \`pipe\` argument.

<a name="RunTarget"></a>
## func [RunTarget](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L145>)

```go
func RunTarget(ctxt context.Context, target string, cmdLineArgs ...string)
```

Runs the supplied target, given that the supplied target is present in the build systems target list. Execution of all further targets/stages will stop if running the supplied target fails.

<a name="StageFunc"></a>
## type [StageFunc](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L29>)

The function that will be executed to perform an operation for a given target. The supplied context is meant to be used to control the runtime of the stage operation.

```go
type StageFunc func(ctxt context.Context, cmdLineArgs ...string) error
```

<a name="Stage"></a>
### func [Stage](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L157-L160>)

```go
func Stage(name string, op func(ctxt context.Context, cmdLineArgs ...string) error) StageFunc
```

Creates a stage that can be added to a build target. Stages define the operations that will take place when a build target is executing. The supplied context can be modified and passed to [Run](<#Run>) functions to deterministically control how long various operations take. This prevents builds from hanging forever.

<a name="TargetAsStage"></a>
### func [TargetAsStage](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L189>)

```go
func TargetAsStage(target string) StageFunc
```

Runs the supplied target as though it were a stage, given that the supplied target is preset in the build systems target list. Execution of all further targets/stages will stop if running the supplied target fails.

<a name="TargetFunc"></a>
## type [TargetFunc](<https://github.com/barbell-math/smoothbrain-bs/blob/main/bs.go#L24>)

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
1. [smoothbrain-hashmap](https://github.com/barbell-math/smoothbrain-hashmap/tree/main/bs/bs.go)
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
