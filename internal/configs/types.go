package configs

type Configuration struct {
	API      *APIConfig
	Database *DatabaseConfig
}

type APIConfig struct {
	Host      string
	Port      int
	Env       string
	LogPath   string
	LogPrefix string
}

type DatabaseConfig struct {
	Path string
}
