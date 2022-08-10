package types

type CnaInterface interface {
	New(CnaOsInterface, ConfigInterface)
	Get(DcYmlValue) (CmdNameArgsValueType, error)
}

type CnaOsInterface interface {
	ReadFile(name string) ([]byte, error)
}

type CmdNameArgsValueType struct {
	Workdir string
	Name    string
	Args    []string
}
