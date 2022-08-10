package types

type ConfigInterface interface {
	New(ConfigOsInterface)
	GetCnaTemplateFname() string
	GetDcYmlFname() string
}

type ConfigOsInterface interface {
	Getenv(string) string
}
