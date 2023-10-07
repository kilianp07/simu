package sum

import (
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
)

const (
	keyResult = "result"
)

type conf struct {
	Members []member `json:"members"`
}

type member struct {
	value  float64
	Key    string  `json:"key"`
	Weight float64 `json:"weight"`
}

type Adapter struct {
	conf     *conf
	logger   *zerolog.Logger
	confPath string

	result float64
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("operator/sum: Adapter created")

	return &Adapter{
		conf:     &conf{},
		logger:   logger,
		confPath: confpath,
	}
}

func (a *Adapter) Configure() error {
	var (
		conf = &conf{}
		err  error
	)

	if err = utils.ReadJsonFile(a.confPath, conf); err != nil {
		a.logger.Fatal().Err(err).Msg("operator/sum: Failed to read config file")
		return err
	}

	a.conf = conf

	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {
	var (
		result float64
	)

	for _, member := range a.conf.Members {
		result += member.value * member.Weight
	}

	a.result = result
}

func (a *Adapter) Output() map[string]float64 {
	return map[string]float64{
		"result": a.result,
	}
}

func (a *Adapter) Input(value float64, key string) {
	for i, member := range a.conf.Members {
		if member.Key == key {
			a.conf.Members[i].value = value
			return
		}
	}

	a.logger.Warn().Msgf("operator/sum: Unexpected input key %s", key)

}
