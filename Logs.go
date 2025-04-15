package sbbs

import (
	"fmt"
	"log"
	"os"
	"strings"
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
