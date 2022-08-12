package run_test

import (
	"fmt"
	"testing"
	"tmux_compose/cmd_name_args"
	"tmux_compose/config"
	"tmux_compose/dc_yml"
	"tmux_compose/exec"
	"tmux_compose/logger"

	"tmux_compose/run"
)

// const loremString = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`

// type stdHandlesType struct {
// 	Stdout *bytes.Buffer
// 	Stderr *bytes.Buffer
// 	Stdin  *bytes.Buffer
// }

// var projDir = getProjDir()
// var dcConfigSample = getDcConfigSample()

// var cmdNameArgsEmpty = types.CmdNameArgsType{Workdir: ``, Name: ``, Args: []string{}}

// func getProjDir() string {
// 	_, testFilename, _, _ := runtime.Caller(0)
// 	projDir, err := filepath.Abs(filepath.Join(filepath.Dir(testFilename), `..`))
// 	if err != nil {
// 		panic(fmt.Sprintf("Empty project directory: '%v'", projDir))
// 	}
// 	return projDir
// }

// func readTestdataFile(fname string) string {
// 	fqfn := filepath.Join(projDir, `testdata`, fname)
// 	outputBytes, _ := ioutil.ReadFile(fqfn)
// 	if len(outputBytes) == 0 {
// 		panic(fmt.Sprintf("Empty test data file: '%v'", fqfn))
// 	}

// 	return string(outputBytes)
// }

// func getDcConfigSample() string {
// 	fname := filepath.Join(`dumbclicker`, `docker-compose.yml`)
// 	return readTestdataFile(fname)
// }

// func makeRunnerCommon() (*stdHandlesType, *Runner) {
// 	var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
// 	stdout, stderr, stdin := &stdoutBuf, &stderrBuf, &stdinBuf
// 	stdHandles := stdHandlesType{stdout, stderr, stdin}

// 	return &stdHandles, &Runner{
// 		DcConfigReader: dc_config.DcConfig{},
// 		ExecStruct:     &exec.ExecStruct{},
// 		OsStruct: &types.OsStruct{
// 			Stdout: stdout,
// 			Stderr: stderr,
// 			Stdin:  stdin,
// 			Exit:   func(code int) {},
// 			Getenv: func(string) string { return `` },
// 			Chdir:  func(string) error { return nil },
// 		},
// 		CmdNameArgs: func(dcConfigReader dc_config.ReaderInterface, getTmplFuncs []func() string) (types.CmdNameArgsType, error) {
// 			return cmdNameArgsEmpty, nil
// 		},
// 	}

// }

// func makeRunner(tle *testLogfuncExitType) (stdHandles *stdHandlesType, runner *Runner) {
// 	stdHandles, runner = makeRunnerCommon()

// 	runner.OsStruct.Exit = tle.exit
// 	runner.LogFunc = tle.fatal

// 	return stdHandles, runner
// }

// func makeRunnerForFatal(cmdName string, tle *testLogfuncExitType) (*stdHandlesType, *Runner) {
// 	stdHandles, runner := makeRunner(tle)

// 	runner.CmdNameArgs = func(dcConfigReader dc_config.ReaderInterface, getTmplFuncs []func() string) (types.CmdNameArgsType, error) {
// 		return cmdNameArgsEmpty, nil
// 	}

// 	return stdHandles, runner
// }

// type OsExecCmdRunDouble struct {
// 	timesOsExecCmdRunWasCalled uint8
// 	nameOsExecCommandWasCalled string
// 	argsOsExecCommandWasCalled []string
// 	Stdout, Stderr             io.Writer
// 	Stdin                      io.Reader
// }

// func (Obj *OsExecCmdRunDouble) Run() error {
// 	Obj.timesOsExecCmdRunWasCalled++
// 	return nil
// }

// type ExecStructDouble struct {
// 	osExecCmdRunDouble *OsExecCmdRunDouble
// 	cmd                *exec.CmdType
// 	Stdout, Stderr     io.Writer
// 	Stdin              io.Reader
// }

// func (execStructDouble *ExecStructDouble) MakeCommand(dryRun *exec.MakeCommandDryRunType,
// 	nameArgs types.CmdNameArgsType) {

