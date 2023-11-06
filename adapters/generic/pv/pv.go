package pv

import (
	"fmt"
	"math"
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
	simulatedTime *time.Time
	iteration     uint

	p_w float64
}

type conf struct {
	Host string `json:"host"`
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("PV: Adapter created")
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

	a.logger.Info().Msg("PV: Adapter configured")

	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {

}

func (a *Adapter) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {

	// loop through all register addresses from req.addr to req.addr + req.Quantity - 1
	for regAddr := req.Addr; regAddr < req.Addr+req.Quantity; regAddr++ {
		switch regAddr {
		case 0:

			val, _ := utils.Uint32ToUint16(math.Float32bits(float32(a.p_w * 1000)))
			res = append(res, val)

		case 1:
			_, val := utils.Uint32ToUint16(math.Float32bits(float32(a.p_w * 1000)))
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

func (a *Adapter) Input(value any, key string) {
	switch key {
	case "p_w":
		_, ok := value.(float64)
		if !ok {
			a.logger.Warn().Msg("PV: invalid input type for p_w")
			return
		}
		a.p_w = value.(float64)
	default:
		a.logger.Warn().Msg("PV: invalid input key")
	}
}

func (a *Adapter) Output() map[string]any {
	return map[string]any{
		"p_w": a.p_w,
	}
}
