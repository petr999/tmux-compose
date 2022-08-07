package types

type DcYmlInterface interface {
	New(DcYmlOsInterface)
	Get() (DcYmlValue, error)
}

type DcFileInfoStruct struct {
	IsDir  func() bool
	IsFile func() bool
}
type DcYmlOsInterface interface {
	Getwd() (dir string, err error)
	ReadFile(name string) ([]byte, error)
	Getenv(string) string
	Stat(name string) (DcFileInfoStruct, error)
}

type DcYmlValue = struct {
	Workdir         string
	DcServicesNames []string
}
