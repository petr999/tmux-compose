package os_struct

import "os"

type ConfigOsStruct struct{}

var envVars map[string]bool = map[string]bool{
	`TMUX_COMPOSE_TEMPLATE_FNAME`: true,
	`TMUX_COMPOSE_DC_YML`:         true,
	`TMUX_COMPOSE_DRY_RUN`:        true,
}

func (osStruct ConfigOsStruct) Getenv(name string) string {
	if envVars[name] {
		return os.Getenv(name)
	}
	return ``
}
