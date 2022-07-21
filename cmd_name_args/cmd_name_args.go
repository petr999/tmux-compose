package cmd_name_args

import "tmux_compose/dc_config"

func CmdNameArgs(dcConfigReader dc_config.Reader) (string, []string) {
	return `id`, make([]string, 0)
}
