package os_struct

import "os"

type DcYmlOsStruct struct{}

func (osStruct DcYmlOsStruct) Chdir(dir string) error               { return os.Chdir(dir) }
func (osStruct DcYmlOsStruct) Getwd() (dir string, err error)       { return os.Getwd() }
func (osStruct DcYmlOsStruct) ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }
