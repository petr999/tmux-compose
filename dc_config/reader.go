package dc_config

import "io/fs"

type ReaderInterface interface {
	Read() DcConfigValueType
	SetFqfn(fqfn string)
}

type FsInterface interface {
	ReadFile(fsys fs.FS, name string) ([]byte, error)
}

type DcConfigOsInterface interface {
	Chdir(dir string) error
	Getwd() (dir string, err error)
}

type FsStruct struct {
}
type DcConfig struct {
	OsStruct DcConfigOsInterface
	FsStruct FsInterface
	Fqfn     string
}

type DcConfigValueType = map[string]interface{}

func (dcConfig DcConfig) Read() DcConfigValueType {
	var value = make(DcConfigValueType)
	return value
}

func (dcConfig DcConfig) SetFqfn(fqfn string) {}
