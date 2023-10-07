package receiver

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	keyValue = "value"
)

type Adapter struct {
	logger *zerolog.Logger
	value  float64
}

func (a *Adapter) Configure() error {
	return nil
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("dummy/receiver: Adapter created")
	return &Adapter{
		logger: logger,
		value:  0,
	}
}

func (a *Adapter) Cycle(*time.Time) {
}

func (a *Adapter) Output() map[string]float64 {
	a.logger.Warn().Msgf("dummy/receiver: Adapter does not send output")
	return nil
}

func (a *Adapter) Input(value float64, key string) {
	switch key {
	case keyValue:
		a.value = value
		a.logger.Info().Msgf("dummy/receiver: Received value %f", a.value)
	default:
		a.logger.Warn().Msgf("dummy/receiver: Adapter does not accept input for key %s", key)
	}
}
