package types

type DcYmlInterface interface {
	New(DcYmlOsInterface)
	Get() (DcYmlValue, error)
}

type DcYmlOsInterface interface {
	Getwd() (dir string, err error)
	ReadFile(name string) ([]byte, error)
	Getenv(string) string
	Stat(name string) (FileInfoStruct, error)
}

type DcYmlValue = struct {
	Workdir         string
	DcServicesNames []string
}
