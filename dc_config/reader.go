package dc_config

import (
	"fmt"
	"path/filepath"
	"regexp"

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

type DcConfigValueType = struct {
	WorkDir         string
	DcServicesNames map[string]interface{} `json:"services"`
} // map[string]interface{}

type DcConfigFileType = struct {
	Services map[string]map[string]interface{}
}

func (dcConfig DcConfig) readConfigFile() (value DcConfigValueType, err error) {
	var buf []byte
	fqfn := dcConfig.Fqfn

	buf, err = dcConfig.OsStruct.ReadFile(fqfn)
	if err != nil {
		return DcConfigValueType{}, fmt.Errorf("reading config file: '%s' error:\n\t%w", fqfn, err)
	}

	err = yaml.Unmarshal(buf, &value)
	if err != nil {
		return DcConfigValueType{}, fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
	}

	return
}

func (dcConfig DcConfig) Read() (DcConfigValueType, error) {
	dir := filepath.Dir(dcConfig.Fqfn)
	err := dcConfig.OsStruct.Chdir(dir)
	if err != nil {
		return DcConfigValueType{}, fmt.Errorf("failed to change to dir: '%v' error:\n\t%w", dir, err)
	}
	workDir, err := dcConfig.OsStruct.Getwd()
	if err != nil {
		return DcConfigValueType{}, fmt.Errorf("failed to get current directory:\n\t%w", err)
	}
	if len(workDir) == 0 {
		return DcConfigValueType{}, fmt.Errorf("work directory is empty")
	}

	value, err := dcConfig.readConfigFile()
	if err != nil {
		return value, err
	}

	if len(value.DcServicesNames) == 0 {
		return DcConfigValueType{}, fmt.Errorf("no service names in config: '%v'", dcConfig.Fqfn)
	} else {
		for serviceName := range value.DcServicesNames {
			if len(serviceName) == 0 || !regexp.MustCompile(`^\w[\w\d]*$`).MatchString(serviceName) {
				return DcConfigValueType{}, fmt.Errorf("empty or inappropriate service name '%v' in config: '%v'", serviceName, dcConfig.Fqfn)
			}
		}
	}

	value.WorkDir = workDir

	return value, err
}

func (dcConfig DcConfig) SetFqfn(fqfn string) {}
