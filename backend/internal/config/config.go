package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"MB_ENV" env-default:"local"` // local|dev|prod
	StoragePath string `yaml:"storage_path" env:"MB_STORAGE_PATH" env-required:"true"`
	Logger      `yaml:"logger"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Port      int           `yaml:"port" env:"MB_HTTP_PORT" env-default:"8080"`
	Timeout   time.Duration `yaml:"timeout" env:"MB_HTTP_TIMEOUT" env-default:"5s"`
	TokenTTL  time.Duration `yaml:"token_ttl" env:"MB_TOKEN_TTL" env-default:"24h"`
	JWTSecret string        `yaml:"jwt_secret" env:"MB_AUTH_SECRET" env-required:"true"`
}

type Logger struct {
	LevelString string     `yaml:"level" env:"MB_LOG_LEVEL" env-default:"info"`
	Level       slog.Level `yaml:"-"` // will be loaded later
}

func (c *Logger) InitLevel() {
	switch c.LevelString {
	case "debug":
		c.Level = slog.LevelDebug
	case "warn":
		c.Level = slog.LevelWarn
	case "error":
		c.Level = slog.LevelError
	default:
		c.Level = slog.LevelInfo
	}
}

// MustLoad returns Config of just PANIC
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file path is empty")
	}

	return MustLoadByPath(path)
}

// MustLoad returns Config of just PANIC
func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or env variable
// Priority: flag > env > default.
// Default = ""
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
