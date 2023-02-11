package config

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}
