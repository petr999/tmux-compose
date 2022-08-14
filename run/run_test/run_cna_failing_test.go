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
	"tmux_compose/types"

	_ "embed"
)

type dcYmlOsGetwdRootDouble struct {
	dcYmlOsFailingDouble
}

// ReadFile implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) ReadFile(name string) ([]byte, error) {
	if name == `/docker-compose.yml` {
		return dcContents, nil
	}
	return []byte{}, fmt.Errorf("Wrong path to Dc ReadFile(): '%v'", name)
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Getwd() (dir string, err error) {
	return `/`, nil
}

// Stat implements types.DcYmlOsInterface
func (osStruct *dcYmlOsGetwdRootDouble) Stat(name string) (dfi types.FileInfoStruct, err error) {
	if name == `/` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return true
			},
			IsFile: func() bool {
				return false
			},
		}, nil
	} else if name == `/docker-compose.yml` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return false
			},
			IsFile: func() bool {
				return true
			},
		}, nil
	}
	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
}

func TestRunCnaWorkdirRoot(t *testing.T) {
	configStruct := config.Construct(ConfigOsDouble{})

	dcYml := dc_yml.Construct(&dcYmlOsGetwdRootDouble{}, configStruct)
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

	if os.ExitData.code != 1 {
		t.Errorf(`Workdir '/' provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Workdir '/' called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	stderrExpected := "Get command name and args error: error finding base dir name '/' same length for work dir: '/'\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing DcOsStruct.Getwd() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type cnaOsFailingStat struct {
	cnaOsFailingDouble
}

// Stat implements types.DcYmlOsInterface
func (osStruct *cnaOsFailingStat) Stat(name string) (dfi types.FileInfoStruct, err error) {
	return dfi, fmt.Errorf("not found")
}

type configOsCnaTemplate struct {
	ConfigOsDouble
}

func (configOsCnaTemplate) Getenv(name string) string {
	if name == `TMUX_COMPOSE_TEMPLATE_FNAME` {
		return `/path/to/dumbclicker/tmux-compose-template.gson`
	}
	return ``
}

func TestRunCnaOsFailingStat(t *testing.T) {

	configStruct := config.Construct(configOsCnaTemplate{})

	dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{}, configStruct)
	cna := cmd_name_args.Construct(&cnaOsFailingStat{}, configStruct)
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

	if os.ExitData.code != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get command name and args error: error reading file '/path/to/dumbclicker/tmux-compose-template.gson': not found\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type cnaOsStatFile struct{}

// Stat implements types.DcYmlOsInterface
func (osStruct *cnaOsStatFile) Stat(name string) (dfi types.FileInfoStruct, err error) {
	if name == `/path/to/dumbclicker/tmux-compose-template.gson` {
		return types.FileInfoStruct{
			IsDir: func() bool {
				return false
			},
			IsFile: func() bool {
				return true
			},
		}, nil
	}
	return dfi, fmt.Errorf("Failed to Stat() path: '%v':", name)
}

type cnaOsFailingReadFile struct {
	cnaOsStatFile
	ReadFileData struct {
		WasCalled int
		Args      []string
	}
}

// ReadFile implements types.CnaOsInterface
func (cnaOsStruct cnaOsFailingReadFile) ReadFile(name string) ([]byte, error) {
	cnaOsStruct.ReadFileData.WasCalled++
	cnaOsStruct.ReadFileData.Args = append(cnaOsStruct.ReadFileData.Args, name)
	if name == `/path/to/dumbclicker/tmux-compose-template.gson` {
		return []byte{}, fmt.Errorf("permission denied")
	} else {
		return []byte{}, fmt.Errorf(`Wrong file name: '%v' not '/path/to/dumbclicker/tmux-compose-template.gson'`, name)
	}
}

type ConfigOsCnaOsGetenvFile struct {
	ConfigOsDouble
}

func (osStruct *ConfigOsCnaOsGetenvFile) Getenv(name string) string {
	if name == `TMUX_COMPOSE_TEMPLATE_FNAME` {
		return `/path/to/dumbclicker/tmux-compose-template.gson`
	}
	return ``
}

func TestRunCnaOsFailingReadFile(t *testing.T) {

	configStruct := config.Construct(&ConfigOsCnaOsGetenvFile{})

	dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{}, configStruct)
	cnaOsStruct := &cnaOsFailingReadFile{}
	cna := cmd_name_args.Construct(cnaOsStruct, configStruct)
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

	if cnaOsStruct.ReadFileData.WasCalled != 1 {
		if os.ExitData.code != 1 {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called not '1' time(s)  but: '%v'`, os.ExitData.code)
		}
	} else {
		if cnaOsStruct.ReadFileData.Args[0] != `/path/to/dumbclicker/tmux-compose-template.gson` {
			t.Errorf(`Failing CnaOsStruct.ReadFile() was called with '%v' arg(s) not with: '%v'`, cnaOsStruct.ReadFileData.Args[0], `/path/to/dumbclicker/tmux-compose-template.gson`)
		}
	}

	if os.ExitData.code != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
	}
	if os.ExitData.wasCalledTimes != 1 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
	}
	if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
	}
	stderrExpected := "Get command name and args error: error reading file '/path/to/dumbclicker/tmux-compose-template.gson': permission denied\n"
	if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected {
		t.Errorf(`Failing CnaOsStruct.ReadFile() made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected)
	}
}

type cnaOsStructFailingTmplParse struct {
	cnaOsStatFile
}

func (cnaOsStruct *cnaOsStructFailingTmplParse) ReadFile(name string) ([]byte, error) {
	return []byte(`{{define}}`), nil
}

type cnaOsStructFailingTmplExecute struct {
	cnaOsStatFile
}

func (cnaOsStruct *cnaOsStructFailingTmplExecute) ReadFile(name string) ([]byte, error) {
	return []byte(`{{.Nonexistent}}`), nil
}

type cnaOsStructFailingTmplJson struct {
	cnaOsStatFile
}

func (cnaOsStruct *cnaOsStructFailingTmplJson) ReadFile(name string) ([]byte, error) {
	return []byte(`{`), nil
}

type cnaOsStructFailingTmplSingleCommand struct {
	cnaOsStatFile
}

func (cnaOsStruct *cnaOsStructFailingTmplSingleCommand) ReadFile(name string) ([]byte, error) {
	return []byte(`[{"Cmd":"","Args":[]},{"Cmd":"","Args":[]}]`), nil
}

func TestRunCnaOsFailingTmpl(t *testing.T) {

	for _, stdErrExpectedAndCnaOsStruct := range []map[string]types.CnaOsInterface{
		{"Get command name and args error: error executing template '{{define}}' on with vars from: '{/path/to/dumbclicker [nginx h2o dumbclicker] dumbclicker}': error reading name template: '{{define}}': template: tmux_compose:1: unexpected \"}}\" in define clause": &cnaOsStructFailingTmplParse{}},
		{"Get command name and args error: error executing template '{{.Nonexistent}}' on with vars from: '{/path/to/dumbclicker [nginx h2o dumbclicker] dumbclicker}': error executing name template: '{{.Nonexistent}}' on '{/path/to/dumbclicker [nginx h2o dumbclicker] dumbclicker}': template: tmux_compose:1:2: executing \"tmux_compose\" at <.Nonexistent>: can't evaluate field Nonexistent in type cmd_name_args.dcvBasedirType": &cnaOsStructFailingTmplExecute{}},
		{"Get command name and args error: error unserializing template '{': unexpected end of JSON input": &cnaOsStructFailingTmplJson{}},
		{`Get command name and args error: error unserializing template '[{"Cmd":"","Args":[]},{"Cmd":"","Args":[]}]': amount of commands '2' is not '1' from '[{  []} {  []}]'`: &cnaOsStructFailingTmplSingleCommand{}},
	} {

		for stderrExpected, cnaOsStruct := range stdErrExpectedAndCnaOsStruct {
			configStruct := config.Construct(&ConfigOsCnaOsGetenvFile{})

			dcYml := dc_yml.Construct(&dcYmlOsGetwdDouble{}, configStruct)
			cna := cmd_name_args.Construct(cnaOsStruct, configStruct)
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

			if os.ExitData.code != 1 {
				t.Errorf(`Failing exec template was provided not '1' to Runner.Os.Exit exit code but: '%v'`, os.ExitData.code)
			}
			if os.ExitData.wasCalledTimes != 1 {
				t.Errorf(`Failing exec template was called Runner.Os.Exit not '1' time: '%v'`, os.ExitData.code)
			}
			if execOsStruct.StdHandlesDouble.Stdout.Len() != 0 {
				t.Errorf(`Failing exec template made stdout not empty: '%s'`, execOsStruct.StdHandlesDouble.Stdout)
			}
			if execOsStruct.StdHandlesDouble.Stderr.String() != stderrExpected+"\n" {
				t.Errorf(`Failing exec template made stderr '%s' not equal to: '%s'`, execOsStruct.StdHandlesDouble.Stderr, stderrExpected+"\n")
			}
		}
	}
}
