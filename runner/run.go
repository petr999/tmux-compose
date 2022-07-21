package runner

import (
	"tmux_compose/dc_config"
)

type LogFunc func(v ...any)

func Run(dcConfigReader dc_config.Reader, execStruct execInterface, osStruct *OsStruct, logFunc LogFunc) {
	cmdName, args := cmdNameArgs(dcConfigReader)

	cmd := execStruct.GetCommand(cmdName, args...)
	stdCommute(&cmd, osStruct)

	err := cmd.Obj.Run()
	if err != nil {
		logFunc(err)
	}

}
