package types

type ConfigInterface interface {
	New(ConfigOsInterface)
	GetCnaTemplateFname() string
	GetDcYmlFname() string
	GetDryRun() string
}

type ConfigOsInterface interface {
	Getenv(string) string
}
