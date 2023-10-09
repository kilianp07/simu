package sender

import (
	"math/rand"
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
	logger.Info().Msg("dummy/sender: Adapter created")
	return &Adapter{
		logger: logger,
		value:  0,
	}
}

func (a *Adapter) Cycle(*time.Time) {
	a.value = rand.Float64()
}

func (a *Adapter) Output() map[string]any {
	a.logger.Debug().Msgf("dummy/sender: Sending value %f", a.value)
	return map[string]any{
		keyValue: a.value,
	}
}

func (a *Adapter) Input(value any, key string) {
	a.logger.Warn().Msg("dummy/sender: Adapter does not accept input")
}
