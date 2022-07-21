package main

import (
	"log"
	"os"
	"tmux_compose/dc_config"
	"tmux_compose/runner"
)

func main() {
	osStruct := runner.OsStruct{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	execStruct := runner.ExecStruct{}

	runner.Run(dc_config.DcConfig{}, &execStruct, &osStruct, log.Fatal)
}
