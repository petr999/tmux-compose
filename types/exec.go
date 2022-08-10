package types

type ExecInterface interface {
	// MakeCommand(*MakeCommandDryRunType,
	// 	nameArgsType,
	// )
	New(ExecOsInterface, StdHandlesType, ConfigInterface)
	GetCommand(cna CmdNameArgsValueType) *CmdType
	GetSelector() any
	// execCommand(cna CmdNameArgsValueType) CmdInterface
	// dryRun(cna CmdNameArgsValueType) CmdInterface
}

type CmdInterface interface {
	Run() error
}

type CmdType struct {
	Obj        CmdInterface
	Stdhandles StdHandlesType
	Run        func() error
}
