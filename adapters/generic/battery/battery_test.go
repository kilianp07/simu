package battery

import (
	"fmt"
	"testing"
	"time"

	"github.com/kilianp07/simu/logger"
	"github.com/rs/zerolog"
	"github.com/simonvetter/modbus"
)

const (
	confpath = "../../../test/battery/battery.json"

	// Readings
	socReg         = 0
	sohReg         = 1
	capacityReg    = 2
	activePowerReg = 4

	// Writings
	setpointReg = 0
)

func new() (*Adapter, error) {
	var (
		simulatedTime *time.Time
		logger        *zerolog.Logger = logger.Get()
	)

	now := time.Now()
	simulatedTime = &now

	a := New(confpath, simulatedTime, logger)

	if err := a.Configure(); err != nil {
		return nil, err
	}

	return a, nil
}

func TestCyle(t *testing.T) {
	var (
		client        *modbus.ModbusClient
		a             *Adapter
		err           error
		simulatedTime time.Time

		soc         uint16
		soh         uint16
		capacity    uint32
		setpoint    int32
		activePower uint32
	)

	if a, err = new(); err != nil {
		t.Fatalf("failed to create adapter: %v", err)
	}

	simulatedTime = time.Now()

	client, err = modbus.NewClient(
		&modbus.ClientConfiguration{
			URL:     fmt.Sprintf("tcp://%s", a.conf.Host),
			Timeout: 1 * time.Second,
		},
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	a.Cycle(&simulatedTime)

	err = client.Open()
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}
	defer client.Close()

	/*
		Readings
	*/
	// Read SOC
	soc, err = client.ReadRegister(socReg, modbus.INPUT_REGISTER)
	if err != nil {
		t.Fatalf("failed to read register: %v", err)
	}
	if soc != uint16(a.conf.Soc) {
		t.Fatalf("expected %v, got %v", a.conf.Soc, soc)
	}

	// Read SOH
	soh, err = client.ReadRegister(sohReg, modbus.INPUT_REGISTER)
	if err != nil {
		t.Fatalf("failed to read register: %v", err)
	}
	if soh != uint16(a.conf.Soh) {
		t.Fatalf("expected %v, got %v", a.conf.Soh, soh)
	}

	// Read Capacity
	capacity, err = client.ReadUint32(capacityReg, modbus.INPUT_REGISTER)
	if err != nil {
		t.Fatalf("failed to read register: %v", err)
	}
	if float64(int32(capacity)) != a.conf.Capacity_Wh {
		t.Fatalf("expected %v, got %v", a.conf.Capacity_Wh, capacity)
	}

	/*
		Writings
	*/
	// Write Positive Setpoint
	setpoint = 10
	err = client.WriteUint32(setpointReg, uint32(setpoint))
	if err != nil {
		t.Fatalf("failed to write register: %v", err)
	}

	simulatedTime = simulatedTime.Add(1 * time.Second)
	a.Cycle(&simulatedTime)

	if int32(a.setPoint_w) != setpoint {
		t.Fatalf("expected %v, got %v", setpoint, a.setPoint_w)
	}

	// Write Negative Setpoint
	setpoint = -10
	err = client.WriteUint32(setpointReg, uint32(setpoint))
	if err != nil {
		t.Fatalf("failed to write register: %v", err)
	}

	simulatedTime = simulatedTime.Add(1 * time.Second)
	a.Cycle(&simulatedTime)

	if int32(a.setPoint_w) != setpoint {
		t.Fatalf("expected %v, got %v", setpoint, a.setPoint_w)
	}

	simulatedTime = simulatedTime.Add(1 * time.Second)
	a.Cycle(&simulatedTime)

	// Read Active Power
	activePower, err = client.ReadUint32(activePowerReg, modbus.INPUT_REGISTER)
	if err != nil {
		t.Fatalf("failed to read register: %v", err)
	}

	if float64(int32(activePower)) == 0 {
		t.Fatalf("expected not null active power, got %v", activePower)
	}

}
