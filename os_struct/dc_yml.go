package os_struct

import (
	"fmt"
	"os"
	"tmux_compose/types"
)

type DcYmlOsStruct struct{}

func (osStruct DcYmlOsStruct) Getwd() (dir string, err error)       { return os.Getwd() }
func (osStruct DcYmlOsStruct) ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }
func (osStruct DcYmlOsStruct) Stat(name string) (dfi types.FileInfoStruct, err error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		err = fmt.Errorf(`stat file name: '%v' error: '%v'`, name, err)
		return
	}
	return types.FileInfoStruct{IsDir: func() bool { return fileInfo.IsDir() }}, nil
}

// type FileInfoStruct struct{}
// func (fileInfo FileInfoStruct) IsDir() bool {}
