package configfx

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"
)

// Module provides a config.Provider.
var Module = fx.Provide(New)

// Params defines the dependencies of the configfx module.
type Params struct {
	fx.In
}

// Result defines the objects that the configfx module provides.
type Result struct {
	fx.Out

	Provider config.Provider
}

func New(p Params) (Result, error) {
	var opts []config.YAMLOption

	if o := tryFile("config.yml"); o != nil {
		opts = append(opts, o)
	}
	if o := tryFile("config.yaml"); o != nil {
		opts = append(opts, o)
	}

	if o := tryFile("secrets.yml"); o != nil {
		opts = append(opts, o)
	}
	if o := tryFile("secrets.yaml"); o != nil {
		opts = append(opts, o)
	}

	if o := tryFile(".env"); o != nil {
		opts = append(opts, o)
	}

	if f := os.Getenv("CONFIG_FILE"); f != "" {
		opts = append(opts, config.File(f))
	}

	provider, err := config.NewYAML(opts...)

	return Result{Provider: provider}, err
}

func tryFile(path string) config.YAMLOption {
	fi, err := os.Stat(path)
	if err == nil {
		if !fi.IsDir() {
			return config.File(path)
		}
	}
	return nil
}
