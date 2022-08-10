package config

import "os"

type ConfigStruct struct{}

func (ConfigStruct) GetCnaTemplateFname() string {
	return os.Getenv(`TMUX_COMPOSE_TEMPLATE_FNAME`)
}
