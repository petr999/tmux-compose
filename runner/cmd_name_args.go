package runner

import "tmux_compose/dc_config"

func cmdNameArgs(dcConfigReader dc_config.Reader) (string, []string) {
	return `id`, make([]string, 0)
}
