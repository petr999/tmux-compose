package runner

import (
	"log"
	"os"
	"os/exec"
)

func Run() {
	cmd := exec.Command(`id`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
