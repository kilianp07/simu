package battery

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
	"github.com/simonvetter/modbus"
)

type Adapter struct {
	server        *modbus.ModbusServer
	conf          *Conf
	confPath      string
	logger        *zerolog.Logger
	simulatedTime *time.Time

	p_w         float64
	soc         uint
	soh         uint
	capacity_Wh float64

	setPoint_w         float64
	setPoint1          uint16
	setPoint2          uint16
	lock               sync.RWMutex
	filter             *utils.LowPassFilter
	integraleEnergy_wh float64
	initialEnergy_wh   float64
}

type Conf struct {
	Soc          uint    `json:"soc"`
	Soh          uint    `json:"soh"`
	Capacity_Wh  float64 `json:"capacity_wh"`
	PCharge_w    float64 `json:"p_charge_w"`
	PDischarge_w float64 `json:"p_discharge_w"`
	Host         string  `json:"host"`
	Attenuation  float64 `json:"attenuation"`
}

func New(confPath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter {
	logger.Info().Msg("Battery: Adapter created")
	a := &Adapter{
		conf:          &Conf{},
		logger:        logger,
		simulatedTime: simulatedTime,
		confPath:      confPath,
		filter:        &utils.LowPassFilter{},
	}
	return a
}

func (a *Adapter) Configure() error {
	var err error
	conf := &Conf{}

	if err = utils.ReadJsonFile(a.confPath, conf); err != nil {
		a.logger.Error().Err(err).Msg("Battery: failed to read config file")
		return err
	}

	a.conf = conf

	// Init variables
	a.soc = a.conf.Soc
	a.soh = a.conf.Soh
	a.capacity_Wh = a.conf.Capacity_Wh * float64(a.soh) / 100
	a.initialEnergy_wh = a.capacity_Wh * float64(a.soc) / 100
	a.filter = utils.NewLowPassFilter(a.conf.Attenuation)

	a.server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL:        fmt.Sprintf("tcp://%s", a.conf.Host),
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, a)
	if err != nil {
		a.logger.Error().Err(err).Msg("Battery: failed to create modbus server")
		return err
	}

	if err = a.server.Start(); err != nil {
		a.logger.Error().Err(err).Msg("Battery: failed to start modbus server")
		return err
	}

	a.logger.Info().Msgf("Battery: Modbus server started on %s", a.conf.Host)
	return nil
}

func (a *Adapter) Cycle(simulatedTime *time.Time) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.simulatedTime != nil {
		deltaTime := simulatedTime.Sub(*a.simulatedTime)
		a.integraleEnergy_wh += a.integrate(a.p_w, deltaTime)
	}
	a.simulatedTime = simulatedTime
	a.setPoint_w = float64(math.Float32frombits(utils.Uint16ToUint32(a.setPoint1, a.setPoint2)))
	a.computeSetpoint()

	// Update power based on setpoint and filter

	if a.soc >= 100 && a.setPoint_w < 0 || a.soc <= 0 && a.setPoint_w > 0 {
		a.setPoint_w = 0
		a.logger.Info().Msg("Battery: Skipping setpoint due to SOC limits")
	} else {
		a.p_w = a.filter.Update(a.setPoint_w)
	}

	a.computeSoc()
	a.logger.Debug().Msgf("Battery: setpoint %f, p_w %f, soc %d, integral %f", a.setPoint_w, a.p_w, a.soc, a.integraleEnergy_wh)
}

func (a *Adapter) computeSoc() {
	a.soc = uint(((a.initialEnergy_wh - a.integraleEnergy_wh) / a.capacity_Wh) * 100)
	if a.soc > 100 {
		a.soc = 100
	}
}

func (a *Adapter) computeSetpoint() {
	if a.setPoint_w > 0 && a.setPoint_w > a.conf.PDischarge_w {
		a.setPoint_w = a.conf.PDischarge_w
	} else if a.setPoint_w < 0 && math.Abs(a.setPoint_w) > a.conf.PCharge_w {
		a.setPoint_w = -a.conf.PCharge_w
	}
}

func (a *Adapter) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {

	// loop through all register addresses from req.addr to req.addr + req.Quantity - 1
	for regAddr := req.Addr; regAddr < req.Addr+req.Quantity; regAddr++ {
		switch regAddr {
		case 0:
			res = append(res, uint16(a.soc))

		case 1:
			res = append(res, uint16(a.soh))

		case 2:

			val, _ := utils.Uint32ToUint16(math.Float32bits(float32(a.capacity_Wh)))
			res = append(res, val)
		case 3:
			_, val := utils.Uint32ToUint16(math.Float32bits(float32(a.capacity_Wh)))
			res = append(res, val)
		case 4:
			val, _ := utils.Uint32ToUint16(math.Float32bits(float32(a.p_w)))
			res = append(res, val)
		case 5:
			_, val := utils.Uint32ToUint16(math.Float32bits(float32(a.p_w)))
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
	var (
		regAddr uint16
	)
	a.lock.Lock()
	defer a.lock.Unlock()

	for i := 0; i < int(req.Quantity); i++ {
		// compute the target register address
		regAddr = req.Addr + uint16(i)

		switch regAddr {
		case 0:
			if req.IsWrite {
				a.setPoint1 = req.Args[i]
			}
			res = append(res, a.setPoint1)
		case 1:
			if req.IsWrite {
				a.setPoint2 = req.Args[i]
			}
			res = append(res, a.setPoint2)
		default:
			err = modbus.ErrIllegalDataAddress
			a.logger.Error().Err(err).Msgf("Battery: Holding register %d illegal address", regAddr)
		}

	}
	return res, nil
}

func (a *Adapter) integrate(p float64, delta time.Duration) float64 {
	return p * delta.Hours()
}

func (a *Adapter) Output() map[string]any {
	return map[string]any{
		"p_w": a.p_w,
	}
}

func (a *Adapter) Input(value any, key string) {
	a.logger.Warn().Msg("Battery: Adapter does not accept input")
}
