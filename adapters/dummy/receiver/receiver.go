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

func (a *Adapter) Output() map[string]any {
	a.logger.Warn().Msgf("dummy/receiver: Adapter does not send output")
	return nil
}

func (a *Adapter) Input(value any, key string) {
	switch key {
	case keyValue:
		if _, ok := value.(float64); !ok {
			a.logger.Error().Msgf("dummy/receiver: Unexpected input type for key %s", value)
			return
		}
		a.value = value.(float64)
		a.logger.Info().Msgf("dummy/receiver: Received value %f", a.value)
	default:
		a.logger.Warn().Msgf("dummy/receiver: Adapter does not accept input for key %s", key)
	}
}
