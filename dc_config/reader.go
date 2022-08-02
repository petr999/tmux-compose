package dc_config

import (
	"fmt"
	"path/filepath"

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

func (dcConfig DcConfig) readConfigFile() (value DcConfigValueType, err error) {
	var buf []byte
	fqfn := dcConfig.Fqfn

	buf, err = dcConfig.OsStruct.ReadFile(fqfn)
	if err != nil {
		return nil, fmt.Errorf("reading config file: '%s' error:\n\t%w", fqfn, err)
	}

	err = yaml.Unmarshal(buf, &value)
	if err != nil {
		return nil, fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
	}

	return
}

func (dcConfig DcConfig) Read() (DcConfigValueType, error) {
	dir := filepath.Dir(dcConfig.Fqfn)
	err := dcConfig.OsStruct.Chdir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to change to dir: '%s' error:\n\t%w", dir, err)
	}

	return dcConfig.readConfigFile()
}

func (dcConfig DcConfig) SetFqfn(fqfn string) {}