// 	osExecCmdRunDouble := execStructDouble.osExecCmdRunDouble
// 	osExecCmdRunDouble.nameOsExecCommandWasCalled = nameArgs.Name
// 	osExecCmdRunDouble.argsOsExecCommandWasCalled = nameArgs.Args

// 	cmd := &exec.CmdType{
// 		Obj:    osExecCmdRunDouble,
// 		Stdout: &execStructDouble.Stdout,
// 		Stderr: &execStructDouble.Stderr,
// 		Stdin:  &execStructDouble.Stdin,
// 	}

// 	execStructDouble.SetCommand(cmd)
// }

// func (execStructDouble *ExecStructDouble) GetCommand() *exec.CmdType {
// 	return execStructDouble.cmd
// }

// func (execStructDouble *ExecStructDouble) SetCommand(cmd *exec.CmdType) {
// 	execStructDouble.cmd = cmd
// }

// func makeRunnerForCmdRun(osExecCmdRunDouble *OsExecCmdRunDouble, tle *testLogfuncExitType) *Runner {
// 	_, runner := makeRunner(tle)

// 	runner.ExecStruct = &ExecStructDouble{
// 		osExecCmdRunDouble: osExecCmdRunDouble,
// 		Stdout:             osExecCmdRunDouble.Stdout,
// 		Stderr:             osExecCmdRunDouble.Stderr,
// 		Stdin:              osExecCmdRunDouble.Stdin,
// 	}

// 	return runner
// }

// type testLogfuncExitType struct {
// 	timesLogFuncWasCalled *int
// 	logFuncArgs           *[]any
// 	fatal                 func(s string)
// 	timesExitWasCalled    *int
// 	exitCode              *int
// 	exit                  func(code int)
// }

// func (tle *testLogfuncExitType) LogfuncAndExitTestWascalledsAndArgs(t *testing.T, tlfwcExpected int, lfaExpected []string, tewcExpected int, exitCodeExpected int) {

// 	if *tle.timesLogFuncWasCalled != tlfwcExpected {
// 		t.Errorf(`Log func was called '%v' times`, *tle.timesLogFuncWasCalled)
// 	}
// 	if len(*tle.logFuncArgs) != len(lfaExpected) {
// 		t.Errorf(`Wrong argument of log function: '%v'`, *tle.logFuncArgs)
// 	}
// 	if *tle.timesExitWasCalled != tewcExpected {
// 		t.Errorf(`Exit func was called '%v' times`, *tle.timesExitWasCalled)
// 	}
// 	if *tle.exitCode != exitCodeExpected {
// 		t.Errorf(`Wrong argument of Exit function: '%v'`, *tle.exitCode)
// 	}

// }

// func getTestLogfuncExitType() testLogfuncExitType {
// 	timesLogFuncWasCalled := 0
// 	var logFuncArgs []any
// 	var fatal LogFuncType = func(s string) {
// 		timesLogFuncWasCalled++
// 		logFuncArgs = append(logFuncArgs, s)
// 	}

// 	timesExitWasCalled := 0
// 	var exitCode int
// 	exit := func(code int) {
// 		timesExitWasCalled++
// 		exitCode = code
// 	}
// 	return testLogfuncExitType{
// 		timesLogFuncWasCalled: &timesLogFuncWasCalled,
// 		logFuncArgs:           &logFuncArgs,
// 		fatal:                 fatal,
// 		timesExitWasCalled:    &timesExitWasCalled,
// 		exitCode:              &exitCode,
// 		exit:                  exit,
// 	}
// }

type dcYmlOsFailingGetwdDouble struct {
	dcYmlOsFailingDouble
	GetwdData struct{ WascalledTimes int }
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingGetwdDouble) Getwd() (dir string, err error) {
	osStruct.GetwdData.WascalledTimes++
	return ``, fmt.Errorf(`current working directory not found`)
}

