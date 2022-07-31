package dc_config

import (
	"fmt"

	"github.com/ghodss/yaml"
)

type ReaderInterface interface {
	Read() (DcConfigValueType, error)
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

func (dcConfig DcConfig) Read() (DcConfigValueType, error) {
	var err error
	var buf []byte
	var value DcConfigValueType
	fqfn := dcConfig.Fqfn
	buf, err = dcConfig.OsStruct.ReadFile(fqfn)
	if err != nil {
		return nil, fmt.Errorf("reading config file: '%s' error:\n\t%w", fqfn, err)
	}

	err = yaml.Unmarshal(buf, &value)
	if err != nil {
		return nil, fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
	}
	return value, nil
}

func (dcConfig DcConfig) SetFqfn(fqfn string) {}
