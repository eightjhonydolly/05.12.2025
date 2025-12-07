package config

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

func LoadConfig(configPath string) (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}, nil
}
