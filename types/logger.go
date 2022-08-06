package types

type LoggerInterface interface {
	New(stdHandles StdHandlesType)
	Log(s string)
}
