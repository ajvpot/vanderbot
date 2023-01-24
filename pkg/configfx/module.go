package configfx

import (
	"fmt"
	"os"
	"path"

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

func envLookup(v string) (string, bool) {
	return os.LookupEnv(v)
}

func New(p Params) (Result, error) {
	var opts = []config.YAMLOption{
		config.Expand(envLookup),
	}

	opts = append(opts, tryFiles(".env")...)

	opts = append(opts, tryFiles("config")...)
	opts = append(opts, tryFiles("secrets")...)

	if f := os.Getenv("CONFIG_FILE"); f != "" {
		opts = append(opts, config.File(f))
	}

	provider, err := config.NewYAML(opts...)

	return Result{Provider: provider}, err
}

func tryFiles(name string) (out []config.YAMLOption) {
	if o := tryFile(fmt.Sprintf("%s.yml", name)); o != nil {
		out = append(out, o)
	}
	if o := tryFile(fmt.Sprintf("%s.yaml", name)); o != nil {
		out = append(out, o)
	}
	if kdp := os.Getenv("KO_DATA_PATH"); kdp != "" {
		if o := tryFile(path.Join(kdp, fmt.Sprintf("%s.yaml", name))); o != nil {
			out = append(out, o)
		}
		if o := tryFile(path.Join(kdp, fmt.Sprintf("%s.yml", name))); o != nil {
			out = append(out, o)
		}
	}
	return out
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
