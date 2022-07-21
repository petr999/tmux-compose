package run

import (
	"bytes"
	"testing"
	"tmux_compose/dc_config"
	"tmux_compose/runexec"
)

func TestRunFatal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var stdin bytes.Buffer

	timesLogFuncWasCalled := 0
	var logFuncArgs []any
	fatal := func(v ...any) {
		timesLogFuncWasCalled++
		logFuncArgs = v
	}
	timesExitWasCalled := 0
	var exitCode int
	exit := func(code int) {
		timesExitWasCalled++
		exitCode = code
	}

	runner := Runner{
		CmdNameArgs: func(dcConfigReader dc_config.Reader) (string, []string) {
			return `/\\nonexistent`, make([]string, 0)
		},
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     runexec.ExecStruct{},
		OsStruct: &runexec.OsStruct{
			Stdout: &stdout,
			Stderr: &stderr,
			Stdin:  &stdin,
			Exit:   exit,
		},
		LogFunc: fatal,
	}

	runner.Run()

	if timesLogFuncWasCalled != 1 {
		t.Errorf(`Log func was called %v times`, timesLogFuncWasCalled)
	}
	if len(logFuncArgs) != 1 {
		t.Errorf(`Wrong argument of log function: %v`, logFuncArgs)
	}
	if timesExitWasCalled != 1 {
		t.Errorf(`Exit func was called %v times`, timesExitWasCalled)
	}
	if exitCode != 1 {
		t.Errorf(`Wrong argument of Exit function: %v`, exitCode)
	}

}
