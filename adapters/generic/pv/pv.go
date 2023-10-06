package pv

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
	"github.com/simonvetter/modbus"
)

type Adapter struct {
	server        *modbus.ModbusServer
	Conf          *conf
	confPath      string
	logger        *zerolog.Logger
	data          [][]string
	simulatedTime *time.Time
	iteration     uint

	p_w float64
}

type conf struct {
	StartDate  string `json:"start_date"`
	TimeFormat string `json:"time_format"`
	CsvPath    string `json:"csv_path"`
	P_col      uint   `json:"p_col"`
	Date_col   uint   `json:"date_col"`
	Host       string `json:"host"`
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {

	a := &Adapter{
		Conf:          &conf{},
		logger:        logger,
		simulatedTime: simulatedTime,
		confPath:      confpath,
		iteration:     1,
	}

	return a
}
func (a *Adapter) Configure() error {
	var (
		err error
	)
	conf := &conf{}

	if err = utils.ReadJsonFile(a.confPath, &conf); err != nil {
		a.logger.Fatal().Err(err).Msg("PV: failed to read config file")
		return err
	}
	a.Conf = conf

	a.server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL: fmt.Sprintf("tcp://%s", a.Conf.Host),
		// close idle connections after 30s of inactivity
		Timeout: 30 * time.Second,
		// accept 5 concurrent connections max.
		MaxClients: 5,
	}, a)

	if err != nil {
		a.logger.Fatal().Err(err).Msg("PV: failed to create modbus server")
		return err
	}

	if err = a.server.Start(); err != nil {
		a.logger.Fatal().Err(err).Msg("PV: failed to start modbus server")
		return err
	}

	if err := a.getCsvData(); err != nil {
		a.logger.Fatal().Err(err).Msg("PV: failed to get csv data")
		return err
	}

	return nil
}

func (a *Adapter) getCsvData() error {
	var (
		data [][]string
		err  error
	)

	if data, err = utils.ReadCsvFile(a.Conf.CsvPath); err != nil {
		a.logger.Fatal().Err(err).Msg("PV: failed to read csv file")
		return err
	}

	a.data = data

	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {

	delta := simulatedTime.Sub(*a.simulatedTime)

	nextDate, err := time.Parse(a.Conf.TimeFormat, a.data[a.iteration+1][a.Conf.Date_col])
	if err != nil {
		a.logger.Err(err).Msg("PV: failed to parse date")
	}

	actualDate, err := time.Parse(a.Conf.TimeFormat, a.data[a.iteration][a.Conf.Date_col])
	if err != nil {
		a.logger.Err(err).Msg("PV: failed to parse date")
	}

	if nextDate.After(actualDate.Add(delta)) {
		a.iteration++
	}

	a.p_w, err = strconv.ParseFloat(a.data[a.iteration][a.Conf.P_col], 64)
	if err != nil {
		a.logger.Err(err).Msg("PV: failed to parse power")
	}

}
func (a *Adapter) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {

	// loop through all register addresses from req.addr to req.addr + req.Quantity - 1
	for regAddr := req.Addr; regAddr < req.Addr+req.Quantity; regAddr++ {
		switch regAddr {
		case 0:
			val, _ := utils.Uint32ToUint16(uint32(a.p_w * 100))
			res = append(res, val)

		case 1:
			_, val := utils.Uint32ToUint16(uint32(a.p_w * 100))
			res = append(res, val)

		default:
			a.logger.Warn().Msg("PV: illegal data address")
			err = modbus.ErrIllegalDataAddress
			return
		}
	}

	return res, err
}

func (a *Adapter) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	return nil, modbus.ErrIllegalFunction
}

func (a *Adapter) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	return nil, modbus.ErrIllegalFunction
}

func (a *Adapter) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	return nil, modbus.ErrIllegalFunction
}
