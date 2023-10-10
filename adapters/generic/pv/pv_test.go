package pv

import (
	"testing"
	"time"

	"github.com/kilianp07/simu/logger"
	"github.com/kilianp07/simu/utils"
	"github.com/simonvetter/modbus"
)

const (
	confpath = "../../../test/adapters/generic/pv/config.json"
)

func newAdapter() (*Adapter, error) {
	var (
		simulatedTime *time.Time
	)

	now := time.Now()
	simulatedTime = &now

	a := New(confpath, simulatedTime, logger.Get())

	if err := a.Configure(); err != nil {
		return nil, err
	}

	return a, nil
}

func TestHandleInputRegisters(t *testing.T) {
	// Example test data
	p_w := 0.75 // Test input value

	// Create a new Adapter
	a, err := newAdapter()
	if err != nil {
		t.Fatal(err)
	}

	// Input the test value into the Adapter
	a.Input(p_w, "p_w")

	// Create a mock modbus InputRegistersRequest
	req := &modbus.InputRegistersRequest{
		Addr:     0, // Mock address
		Quantity: 2, // Mock quantity
	}

	// Call the HandleInputRegisters method with the mock request
	res, err := a.HandleInputRegisters(req)

	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected 2 values in response, got %d", len(res))
	}

	// Check if the response values match the expected values
	expectedValue := uint32(p_w * 100)

	if utils.Uint16ToUint32(res[0], res[1]) != expectedValue {
		t.Fatalf("Expected value %d, got %d", expectedValue, utils.Uint16ToUint32(res[0], res[1]))
	}

	_ = a.server.Stop()
}

func TestInput(t *testing.T) {
	// Example test data
	p_w := 0.75 // Test input value

	// Create a new Adapter
	a, err := newAdapter()
	if err != nil {
		t.Fatal(err)
	}

	// Input the test value into the Adapter
	a.Input(p_w, "p_w")

	// Check if the internal state of the Adapter was updated correctly
	if a.p_w != p_w {
		t.Fatalf("Expected p_w value %f, got %f", p_w, a.p_w)
	}
}
