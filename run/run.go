package run

import (
	"fmt"
	"tmux_compose/types"
)

const DryRunEnvVarName = `TMUX_COMPOSE_DRY_RUN`

// type LogFuncType func(s string)

// type CmdNameArgsFunc func(dcConfigReader dc_config.ReaderInterface, osStruct cmd_name_args.OsStructCmdNameArgs) (types.CmdNameArgsType, error)

type Runner struct {
	CmdNameArgs types.CnaInterface
	DcYml       types.DcYmlInterface
	Exec        types.ExecInterface
	Os          types.RunnerOsInterface
	Logger      types.LoggerInterface
}

func (runner *Runner) runWithExitcode() int {
	log := runner.Logger.Log

	dcYmlValue, err := runner.DcYml.Get()
	if err != nil {
		log(fmt.Sprintf("%v\n", err))
		return 1
	}
	cna, err := runner.CmdNameArgs.Get(dcYmlValue)
	if err != nil {
		log(fmt.Sprintf("%v,\n", err))
		return 1
	}

	cmd := runner.Exec.GetCommand(cna)

	err = cmd.Obj.Run()
	if err != nil {
		log(fmt.Sprintf("%v,\n", err))
		return 1
	}

	// DcConfigReader, ExecStruct, OsStruct, LogFunc := runner.DcConfigReader, runner.ExecStruct, runner.OsStruct, runner.LogFunc
	// CmdNameArgs := runner.CmdNameArgs
	// cmdNameArgs, err := CmdNameArgs(DcConfigReader, cmd_name_args.OsStructCmdNameArgs)
	// if err != nil {
	// 	LogFunc(fmt.Sprintf("%v,\n", err))
	// 	return 1
	// }

	// ExecStruct.MakeCommand(&exec.MakeCommandDryRunType{DryRun: OsStruct.Getenv(DryRunEnvVarName), OsStruct: OsStruct},
	// 	cmdNameArgs)
	// cmd := ExecStruct.GetCommand()
	// cmd.StdCommute(OsStruct)

	// err = cmd.Run()

	// if err != nil {
	// 	LogFunc(fmt.Sprintf("%v,\n", err))
	// 	return 1
	// }

	return 0
}

func (runner *Runner) Run() {
	runner.Os.Exit(runner.runWithExitcode())
}
