package run

import (
	"bytes"
	"testing"
	"tmux_compose/dc_config"
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

	runner := Runner{
		CmdNameArgs: func(dcConfigReader dc_config.Reader) (string, []string) {
			return `/\\nonexistent`, make([]string, 0)
		},
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     ExecStruct{},
		OsStruct: &OsStruct{
			Stdout: &stdout,
			Stderr: &stderr,
			Stdin:  &stdin,
		},
		LogFunc: fatal,
	}

	runner.Run()

	if timesLogFuncWasCalled != 1 {
		t.Errorf(`Log func was called %v times`, timesLogFuncWasCalled)
	}
	if logFuncArgs == nil {
		t.Errorf(`Wrong argument of log function: %s`, logFuncArgs)
	}
}
