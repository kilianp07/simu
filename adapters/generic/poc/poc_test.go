package poc

import (
	"fmt"
	"testing"
	"time"

	"github.com/kilianp07/simu/logger"
	"github.com/rs/zerolog"
	"github.com/simonvetter/modbus"
)

const (
	confpath = "../../../test/adapters/generic/poc/poc.json"

	// Readings
	p_wReg = 0
)

func newAdapter() (*Adapter, error) {
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

func TestCycle(t *testing.T) {
	var (
		a             *Adapter
		client        *modbus.ModbusClient
		simulatedTime = time.Now()
		value         = 100.0
		result        uint32
		err           error
	)

	if a, err = newAdapter(); err != nil {
		t.Fatal(err)
	}

	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("tcp://%s", a.conf.Host),
		Timeout: 1 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Open()
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}
	defer client.Close()

	// Set Input
	a.Input(value, keyPw)
	a.Cycle(&simulatedTime)

	// Read register
	if result, err = client.ReadUint32(p_wReg, modbus.INPUT_REGISTER); err != nil {
		t.Fatal(err)
	}

	if result != uint32(value) {
		t.Fatalf("Expected %f, got %d", value, result)
	}

}
