package os_struct

import "os"

type RunnerOsStruct struct{}

func (osStruct RunnerOsStruct) Exit(code int) {
	os.Exit(code)
}
