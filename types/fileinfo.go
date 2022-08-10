package types

type FileInfoStruct struct {
	IsDir  func() bool
	IsFile func() bool
}
