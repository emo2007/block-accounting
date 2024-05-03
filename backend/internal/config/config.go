package config

type Config struct {
	Common CommonConfig
	Rest   RestConfig
	DB     DBConfig
	Eth    EthConfig
}

type CommonConfig struct {
	LogLevel     string
	LogLocal     bool
	LogFile      string
	LogAddSource bool
}

type RestConfig struct {
	Address string
	TLS     bool
}

type DBConfig struct {
	Host      string
	EnableSSL bool
	Database  string
	User      string
	Secret    string
}

type EthConfig struct {
	// todo
}
