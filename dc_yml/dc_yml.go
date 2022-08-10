package dc_yml

import (
	"fmt"
	"path/filepath"
	"regexp"
	"tmux_compose/types"

	"gopkg.in/yaml.v2"
)

type DcYmlValue = struct {
	Workdir         string
	DcServicesNames []string
	config          types.ConfigOsInterface
}

type DcYml struct {
	osStruct types.DcYmlOsInterface
	config   types.ConfigInterface
}

func (dcYml *DcYml) New(osStruct types.DcYmlOsInterface, config types.ConfigInterface) {
	dcYml.config = config
	dcYml.osStruct = osStruct
}

func Construct(osStruct types.DcYmlOsInterface, config types.ConfigInterface) *DcYml {
	dcYml := &DcYml{}
	dcYml.New(osStruct, config)
	return dcYml
}

func (dcYml *DcYml) getServiceNames(fName string) (dcYmlValue types.DcYmlValue, err error) {
	dcYmlValue.Workdir = filepath.Dir(fName)
	buf, err := dcYml.osStruct.ReadFile(fName)
	if err != nil {
		return
	}

	mapSlices := yaml.MapSlice{}

	err = yaml.Unmarshal(buf, &mapSlices)
	if err != nil {
		return // DcConfigValueType{}, fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
	}

	services, err := getServices(mapSlices)
	if err != nil { // err = fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
		return
	}

	dcYmlValue.DcServicesNames = services

	return
}

// if len(dcYml.osStruct.Getenv(dcEnvVar)) > 0 {
func (dcYml *DcYml) getByFname(fName string) (dcYmlValue types.DcYmlValue, err error) {
	fPath, err := filepath.Abs(fName)
	if err != nil {
		return dcYmlValue, fmt.Errorf(`normalizing path: '%v' error: '%w'`, fPath, err)
	}

	fileInfo, err := dcYml.osStruct.Stat(fPath)
	if err != nil {
		return dcYmlValue, err
	}

	isDir := false
	if fileInfo.IsDir() {
		isDir = true
		fPath = filepath.Join(fPath, `docker-compose.yml`)
		fPath, err = filepath.Abs(fPath)
		if err != nil {
			return dcYmlValue, err
		}
		fileInfo, err = dcYml.osStruct.Stat(fPath)
		if err != nil {
			return dcYmlValue, err
		}
	}

	if fileInfo.IsFile() {
		dcYmlValue, err = dcYml.getServiceNames(fPath)
		if err != nil {
			err = fmt.Errorf(`get services from '%v': %w`, fPath, err)
			return
		}
	} else {
		if isDir {
			err = fmt.Errorf(`not a file: '%v'`, fPath)
		} else {
			err = fmt.Errorf(`not a dir or file: '%v'`, fPath)
		}
		return
	}

	if len(dcYmlValue.DcServicesNames) == 0 {
		return types.DcYmlValue{}, fmt.Errorf("no service names in config: '%v'", fPath)
	} else {
		for _, serviceName := range dcYmlValue.DcServicesNames {
			if len(serviceName) == 0 || !regexp.MustCompile(`^\w[\w\d]*$`).MatchString(serviceName) {
				return types.DcYmlValue{}, fmt.Errorf("empty or inappropriate service name '%v' in config: '%v'", serviceName, fPath)
			}
		}
	}

	return
}

func (dcYml *DcYml) Get() (dcYmlValue types.DcYmlValue, err error) {
	fPath := dcYml.config.GetDcYmlFname()

	if len(fPath) > 0 {
		return dcYml.getByFname(fPath)
	} else {
		workDir, errGetwd := dcYml.osStruct.Getwd()
		if errGetwd != nil {
			errGetwd = fmt.Errorf(`getting current working directory: '%w'`, errGetwd)
			return dcYmlValue, errGetwd
		}
		return dcYml.getByFname(workDir)
	}
}

// type DcConfig struct {
// 	OsStruct DcConfigOsInterface
// 	Fqfn     string
// }

// type DcConfigValueType = struct {
// 	Workdir         string
// 	DcServicesNames []string
// } // map[string]interface{}

// type DcConfigFileType = struct {
// 	Services map[string]map[string]interface{}
// }

func getServices(mapSlices yaml.MapSlice) ([]string, error) {
	services := []string{}

	for _, slice := range mapSlices {
		if slice.Key == `services` {
			servicesFound, ok := slice.Value.(yaml.MapSlice)
			if !ok {
				return []string{}, fmt.Errorf("services are not a list of strings: '%v'", mapSlices)
			}

			for _, servicesItem := range servicesFound {
				serviceName, ok := servicesItem.Key.(string)
				if !ok {
					return []string{}, fmt.Errorf("service name is not a string: '%v'", servicesItem.Key)
				}
				services = append(services, serviceName)
			}

			continue
		}
	}

	// for _, mapSlice := range mapSlices {
	// 	dcKeyValue, ok := mapSlice.Key.(string)
	// }

	return services, nil
}

// func (dcConfig DcConfig) readConfigFile() (value DcConfigValueType, err error) {
// 	var buf []byte
// 	fqfn := dcConfig.Fqfn

// 	buf, err = dcConfig.OsStruct.ReadFile(fqfn)
// 	if err != nil {
// 		return DcConfigValueType{}, fmt.Errorf("reading config file: '%s' error:\n\t%w", fqfn, err)
// 	}

// 	mapSlices := yaml.MapSlice{}

// 	err = yaml.Unmarshal(buf, &mapSlices)
// 	if err != nil {
// 		return DcConfigValueType{}, fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
// 	}

// 	services, err := getServices(mapSlices)
// 	if err != nil {
// 		err = fmt.Errorf("parsing config file: '%s' error:\n\t%w", fqfn, err)
// 		return
// 	}

// 	value.DcServicesNames = services
// 	return
// }

// func (dcConfig DcConfig) Read() (DcConfigValueType, error) {
// 	dir := filepath.Dir(dcConfig.Fqfn)
// 	err := dcConfig.OsStruct.Chdir(dir)
// 	if err != nil {
// 		return DcConfigValueType{}, fmt.Errorf("failed to change to dir: '%v' error:\n\t%w", dir, err)
// 	}
// 	workDir, err := dcConfig.OsStruct.Getwd()
// 	if err != nil {
// 		return DcConfigValueType{}, fmt.Errorf("failed to get current directory:\n\t%w", err)
// 	}
// 	if len(workDir) == 0 {
// 		return DcConfigValueType{}, fmt.Errorf("work directory is empty")
// 	}

// 	value, err := dcConfig.readConfigFile()
// 	if err != nil {
// 		return value, err
// 	}

// 	if len(value.DcServicesNames) == 0 {
// 		return DcConfigValueType{}, fmt.Errorf("no service names in config: '%v'", dcConfig.Fqfn)
// 	} else {
// 		for _, serviceName := range value.DcServicesNames {
// 			if len(serviceName) == 0 || !regexp.MustCompile(`^\w[\w\d]*$`).MatchString(serviceName) {
// 				return DcConfigValueType{}, fmt.Errorf("empty or inappropriate service name '%v' in config: '%v'", serviceName, dcConfig.Fqfn)
// 			}
// 		}
// 	}

// 	value.Workdir = workDir

// 	return value, err
// }

// func (dcConfig DcConfig) SetFqfn(fqfn string) {}
