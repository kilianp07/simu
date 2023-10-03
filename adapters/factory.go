package adapters

import (
	"time"

	"github.com/kilianp07/simu/adapters/generic/pv"
	"github.com/rs/zerolog"
)

type Adapter interface {
	Configure() error
	Cycle(*time.Time)
}

func New(name string, confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	var adapter Adapter

	switch name {
	case "generic/pv":
		adapter = pv.New(confpath, simulatedTime, logger)
	default:
		logger.Warn().Msgf("Adapter %s not found", name)
		return nil
	}

	return &adapter
}
