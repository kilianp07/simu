package sum

import (
	"testing"
	"time"

	"github.com/kilianp07/simu/logger"
	"github.com/rs/zerolog"
)

const (
	confpath = "../../../test/adapters/operator/sum/sum.json"
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
		simulatedTime = time.Now()
		err           error
		values        = map[string]float64{
			"val1": 1.0,
			"val2": 2.0,
			"val3": 3.0,
			"val4": 4.0,
		}
	)

	if a, err = newAdapter(); err != nil {
		t.Fatal(err)
	}

	a.Cycle(&simulatedTime)

	for key, value := range values {
		a.Input(value, key)
	}

	a.Cycle(&simulatedTime)

	if a.Output()[keyResult] != 3.0 {
		t.Fatalf("Expected %f, got %f", 3.0, a.Output()[keyResult])
	}
}
