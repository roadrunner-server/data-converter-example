package dataconverter

import (
	"time"

	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"go.temporal.io/sdk/converter"
	"go.uber.org/zap"
)

// plugin name
// key in the configuration file should be with this name
const name = "encrypted-data-converter"

type Configurer interface {
	// UnmarshalKey takes a single key and unmarshals it into a Struct.
	UnmarshalKey(name string, out any) error

	// Unmarshal the config into a Struct. Make sure that the tags
	// on the fields of the structure are properly set.
	Unmarshal(out any) error

	// Has checks if config section exists.
	Has(name string) bool

	// GracefulTimeout represents timeout for all servers registered in the endure
	GracefulTimeout() time.Duration

	// RRVersion returns running RR version
	RRVersion() string
}

// Plugin structure should have exactly the `Plugin` name to be found by RR
type Plugin struct {
	log *zap.Logger
	cfg *Config
}

// Init will be called only once
func (p *Plugin) Init(cfg Configurer, log *zap.Logger) error {
	const op = errors.Op("nil-data-converter-init")
	// check for the `nil-data-converter-init` key in the configuration
	// you may skip this check if your plugin does not require any configuration
	if !cfg.Has(name) {
		// special type of error, which tells RR to disable this plugin
		return errors.E(errors.Disabled)
	}

	// initialize your configuration
	p.cfg = &Config{}
	err := cfg.UnmarshalKey(name, p.cfg)
	if err != nil {
		return errors.E(op, err)
	}
	// after unmarshal, init values which are not set with the default values
	err = p.cfg.InitDefaults()
	if err != nil {
		return errors.E(op, err)
	}

	// init logger
	p.log = new(zap.Logger)
	*p.log = *log

	return nil
}

func (p *Plugin) Name() string {
	return name
}

func (p *Plugin) PayloadConverter() converter.PayloadConverter {
	return NewJSONPayloadConverter()
}

func (p *Plugin) Provides() []*dep.Out {
	return []*dep.Out{
		dep.Bind((*converter.PayloadConverter)(nil), p.PayloadConverter()),
	}
}
