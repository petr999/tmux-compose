package run_test

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"tmux_compose/types"
)

type cnaOsFailingDouble struct{}

// ReadFile implements types.CnaOsInterface
func (cnaOsFailingDouble) ReadFile(name string) ([]byte, error) {
	return []byte{}, fmt.Errorf("unimplemented")
}

func (cnaOsFailingDouble) Stat(name string) (fileInfo types.FileInfoStruct, err error) {
	return fileInfo, fmt.Errorf(`unimplemented`)
}

type dcYmlOsFailingDouble struct {
}

// Getwd implements types.DcYmlOsInterface
func (osStruct *dcYmlOsFailingDouble) Getwd() (dir string, err error) {
	err = fmt.Errorf("unimplemented")
	return
}

// ReadFile implements types.DcYmlOsInterface
func (*dcYmlOsFailingDouble) ReadFile(name string) ([]byte, error) {
	return []byte{}, fmt.Errorf("unimplemented")
}

func (*dcYmlOsFailingDouble) Getenv(string) string { return `` }

func (*dcYmlOsFailingDouble) Stat(name string) (dfi types.FileInfoStruct, err error) {
	return dfi, fmt.Errorf("unimplemented")
}

type execOsDouble struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

// Chdir implements types.ExecOsInterface
func (execOsDouble) Chdir(dir string) error {
	return fmt.Errorf("unimplemented")
}

// GetStdHandles implements types.ExecOsInterface
func (execOsDouble) GetStdHandles() types.StdHandlesType {
	panic("unimplemented")
}

// Getenv implements types.ExecOsInterface
func (execOsDouble) Getenv(key string) string {
	return ``
}

// ReadFile implements types.ExecOsInterface
func (execOsDouble) ReadFile(name string) ([]byte, error) {
	return []byte{}, fmt.Errorf("unimplemented")
}

func (execOsDouble) Command(name string, arg ...string) *exec.Cmd {
	return &exec.Cmd{}
}

type stdHandlesDoubleStruct struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Stdin  *bytes.Buffer
}

type stdHandlesDoubleType *stdHandlesDoubleStruct

type execOsFailingDouble struct {
	StdHandlesDouble stdHandlesDoubleType
	execOsDouble
}

// GetStdHandles implements types.ExecOsInterface
func (osStruct *execOsFailingDouble) GetStdHandles() types.StdHandlesType {
	if osStruct.StdHandlesDouble == nil {
		stdout, stderr, stdin := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
		osStruct.StdHandlesDouble = &stdHandlesDoubleStruct{stdout, stderr, stdin}
	}
	return &types.StdHandlesStruct{
		Stdout: osStruct.StdHandlesDouble.Stdout,
		Stderr: osStruct.StdHandlesDouble.Stderr,
		Stdin:  osStruct.StdHandlesDouble.Stdin,
	}
}

type osDouble struct {
	ExitData struct {
		wasCalledTimes int
		code           int
	}
}

func (os *osDouble) Exit(code int) {
	os.ExitData.wasCalledTimes++
	os.ExitData.code = code
}

type ConfigOsDouble struct{}

func (ConfigOsDouble) Getenv(string) string { return `` }
