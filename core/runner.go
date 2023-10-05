package core

import (
	"reflect"
	"time"

	"github.com/kilianp07/simu/adapters"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Runner struct {
	conf          *Conf
	adapters      []adapters.Adapter
	simulatedTime *time.Time
	logger        *zerolog.Logger
	lastCycle     time.Time
}

func (r *Runner) readConfig(confPath string) (*Conf, error) {
	viper.SetConfigFile(confPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	conf := &Conf{}

	if err := viper.Unmarshal(conf, func(m *mapstructure.DecoderConfig) {
		m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			func(
				f reflect.Type,
				t reflect.Type,
				data interface{}) (interface{}, error) {
				if f.Kind() != reflect.String {
					return data, nil
				}
				if t != reflect.TypeOf(time.Time{}) {
					return data, nil
				}

				asString := data.(string)
				if asString == "" {
					return time.Time{}, nil
				}

				// Convert it by parsing
				return time.Parse("2006-01-02", asString)
			},
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	}); err != nil {
		return nil, err
	}

	if conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return conf, nil
}

func (r *Runner) instanciate() {
	for _, adapter := range r.conf.Adapters {
		if a := adapters.New(adapter.Adapter, adapter.ConfPath, r.simulatedTime, r.logger); a != nil {
			r.adapters = append(r.adapters, *a)
		}
	}
}

func (r *Runner) configureAdapters() error {
	for _, adapter := range r.adapters {
		if err := adapter.Configure(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) run() {
	for {
		// Wait for the next cycle
		if time.Since(r.lastCycle) < time.Duration(r.conf.Timestep)*time.Millisecond {
			continue
		}

		// At each cycle, we increment the simulated time
		t := r.simulatedTime.Add(time.Duration(r.conf.Period) * time.Millisecond)
		r.simulatedTime = &t

		for _, adapter := range r.adapters {
			adapter.Cycle(r.simulatedTime)
		}

		r.lastCycle = time.Now()
	}
}
