package os_struct

import "os"

type CnaOsStruct struct{}

func (osStruct CnaOsStruct) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
