package poc

import (
	"fmt"
	"math"
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
	"github.com/simonvetter/modbus"
)

const (
	keyPw = "p_w"
)

type conf struct {
	Host string `json:"host"`
}

type Adapter struct {
	server   *modbus.ModbusServer
	conf     *conf
	logger   *zerolog.Logger
	confPath string

	p_w float32
}

func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("generic/poc: Adapter created")

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
		a.logger.Fatal().Err(err).Msg("generic/poc: Failed to read config file")
		return err
	}

	a.conf = conf

	a.server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL: fmt.Sprintf("tcp://%s", a.conf.Host),
		// close idle connections after 30s of inactivity
		Timeout: 30 * time.Second,
		// accept 5 concurrent connections max.
		MaxClients: 5,
	}, a)

	if err != nil {
		a.logger.Fatal().Err(err).Msg("generic/poc: Failed to create modbus server")
		return err
	}

	if err = a.server.Start(); err != nil {
		a.logger.Fatal().Err(err).Msg("generic/poc: failed to start modbus server")
		return err
	}

	a.logger.Info().Msgf("generic/poc: Modbus server started on %s", a.conf.Host)

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
			err = modbus.ErrIllegalDataAddress
			return
		}
	}

	return res, nil
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

func (a *Adapter) Output() map[string]any {
	return map[string]any{
		"p_w": a.p_w,
	}
}

func (a *Adapter) Input(value any, key string) {
	switch key {
	case keyPw:
		if _, ok := value.(float64); !ok {
			a.logger.Warn().Msgf("generic/poc: Adapter input %s is not a float64", key)
			return
		}
		a.p_w = float32(value.(float64))
	default:
		a.logger.Warn().Msgf("generic/poc: Adapter does not have input %s", key)
	}
}
