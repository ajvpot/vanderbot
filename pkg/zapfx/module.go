package zapfx

import (
	"go.uber.org/fx"
)

const (
	// ConfigurationKey is the portion of the configuration that this package reads.
	ConfigurationKey = "logging"
)

// Module provides a zap logger for structured logging.
//
// In YAML, logging configuration might look like this:
//
//	logging:
//	  level: info
//	  development: false
//	  sampling:
//	    initial: 100
//	    thereafter: 100
//	  encoding: json
var Module = fx.Options(fx.Provide(NewZap, NewMortarProxy))
