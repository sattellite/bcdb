package config

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/cristalhq/aconfig"
)

const project = "bcdb"

type Config struct {
	Debug bool
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

	return &c, nil
}
