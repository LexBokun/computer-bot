package config

import (
	"errors"
	"fmt"
	"github.com/LexBokun/ControlAgent/internal/delivery/telegram"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"

	multimonitor "github.com/LexBokun/ControlAgent/internal/infrastructure/multi-monitor"
	"github.com/LexBokun/ControlAgent/internal/pkg/validator"

	"github.com/spf13/viper"
)

type Config struct {
	MonitorTool multimonitor.Config
	Logger      Logger
	CORS        CORS
	HTTPServer  HTTPServer
	Telegram    telegram.Config
}

type Logger struct {
	Level slog.Level
}
type CORS struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type HTTPServer struct {
	Port int `validate:"required"`
}

//nolint:gochecknoglobals // Singleton
var (
	once   sync.Once
	config *Config
)

// getGoModRoot returns the absolute path to the root directory containing go.mod.
func getGoModRoot() (string, error) {
	output, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}}").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func loadConfig(path string) (*viper.Viper, error) {
	v := viper.New()

	configRAW, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't load config file %s: %w", path, err)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	var configFileNotFoundErr viper.ConfigFileNotFoundError
	if err = v.ReadConfig(configRAW); err != nil {
		if errors.As(err, &configFileNotFoundErr) {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	return v, nil
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		slog.Error("unable to decode into struct", "error", err.Error())
		return nil, err
	}

	err = validator.GetInstance().Struct(&c)
	if err != nil {
		return nil, fmt.Errorf("can't validate config: %s", err.Error())
	}

	return &c, nil
}

func GetInstance() *Config {
	if config == nil {
		once.Do(
			func() {
				var err error

				path := "./config/config.docker.yaml"

				root, err := getGoModRoot()
				if err == nil {
					path = root + "/config/config.local.yaml"
				}

				viperCfg, err := loadConfig(path)
				if err != nil {
					slog.Error("error loading config file:", "error", err.Error())
					return
				}

				cfg, err := parseConfig(viperCfg)
				if err != nil {
					slog.Error("error parsing config file:", "error", err.Error())
					return
				}

				config = cfg
			},
		)
	}

	return config
}
