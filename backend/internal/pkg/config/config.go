package config

import "os"

type Config struct {
	Common   CommonConfig
	Rest     RestConfig
	DB       DBConfig
	ChainAPI ChainAPIConfig
}

type CommonConfig struct {
	LogLevel     string
	LogLocal     bool
	LogFile      string
	LogAddSource bool

	JWTSecret []byte
}

type RestConfig struct {
	Address string
	TLS     bool
}

type DBConfig struct {
	// persistent database config
	Host      string
	EnableSSL bool
	Database  string
	User      string
	Secret    string

	// cache config
	CacheHost   string
	CacheUser   string
	CacheSecret string
}

type ChainAPIConfig struct {
	Host string
}

type QueuesConfig struct {
	ChainJobsQueue Queue
}

type Queue struct {
	Driver string
	Name   string
	Host   string
	User   string
	Secret string
}

func (c *Config) ReadFromEnv() {
	if c.Common.LogLevel == "" {
		c.Common.LogLevel = os.Getenv("BLOCKD_LOG_LEVEL")
	}

	if !c.Common.LogLocal {
		if os.Getenv("BLOCKD_LOG_LOCAL") == "true" {
			c.Common.LogLocal = true
		}
	}

	if c.Common.LogFile == "" {
		c.Common.LogFile = os.Getenv("BLOCKD_LOG_FILE")
	}

	if !c.Common.LogAddSource {
		if os.Getenv("BLOCKD_LOG_ADD_SOURCE") == "true" {
			c.Common.LogAddSource = true
		}
	}

	// os.Getenv("BLOCKD_JWT_SECRET")
	// os.Getenv("BLOCKD_CHAIN_API_URL")
	// os.Getenv("BLOCKD_NUM_INTERNAL_WORKERS")

	// os.Getenv("BLOCKD_REST_ADDRESS")
	// os.Getenv("BLOCKD_REST_ENABLE_TLS")
	// os.Getenv("BLOCKD_REST_CERT_PATH")
	// os.Getenv("BLOCKD_REST_KEY_PATH")

	// os.Getenv("BLOCKD_DB_HOST")
	// os.Getenv("BLOCKD_DB_DATABASE")
	// os.Getenv("BLOCKD_DB_USER")
	// os.Getenv("BLOCKD_DB_SECRET")
	// os.Getenv("BLOCKD_DB_ENABLE_TLS")

	// os.Getenv("BLOCKD_CACHE_HOST")
	// os.Getenv("BLOCKD_CACHE_USER")
	// os.Getenv("BLOCKD_CACHE_SECRET")
}
