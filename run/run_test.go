package run

import (
	"bytes"
	"io"
	"testing"
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
)

type stdHandlesType struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func makeRunner(cmdName string, exit exec.OsStructExit, fatal LogFuncType) (*stdHandlesType, *Runner) {
	var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
	stdout, stderr, stdin := &stdoutBuf, &stderrBuf, &stdinBuf

	runner := Runner{
		CmdNameArgs: func(dcConfigReader dc_config.Reader) (string, []string) {
			return cmdName, make([]string, 0)
		},
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     &exec.ExecStruct{},
		OsStruct: &exec.OsStruct{
			Stdout: stdout,
			Stderr: stderr,
			Stdin:  stdin,
			Exit:   exit,
		},
		LogFunc: fatal,
	}

	stdHandles := stdHandlesType{stdout, stderr, stdin}

	return &stdHandles, &runner
}

func TestRunFatal(t *testing.T) {
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

	stdHandles, runner := makeRunner(`/\\nonexistent`, exit, fatal)

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

	cmd := runner.ExecStruct.GetCommand()
	if cmd == nil {
		t.Error(`Command is nil`)
	} else {
		stdoutActual, stdoutExpected := *cmd.Stdout, stdHandles.Stdout
		if stdoutActual != stdoutExpected {
			t.Errorf(`Stdout in 'cmd' '%p' was not replaced from 'os': '%p', cmp: '%v'`, stdoutActual, stdoutExpected, stdoutActual == stdoutExpected)
		}
		stderrActual, stderrExpected := *cmd.Stderr, stdHandles.Stderr
		if stderrActual != stderrExpected {
			t.Errorf(`Stderr in 'cmd' '%p' was not replaced from 'os': '%p', cmp: '%v'`, stderrActual, stderrExpected, stderrActual == stderrExpected)
		}
		stdinActual, stdinExpected := *cmd.Stdin, stdHandles.Stdin
		if stdinActual != stdinExpected {
			t.Errorf(`Stdin in 'cmd' '%p' was not replaced from 'os': '%p', cmp: '%v'`, stdinActual, stdinExpected, stdinActual == stdinExpected)
		}
	}
}
