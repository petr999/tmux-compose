package os_struct

import (
	"os"
	"os/exec"
	"tmux_compose/types"
)

type ExecOsStruct struct {
	stdHandles types.StdHandlesType
}

func (execOsStruct ExecOsStruct) Getenv(key string) string {
	return os.Getenv(key)
}

func (execOsStruct ExecOsStruct) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (execOsStruct ExecOsStruct) GetStdHandles() (stdHandles types.StdHandlesType) {
	if execOsStruct.stdHandles == nil {
		execOsStruct.stdHandles = &types.StdHandlesStruct{
			Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin,
		}
	}

	return execOsStruct.stdHandles
}

func (execOsStruct ExecOsStruct) Chdir(dir string) error {
	return os.Chdir(dir)
}

func (execOsStruct ExecOsStruct) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
