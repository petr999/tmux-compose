package cmd_name_args

import "tmux_compose/dc_config"

func CmdNameArgs(dcConfigReader dc_config.Reader) (string, []string) {
	// ID=0; try_next=1; trap 'echo "trap pid: ${PID}"; kill -INT $PID; try_next="";' SIGINT; while [ 'x1' == "x${try_next}" ]; do sleep 1 & PID=$!; echo "pid: ${PID}" ; wait $PID; done
	return `id`, make([]string, 0)
}
