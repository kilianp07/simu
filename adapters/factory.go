package adapters

import (
	"time"

	"github.com/kilianp07/simu/adapters/dummy/receiver"
	"github.com/kilianp07/simu/adapters/dummy/sender"
	"github.com/kilianp07/simu/adapters/generic/battery"
	"github.com/kilianp07/simu/adapters/generic/pv"
	"github.com/rs/zerolog"
)

type Adapter interface {
	Configure() error
	Cycle(*time.Time)
	Output() map[string]float64
	Input(value float64, key string)
}

func New(name string, confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	var adapter Adapter

	switch name {
	case "generic/pv":
		adapter = pv.New(confpath, simulatedTime, logger)
	case "generic/battery":
		adapter = battery.New(confpath, simulatedTime, logger)
	case "dummy/sender":
		adapter = sender.New(confpath, simulatedTime, logger)
	case "dummy/receiver":
		adapter = receiver.New(confpath, simulatedTime, logger)
	default:
		logger.Warn().Msgf("Adapter %s not found", name)
		return nil
	}

	return &adapter
}
