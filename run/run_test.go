package run

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
)

const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`

type stdHandlesType struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func makeRunnerCommon() (*stdHandlesType, *Runner) {
	var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
	stdout, stderr, stdin := &stdoutBuf, &stderrBuf, &stdinBuf
	stdHandles := stdHandlesType{stdout, stderr, stdin}

	return &stdHandles, &Runner{
		DcConfigReader: dc_config.DcConfig{},
		ExecStruct:     &exec.ExecStruct{},
		OsStruct: &exec.OsStruct{
			Stdout: stdout,
			Stderr: stderr,
			Stdin:  stdin,
			Exit:   func(code int) {},
		},
		CmdNameArgs: func(dcConfigReader dc_config.Reader) (string, []string) {
			return ``, make([]string, 0)
		},
	}

}

func makeRunnerForFatal(cmdName string, exit exec.OsStructExit, fatal LogFuncType) (*stdHandlesType, *Runner) {
	stdHandles, runner := makeRunnerCommon()

	runner.CmdNameArgs = func(dcConfigReader dc_config.Reader) (string, []string) {
		return cmdName, make([]string, 0)
	}
	runner.OsStruct.Exit = exit
	runner.LogFunc = fatal

	return stdHandles, runner
}

type OsExecCmdRunDouble struct {
	timesOsExecCmdRunWasCalled uint8
	nameOsExecCommandWasCalled string
	argsOsExecCommandWasCalled []string
	Stdout, Stderr             io.Writer
	Stdin                      io.Reader
}

func (Obj *OsExecCmdRunDouble) Run() error {
	Obj.timesOsExecCmdRunWasCalled++
	return nil
}

type ExecStructDouble struct {
	osExecCmdRunDouble *OsExecCmdRunDouble
	cmd                *exec.CmdType
	Stdout, Stderr     io.Writer
	Stdin              io.Reader
}

func (execStructDouble *ExecStructDouble) MakeCommand(name string, arg ...string) {

	osExecCmdRunDouble := execStructDouble.osExecCmdRunDouble
	osExecCmdRunDouble.nameOsExecCommandWasCalled = name
	osExecCmdRunDouble.argsOsExecCommandWasCalled = arg

	cmd := &exec.CmdType{
		Obj:    osExecCmdRunDouble,
		Stdout: &execStructDouble.Stdout,
		Stderr: &execStructDouble.Stderr,
		Stdin:  &execStructDouble.Stdin,
	}

	execStructDouble.SetCommand(cmd)
}

func (execStructDouble *ExecStructDouble) GetCommand() *exec.CmdType {
	return execStructDouble.cmd
}

func (execStructDouble *ExecStructDouble) SetCommand(cmd *exec.CmdType) {
	execStructDouble.cmd = cmd
}

func makeRunnerForCmdRun(osExecCmdRunDouble *OsExecCmdRunDouble) *Runner {
	_, runner := makeRunnerCommon()

	runner.ExecStruct = &ExecStructDouble{
		osExecCmdRunDouble: osExecCmdRunDouble,
		Stdout:             osExecCmdRunDouble.Stdout,
		Stderr:             osExecCmdRunDouble.Stderr,
		Stdin:              osExecCmdRunDouble.Stdin,
	}

	return runner
}

func TestRunFatalAndStdHandles(t *testing.T) {
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

	stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, exit, fatal)

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

func TestCmdRunWasCalled(t *testing.T) {

	osExecCmdRunDouble := &OsExecCmdRunDouble{}
	runner := makeRunnerForCmdRun(osExecCmdRunDouble)

	runner.Run()

	if osExecCmdRunDouble.timesOsExecCmdRunWasCalled != 1 {
		t.Errorf(`os/exec.Command().Run() was called not once but '%v' times`, osExecCmdRunDouble.timesOsExecCmdRunWasCalled)
	}
	if osExecCmdRunDouble.nameOsExecCommandWasCalled != `` {
		t.Errorf(`os/exec.Run.Command() was called with name other than '%v'`, osExecCmdRunDouble.nameOsExecCommandWasCalled)
	}
	if !reflect.DeepEqual(osExecCmdRunDouble.argsOsExecCommandWasCalled, []string{}) {
		t.Errorf(`os/exec.Run.Command() was called with arg '%v' other than '%v'`, osExecCmdRunDouble.argsOsExecCommandWasCalled, []string{})
	}
}

func getLongArgs(amount int) []string {
	longArgs := make([]string, amount)
	loremArr := regexp.MustCompile(`\s+`).Split(loremString, amount)
	half := amount / 2
	iter := func(i int, sign bool) string {
		var idx int
		j := (i + 1) / 2
		if sign {
			idx = half + j
		} else {
			idx = half - j
		}
		return loremArr[idx]
	}

	for i := range longArgs {
		sign := false
		if i%2 != 0 {
			sign = true
		}
		longArgs[i] = iter(i, sign)
	}
	if len(longArgs) < amount {
		i := half + 1
		longArgs = append(longArgs, iter(i, true))
	}

	perms := rand.Perm(amount)
	for i, v := range perms {
		longArgs[i] = longArgs[i] + strings.Repeat(` `, v)
	}

	return longArgs
}

func TestCmdRunWasCalledWithArgs(t *testing.T) {

	const amount = 55

	longArgs := getLongArgs(amount)

	// name, args
	cmdsArgs := []map[string][]string{
		{`testingSchmesting`: {`--arg00`}},
		{`testing Schmesting`: {}},
		{``: {`--arg00`}},
		{``: {`--vaudeville 577`}},
		{`testingSchmesting`: {`--arg00`, `abcd`, `vaudeville 577`, `35`, `-88.55000`}},
		{`docker-compose`: longArgs},
	}

	for _, cmdArgs := range cmdsArgs {
		for name, args := range cmdArgs {
			osExecCmdRunDouble := &OsExecCmdRunDouble{}
			runner := makeRunnerForCmdRun(osExecCmdRunDouble)
			runner.CmdNameArgs = func(dcConfigReader dc_config.Reader) (string, []string) {
				return name, args
			}

			runner.Run()

			if osExecCmdRunDouble.timesOsExecCmdRunWasCalled != 1 {
				t.Errorf(`os/exec.Command().Run() was called not once but '%v' times`, osExecCmdRunDouble.timesOsExecCmdRunWasCalled)
			}
			if osExecCmdRunDouble.nameOsExecCommandWasCalled != name {
				t.Errorf(`os/exec.Run.Command() was called with name other than '%v'`, osExecCmdRunDouble.nameOsExecCommandWasCalled)
			}
			if !reflect.DeepEqual(osExecCmdRunDouble.argsOsExecCommandWasCalled, args) {
				t.Errorf(`os/exec.Run.Command() was called with arg '%v' other than '%v'`, osExecCmdRunDouble.argsOsExecCommandWasCalled, []string{})
			}
		}
	}
}

// func TestStdoutByConfig(t *testing.T) {
// 	osStruct := exec.OsStruct{
// 		Getenv: func(key string) string {
// 			if key == DryRunEnvVarName {
// 				return `1`
// 			}
// 			return ``
// 		},
// 	}
// 	runner := makeRunner(Runner{OsStruct: osStruct})
// }
