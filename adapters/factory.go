package adapters

import (
	"time"

	"github.com/kilianp07/simu/adapters/dummy/receiver"
	"github.com/kilianp07/simu/adapters/dummy/sender"
	"github.com/kilianp07/simu/adapters/generic/battery"
	"github.com/kilianp07/simu/adapters/generic/poc"
	"github.com/kilianp07/simu/adapters/generic/pv"
	csvreader "github.com/kilianp07/simu/adapters/operator/csvReader"
	"github.com/kilianp07/simu/adapters/operator/sum"
	"github.com/rs/zerolog"
)

type Adapter interface {
	Configure() error
	Cycle(*time.Time)
	Output() map[string]any
	Input(value any, key string)
}

func New(name string, confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	var adapter Adapter

	switch name {
	case "generic/pv":
		adapter = pv.New(confpath, simulatedTime, logger)
	case "generic/battery":
		adapter = battery.New(confpath, simulatedTime, logger)
	case "generic/poc":
		adapter = poc.New(confpath, simulatedTime, logger)
	case "operator/reader":
		adapter = csvreader.New(confpath, simulatedTime, logger)
	case "operator/sum":
		adapter = sum.New(confpath, simulatedTime, logger)
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
