package types

type CnaInterface interface {
	New(CnaOsInterface, ConfigInterface)
	Get(DcYmlValue) (CmdNameArgsValueType, error)
}

type CnaOsInterface interface {
	ReadFile(name string) ([]byte, error)
	Stat(name string) (FileInfoStruct, error)
}

type CmdNameArgsValueType struct {
	Workdir string
	Cmd     string
	Args    []string
}
