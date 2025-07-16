package configs

type Configuration struct {
	API      *APIConfig
	Database *DatabaseConfig
	TC       *TCConfig
}

type APIConfig struct {
	Host      string
	Port      uint32
	Env       string
	LogPath   string
	LogPrefix string
	AccessKey string
	PIDPath   string
	PIDFile   string
}

type DatabaseConfig struct {
	Path string
	File string
}

type TCConfig struct {
	InterfaceName  string
	HTBVersion     uint32 // HTB version (default: 3)
	Rate2Quantum   uint32 // Rate to quantum ratio (default: 10)
	DefaultClassID uint32 // Default class ID (default: 1)
	EnableLogging  bool   // Enable TC logging
	LogLevel       string // Log level for TC operations
}
