package config

import (
	"tmux_compose/types"
)

type ConfigStruct struct {
	osStruct types.ConfigOsInterface
}

func Construct(osStruct types.ConfigOsInterface) *ConfigStruct {
	config := &ConfigStruct{}
	config.New(osStruct)
	return config
}

func (config *ConfigStruct) New(osStruct types.ConfigOsInterface) {
	config.osStruct = osStruct
}

func (config *ConfigStruct) GetCnaTemplateFname() string {
	return config.osStruct.Getenv(`TMUX_COMPOSE_TEMPLATE_FNAME`)
}

func (config *ConfigStruct) GetDcYmlFname() string {
	return config.osStruct.Getenv(`TMUX_COMPOSE_DC_YML`)
}

func (config *ConfigStruct) GetDryRun() string {
	return config.osStruct.Getenv(`TMUX_COMPOSE_DRY_RUN`)
}

func (config *ConfigStruct) GetShell() (val string) {
	val = config.osStruct.Getenv(`SHELL`)
	// if len(val) == 0 {
	// 	val = `/usr/bin/env bash`
	// }
	return
}
