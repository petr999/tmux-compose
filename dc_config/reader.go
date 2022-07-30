package dc_config

type ReaderInterface interface {
	Read() DcConfigValueType
	SetFqfn(fqfn string)
}

type DcConfigOsInterface interface {
	Chdir(dir string) error
	Getwd() (dir string, err error)
	ReadFile(name string) ([]byte, error)
}

type DcConfig struct {
	OsStruct DcConfigOsInterface
	Fqfn     string
}

type DcConfigValueType = map[string]interface{}

func (dcConfig DcConfig) Read() DcConfigValueType {
	var value = make(DcConfigValueType)
	return value
}

func (dcConfig DcConfig) SetFqfn(fqfn string) {}
