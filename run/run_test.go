package run

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"tmux_compose/dc_config"
	"tmux_compose/run/exec"
)

const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`

type stdHandlesType struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Stdin  *bytes.Buffer
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
			Getenv: func(string) string { return `` },
			Chdir:  func(string) error { return nil },
		},
		CmdNameArgs: func(dcConfigReader dc_config.Reader) (string, []string) {
			return ``, []string{}
		},
	}

}

func makeRunner(tle *testLogfuncExitType) (stdHandles *stdHandlesType, runner *Runner) {
	stdHandles, runner = makeRunnerCommon()

	runner.OsStruct.Exit = tle.exit
	runner.LogFunc = tle.fatal

	return stdHandles, runner
}

func makeRunnerForFatal(cmdName string, tle *testLogfuncExitType) (*stdHandlesType, *Runner) {
	stdHandles, runner := makeRunner(tle)

	runner.CmdNameArgs = func(dcConfigReader dc_config.Reader) (string, []string) {
		return cmdName, make([]string, 0)
	}

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

func (execStructDouble *ExecStructDouble) MakeCommand(dryRun *exec.MakeCommandDryRunType,
	nameArgs exec.NameArgsType) {
	// osStruct *exec.OsStruct, name string, arg ...string) {

	osExecCmdRunDouble := execStructDouble.osExecCmdRunDouble
	osExecCmdRunDouble.nameOsExecCommandWasCalled = nameArgs.Name
	osExecCmdRunDouble.argsOsExecCommandWasCalled = nameArgs.Args

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

func makeRunnerForCmdRun(osExecCmdRunDouble *OsExecCmdRunDouble, tle *testLogfuncExitType) *Runner {
	_, runner := makeRunner(tle)

	runner.ExecStruct = &ExecStructDouble{
		osExecCmdRunDouble: osExecCmdRunDouble,
		Stdout:             osExecCmdRunDouble.Stdout,
		Stderr:             osExecCmdRunDouble.Stderr,
		Stdin:              osExecCmdRunDouble.Stdin,
	}

	return runner
}

type testLogfuncExitType struct {
	timesLogFuncWasCalled *int
	logFuncArgs           *[]any
	fatal                 func(s string)
	timesExitWasCalled    *int
	exitCode              *int
	exit                  func(code int)
}

func (tle *testLogfuncExitType) LogfuncAndExitTestWascalledsAndArgs(t *testing.T, tlfwcExpected int, lfaExpected []any, tewcExpected int, exitCodeExpected int) {

	if *tle.timesLogFuncWasCalled != tlfwcExpected {
		t.Errorf(`Log func was called %v times`, *tle.timesLogFuncWasCalled)
	}
	if len(*tle.logFuncArgs) != len(lfaExpected) {
		t.Errorf(`Wrong argument of log function: %v`, *tle.logFuncArgs)
	}
	if *tle.timesExitWasCalled != tewcExpected {
		t.Errorf(`Exit func was called %v times`, *tle.timesExitWasCalled)
	}
	if *tle.exitCode != exitCodeExpected {
		t.Errorf(`Wrong argument of Exit function: %v`, *tle.exitCode)
	}

}

func getTestLogfuncExitType() testLogfuncExitType {
	timesLogFuncWasCalled := 0
	var logFuncArgs []any
	var fatal LogFuncType = func(s string) {
		timesLogFuncWasCalled++
		logFuncArgs = append(logFuncArgs, s)
	}

	timesExitWasCalled := 0
	var exitCode int
	exit := func(code int) {
		timesExitWasCalled++
		exitCode = code
	}
	return testLogfuncExitType{
		timesLogFuncWasCalled: &timesLogFuncWasCalled,
		logFuncArgs:           &logFuncArgs,
		fatal:                 fatal,
		timesExitWasCalled:    &timesExitWasCalled,
		exitCode:              &exitCode,
		exit:                  exit,
	}
}

func TestRunFatalAndStdHandles(t *testing.T) {
	tle := getTestLogfuncExitType()
	stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	runner.Run()

	tle.LogfuncAndExitTestWascalledsAndArgs(t, 1, []any{`some error`}, 1, 1)

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
	tle := getTestLogfuncExitType()
	osExecCmdRunDouble := &OsExecCmdRunDouble{}
	runner := makeRunnerForCmdRun(osExecCmdRunDouble, &tle)

	runner.Run()

	tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []any{}, 0, 0)

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
			tle := getTestLogfuncExitType()
			runner := makeRunnerForCmdRun(osExecCmdRunDouble, &tle)
			runner.CmdNameArgs = func(dcConfigReader dc_config.Reader) (string, []string) {
				return name, args
			}

			runner.Run()

			tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []any{}, 0, 0)

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

func makeRunnerDry(Getenv func(key string) string, tle *testLogfuncExitType) (*stdHandlesType, *Runner) {
	// var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
	// stdout, stderr, stdin := &stdoutBuf, &stderrBuf, &stdinBuf

	stdHandles, runner := makeRunner(tle)
	runner.OsStruct.Exit = func(code int) {}
	runner.OsStruct.Getenv = Getenv
	fatal := runner.LogFunc
	runner.LogFunc = func(s string) {
		stdHandles.Stderr.WriteString(s)
		fatal(s)
	}

	// fmt.Printf("stdout: '%p', runner.OsStruct.Stdout: '%p' \n", stdout, runner.OsStruct.Stdout)
	// fmt.Printf("stderr: '%p', runner.OsStruct.Stderr: '%p' \n", stderr, runner.OsStruct.Stderr)

	// runner.OsStruct = &osStruct
	return stdHandles, runner
}

var dryGetenv = func(key string) string {
	if key == DryRunEnvVarName {
		return `1`
	}
	return ``
}

func TestStdoutByConfig(t *testing.T) {

	actByGetenv := func(getenv func(string) string) {
		tle := getTestLogfuncExitType()

		stdHandles, runner := makeRunnerDry(getenv, &tle)
		stdout, stderr, _ := stdHandles.Stdout, stdHandles.Stderr, stdHandles.Stdin

		runner.Run()

		tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []any{}, 0, 0)

		if stdout.Len() == 0 {
			t.Errorf("Empty stdout: '%v'", stdout)
		}

		emptyCmd := `["",[]]` + "\n"
		if stdout.String() != emptyCmd {
			t.Errorf("No match of stdout '%v' to empty command: '%v'", stdout, emptyCmd)
		}

		if stderr.Len() != 0 {
			t.Errorf("Non-empty stderr: '%v'", stderr)
		}
	}

	actByGetenv(dryGetenv)
	os.Setenv(DryRunEnvVarName, `1`)
	actByGetenv(os.Getenv)
	os.Unsetenv(DryRunEnvVarName)

}

func TestStdoutByCommand(t *testing.T) {
	namesArgs := [][][]string{{{`docker-compose`}, {`up`}},
		{{`docker-compose`}, {`up`, `-d`}}}
	for _, nameArgs := range namesArgs {
		name, args := nameArgs[0][0], nameArgs[1]
		cmdNameArgs := func(dcConfigReader dc_config.Reader) (string, []string) {
			return name, args
		}

		tle := getTestLogfuncExitType()

		stdHandles, runner := makeRunnerDry(dryGetenv, &tle)
		stdout, stderr := stdHandles.Stdout, stdHandles.Stderr
		runner.CmdNameArgs = cmdNameArgs

		runner.Run()

		tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []any{}, 0, 0)

		if stdout.Len() == 0 {
			t.Errorf("Empty stdout: '%v'", stdout)
		}

		cmdOutputArr, _ := json.Marshal([]any{name, args})
		cmdOutput := string(cmdOutputArr[:]) + "\n"
		if stdout.String() != cmdOutput {
			t.Errorf("No match of stdout '%v' to empty command: '%v'", stdout, cmdOutput)
		}

		if stderr.Len() != 0 {
			t.Errorf("Non-empty stderr: '%v'", stderr)
		}
	}
}

// func getCmdNameArgsByDcyml(osStruct *exec.OsStruct) CmdNameArgsType {
// 	panic("unimplemented")
// }

// func TestCommandByDcyml(t *testing.T) {
// 	cmdNameArgs := getCmdNameArgsByDcyml()
// 	tle := getTestLogfuncExitType()

// 	stdHandles, runner := makeRunnerDry(dryGetenv, &tle)
// 	// stdout, stderr := stdHandles.Stdout, stdHandles.Stderr

// 	cmdNameArgs := getCmdNameArgsByDcyml(runner.OsStruct)
// 	runner.CmdNameArgs = cmdNameArgs

// 	runner.Run()

// 	tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []any{}, 0, 0)

// }
