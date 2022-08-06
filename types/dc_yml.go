package types

type DcYmlInterface interface {
	New(DcYmlOsInterface)
	Get() (DcYmlValue, error)
}

type DcYmlOsInterface interface {
	Chdir(dir string) error
	Getwd() (dir string, err error)
	ReadFile(name string) ([]byte, error)
}

type DcYmlValue = struct {
	Workdir         string
	DcServicesNames []string
}
