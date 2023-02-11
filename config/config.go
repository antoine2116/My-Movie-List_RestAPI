package config

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	URI string
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}
