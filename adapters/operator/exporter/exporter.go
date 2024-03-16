package exporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
)

type Adapter struct {
	file     *os.File
	conf     *conf
	logger   *zerolog.Logger
	confPath string
	writer   *csv.Writer

	result float64
}

type conf struct {
	Members []member `json:"members"`
	CsvPath string   `json:"csv_path"`
}

type member struct {
	value float64
	Key   string `json:"key"`
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("operator/exporter: Adapter created")

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
		a.logger.Fatal().Err(err).Msg("operator/exporter: Failed to read config file")
		return fmt.Errorf("operator/exporter: Failed to read config file: %w", err)
	}

	a.conf = conf

	if a.createFile() != nil {
		return fmt.Errorf("operator/exporter: Failed to create CSV file: %w", err)
	}

	members := make([]string, len(a.conf.Members))
	for i, member := range a.conf.Members {
		members[i] = member.Key
	}
	a.logger.Info().Msg("operator/exporter: Configured with members: " + fmt.Sprint(members))

	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {
	var (
		values = []string{simulatedTime.Format(time.RFC3339)}
		err    error
	)

	for _, member := range a.conf.Members {
		values = append(values, fmt.Sprintf("%f", member.value))
	}
	a.logger.Debug().Msgf("operator/exporter: Writing values %v", values)

	if err = a.writer.Write(values); err != nil {
		a.logger.Error().Err(err).Msg("operator/exporter: Failed to write values")
	}
	a.writer.Flush()

}

func (a *Adapter) Output() map[string]any {
	return map[string]any{
		"result": a.result,
	}
}

func (a *Adapter) Input(value any, key string) {
	for i, member := range a.conf.Members {
		if member.Key == key {
			v, ok := value.(float64)
			if !ok {
				a.logger.Error().Msgf("operator/exporter: Unexpected input type for key %s", value)
				return
			}
			a.conf.Members[i].value = v
			a.logger.Debug().Msgf("operator/exporter: Updated value for key %s: %f", key, v)
			return
		}
	}

	a.logger.Warn().Msgf("operator/exporter: Unexpected input key %s", key)

}

func (a *Adapter) createFile() error {
	var (
		err     error
		headers = make([]string, len(a.conf.Members)+1)
	)
	a.file, err = os.Create(a.conf.CsvPath)
	if err != nil {
		a.logger.Error().Err(err).Msg("operator/exporter: Failed to create file")
		return fmt.Errorf("operator/exporter: Failed to create file: %w", err)
	}

	a.writer = csv.NewWriter(a.file)

	headers[0] = "timestamp"
	for i, member := range a.conf.Members {
		headers[i+1] = member.Key
	}

	if err = a.writer.Write(headers); err != nil {
		a.logger.Error().Err(err).Msg("operator/exporter: Failed to write headers")
		return fmt.Errorf("operator/exporter: Failed to write headers: %w", err)
	}
	a.writer.Flush()
	return nil
}
