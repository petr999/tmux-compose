package types

import (
	"io"
	"os"
)

type OsStructExit func(code int)
type OsStructGetEnv func(key string) string
type OsStructChdir func(dir string) error
type OsStructGetwd func() (dir string, err error)
type OsStructReadFile func(name string) ([]byte, error)

type OsStruct struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Stdin    io.Reader
	Exit     OsStructExit
	Getenv   OsStructGetEnv
	Chdir    OsStructChdir
	Getwd    OsStructGetwd
	ReadFile OsStructReadFile
}

func MakeOsStruct() *OsStruct {
	return &OsStruct{Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin, Exit: os.Exit, Getenv: os.Getenv, Chdir: os.Chdir, Getwd: os.Getwd, ReadFile: os.ReadFile}
}
