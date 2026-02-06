package config

var globalConfig *Config

type Config struct {
	API    APIConfig
	Git    GitConfig
	Sync   SyncConfig
	MCP    MCPConfig
	Logging LoggingConfig
	Cache  CacheConfig
}

type APIConfig struct {
	URL     string
	Key     string
	Timeout int
}

type GitConfig struct {
	Repository string
	Branch    string
	User      string
	Email     string
}

type SyncConfig struct {
	AutoPush  bool
	AutoPull  bool
	Interval  int
}

type MCPConfig struct {
	Enabled  bool
	CacheDir string
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

type CacheConfig struct {
	Dir string
}

func Load(configPath string) (*Config, error) {
	return &Config{}, nil
}

func Get() *Config {
	return globalConfig
}

func SetGlobal(cfg *Config) {
	globalConfig = cfg
}
