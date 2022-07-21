package dc_config

type Reader interface {
	Read() DcConfigValueType
}

type DcConfig struct{}

type DcConfigValueType = map[string]string

var value = make(DcConfigValueType)

func (dcConfig DcConfig) Read() DcConfigValueType {
	return value
}