func TestRunDcOsGetwdFail(t *testing.T) { // AndStdHandles {
	// tle := getTestLogfuncExitType()
	// stdHandles, runner := makeRunnerForFatal(`/\\nonexistent`, &tle)

	dcYmlOsStruct := &dcYmlOsFailingGetwdDouble{}
	configStruct := config.Construct(ConfigOsDouble{})
	dcYml := dc_yml.Construct(dcYmlOsStruct, configStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingDouble{}, configStruct)
	execOsStruct := &execOsFailingDouble{}
	exec := exec.Construct(execOsStruct, configStruct)

	os := &osDouble{}
	logger := logger.Construct(execOsStruct.GetStdHandles())

	runner := run.Runner{
		CmdNameArgs: cna,
		DcYml:       dcYml,
		Exec:        exec,
		Os:          os,
		Logger:      logger,
	}

	runner.Run()

	if dcYmlOsStruct.GetwdData.WascalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called not '1' time but: '%v'`, dcYmlOsStruct.GetwdData.WascalledTimes)
	}
	if os.ExitData.code != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing DcOsStruct.Getwd() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing DcOsStruct.Getwd() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get docker-compose config error: getting current working directory: 'current working directory not found'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

// func getLongArgs(amount int) []string {
// 	longArgs := make([]string, amount)
// 	loremArr := regexp.MustCompile(`\s+`).Split(loremString, amount)
// 	half := amount / 2
// 	iter := func(i int, sign bool) string {
// 		var idx int
// 		j := (i + 1) / 2
// 		if sign {
// 			idx = half + j
// 		} else {
// 			idx = half - j
// 		}
// 		return loremArr[idx]
// 	}

// 	for i := range longArgs {
// 		sign := false
// 		if i%2 != 0 {
// 			sign = true
// 		}
// 		longArgs[i] = iter(i, sign)
// 	}
// 	if len(longArgs) < amount {
// 		i := half + 1
// 		longArgs = append(longArgs, iter(i, true))
// 	}

// 	perms := rand.Perm(amount)
// 	for i, v := range perms {
// 		longArgs[i] = longArgs[i] + strings.Repeat(` `, v)
// 	}

// 	return longArgs
// }

// func TestCmdRunWasCalledWithArgs(t *testing.T) {
// 	const amount = 55
// 	longArgs := getLongArgs(amount)

// 	// name, args
// 	cmdsArgs := []map[string][]string{
// 		{`testingSchmesting`: {`--arg00`}},
// 		{`testing Schmesting`: {}},
// 		{``: {`--arg00`}},
// 		{``: {`--vaudeville 577`}},
// 		{`testingSchmesting`: {`--arg00`, `abcd`, `vaudeville 577`, `35`, `-88.55000`}},
// 		{`docker-compose`: longArgs},
// 	}

// 	for _, cmdArgs := range cmdsArgs {
// 		for name, args := range cmdArgs {
// 			osExecCmdRunDouble := &OsExecCmdRunDouble{}
// 			tle := getTestLogfuncExitType()
// 			runner := makeRunnerForCmdRun(osExecCmdRunDouble, &tle)
// 			runner.CmdNameArgs = func(dcConfigReader dc_config.ReaderInterface, getTmplFuncs []func() string) (types.CmdNameArgsType, error) {
// 				return types.CmdNameArgsType{Workdir: ``, Name: name, Args: args}, nil
// 			}

// 			runner.Run()

// 			tle.LogfuncAndExitTestWascalledsAndArgs(t, 0, []string{}, 0, 0)

// 			if osExecCmdRunDouble.timesOsExecCmdRunWasCalled != 1 {
// 				t.Errorf(`os/exec.Command().Run() was called not once but '%v' times`, osExecCmdRunDouble.timesOsExecCmdRunWasCalled)
// 			}
// 			if osExecCmdRunDouble.nameOsExecCommandWasCalled != name {
// 				t.Errorf(`os/exec.Run.Command() was called with name other than '%v'`, osExecCmdRunDouble.nameOsExecCommandWasCalled)
// 			}
// 			if !reflect.DeepEqual(osExecCmdRunDouble.argsOsExecCommandWasCalled, args) {
// 				t.Errorf(`os/exec.Run.Command() was called with arg '%v' other than '%v'`, osExecCmdRunDouble.argsOsExecCommandWasCalled, []string{})
// 			}
// 		}
// 	}
// }
