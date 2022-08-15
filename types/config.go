package types

type ConfigInterface interface {
	New(ConfigOsInterface)
	GetCnaTemplateFname() string
	GetDcYmlFname() string
	GetDryRun() string
	GetShell() string
}

type ConfigOsInterface interface {
	Getenv(string) string
}
