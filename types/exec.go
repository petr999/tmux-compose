package types

type ExecInterface interface {
	// MakeCommand(*MakeCommandDryRunType,
	// 	nameArgsType,
	// )
	New(ExecOsInterface, StdHandlesType)
	GetCommand(cna CmdNameArgsValueType) *CmdType
	// execCommand(cna CmdNameArgsValueType) CmdInterface
	// dryRun(cna CmdNameArgsValueType) CmdInterface
}

type CmdInterface interface {
	Run() error
}

type CmdType struct {
	Obj        CmdInterface
	Stdhandles StdHandlesType
}
