package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"

	"github.com/cristalhq/aconfig"
)

const project = "bcdb"

type Config struct {
	Debug bool

	Server struct {
		Address    string `default:"127.0.0.1"`
		Port       string `default:"8080"`
		MaxClients int    `default:"10"`
	}
}

func Load() (*Config, error) {
	var c Config
	// get config directory
	cfgPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	files := []string{
		"./config.toml",
		filepath.Join(cfgPath, project, "config.toml"),
	}

	switch runtime.GOOS {
	case "darwin":
		homePath, err := os.UserHomeDir()
		if err == nil {
			files = append(files, filepath.Join(homePath, ".config", project, "config.toml"))
		}
	case "windows":
		homePath, err := os.UserHomeDir()
		if err == nil {
			files = append(files, filepath.Join(homePath, "AppData", "Roaming", project, "config.toml"))
		}
	case "linux":
		files = append(files, "/etc/"+project+"/config.toml")
	}

	// remove duplicates
	files = slices.Compact(files)

	loader := aconfig.LoaderFor(&c, aconfig.Config{
		SkipFlags: true,
		Files:     files,
	})

	cfgErr := loader.Load()
	if cfgErr != nil {
		return nil, cfgErr
	}

	return &c, c.validate()
}

func (c *Config) validate() error {
	if c.Server.Port == "" {
		return errors.New("server port is required")
	} else {
		port, err := strconv.ParseInt(c.Server.Port, 10, 64)
		if err != nil {
			return errors.New("server port must be a number")
		}
		if port < 1 || port > 65535 {
			return errors.New("server port must be between 1 and 65535")
		}
	}

	if c.Server.MaxClients < 1 {
		return errors.New("server max clients must be greater than 0")
	}

	return nil
}
