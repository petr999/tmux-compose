package cmd_name_args

import (
	"fmt"
	"tmux_compose/dc_config"
	"tmux_compose/types"
)

func CmdNameArgs(dcConfigReader dc_config.ReaderInterface) (types.CmdNameArgsType, error) {
	_, err := dcConfigReader.Read()
	if err != nil {
		return types.CmdNameArgsType{Workdir: ``, Name: ``, Args: nil}, fmt.Errorf(`error reading config: '%w'`, err)
	}

	// ID=0; try_next=1; trap 'echo "trap pid: ${PID}"; kill -INT $PID; try_next="";' SIGINT; while [ 'x1' == "x${try_next}" ]; do sleep 1 & PID=$!; echo "pid: ${PID}" ; wait $PID; done
	return types.CmdNameArgsType{Workdir: ``, Name: `tmux`, Args: make([]string, 0)}, nil
}
