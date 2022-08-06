package types

type ExecInterface interface {
	// MakeCommand(*MakeCommandDryRunType,
	// 	nameArgsType,
	// )
	New(ExecOsInterface, StdHandlesType)
	GetCommand(CmdNameArgsValueType) *CmdType
}

type CmdInterface struct {
	Run func() error
}

type CmdType struct {
	Obj interface {
		Run() error
	}
	Stdhandles StdHandlesType
}
