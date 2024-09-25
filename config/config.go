package config

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/cristalhq/aconfig"
	"github.com/go-playground/validator/v10"
)

const project = "bcdb"

type Config struct {
	Debug bool `default:"false"`

	Mode string `default:"server" validate:"required,oneof=server client"`

	Server Server
}
type Server struct {
	Address    string `default:"127.0.0.1" validate:"omitempty,ip"`
	Port       int    `default:"8080" validate:"required,numeric,gt=0,lt=65536"`
	MaxClients int    `default:"10" validate:"required,gt=0"`
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
		Files: files,
	})

	cfgErr := loader.Load()
	if cfgErr != nil {
		return nil, cfgErr
	}

	return &c, c.validate()
}

func (c *Config) validate() error {
	validate := validator.New()
	vErr := validate.Struct(c)
	if vErr != nil {
		var errs validator.ValidationErrors
		if errors.As(vErr, &errs) {
			return errs
		}
		return vErr
	}

	return nil
}
