package csvreader

import (
	"strconv"
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
)

type Adapter struct {
	conf     *conf
	confPath string
	logger   *zerolog.Logger
	data     [][]string

	simulatedTime *time.Time
	value         float64
	iteration     uint
	delta         time.Duration
}

type conf struct {
	StartDate  string `json:"start_date"`
	TimeFormat string `json:"time_format"`
	CsvPath    string `json:"csv_path"`
	DataCol    uint   `json:"data_column"`
	DateCol    uint   `json:"date_column"`
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("Reader: Adapter created")
	a := &Adapter{
		conf:          &conf{},
		logger:        logger,
		simulatedTime: simulatedTime,
		confPath:      confpath,
		iteration:     1,
	}

	return a
}

func (a *Adapter) Configure() error {
	var (
		err  error
		conf = &conf{}
	)

	if err = utils.ReadJsonFile(a.confPath, &conf); err != nil {
		a.logger.Fatal().Err(err).Msg("Reader: failed to read config file")
		return err
	}
	a.conf = conf

	if err := a.getData(); err != nil {
		a.logger.Fatal().Err(err).Msg("Reader: failed to get csv data")
		return err
	}

	a.logger.Info().Msg("Reader: configured")
	return nil
}

func (a *Adapter) getData() error {
	var (
		data [][]string
		err  error
	)

	if data, err = utils.ReadCsvFile(a.conf.CsvPath); err != nil {
		a.logger.Fatal().Err(err).Msg("Reader: failed to read csv file")
		return err
	}

	a.data = data

	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {
	defer a.updateTime(simulatedTime)
	if a.simulatedTime != nil {
		deltaTime := simulatedTime.Sub(*a.simulatedTime)
		a.delta += deltaTime
	} else {
		a.logger.Warn().Msg("Previous simulated time is nil, skipping delta calculation")
	}

	if a.iteration+1 < uint(len(a.data)) {

		nextDate, err := time.Parse(a.conf.TimeFormat, a.data[a.iteration+1][int(a.conf.DateCol)])
		if err != nil {
			a.logger.Err(err).Msg("Reader: failed to parse next date")
			return
		}

		currentDate, err := time.Parse(a.conf.TimeFormat, a.data[a.iteration][int(a.conf.DateCol)])
		if err != nil {
			a.logger.Err(err).Msg("Reader: failed to parse current date")
			return
		}

		expectedDelta := nextDate.Sub(currentDate)
		if a.delta >= expectedDelta {
			a.iteration++
			a.delta -= expectedDelta
			a.logger.Info().Msgf("Reader: Advancing to iteration %d", a.iteration)
		}
	}
	// Mise à jour de la valeur courante
	if value, err := strconv.ParseFloat(a.data[a.iteration][int(a.conf.DataCol)], 64); err == nil {
		a.value = value
	} else {
		a.logger.Err(err).Msg("Reader: Failed to parse current value")
	}
}

func (a *Adapter) Input(value any, key string) {
	a.logger.Warn().Msg("Reader: Adapter has no input")
}

func (a *Adapter) Output() map[string]any {
	return map[string]any{
		"value": a.value,
	}
}

func (a *Adapter) updateTime(t *time.Time) {
	a.simulatedTime = t
}
