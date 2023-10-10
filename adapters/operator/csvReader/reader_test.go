package csvreader

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestCycle(t *testing.T) {
	// Create a logger for testing (you can customize this if needed)
	logger := zerolog.New(nil)

	// Example test data
	testData := [][]string{
		{"date", "value"},
		{"2023-10-10 12:00:00", "1.0"},
		{"2023-10-10 12:15:00", "2.0"},
		{"2023-10-10 12:30:00", "3.0"},
	}

	conf := &conf{
		StartDate:  "2023-10-10 12:00:00",
		TimeFormat: "2006-01-02 15:04:05",
		CsvPath:    "../../../test/adapters/operator/reader/reader.csv",
		DataCol:    1,
		DateCol:    0,
	}

	adapter := &Adapter{
		conf:   conf,
		logger: &logger,
		data:   testData,
	}

	// Set the initial simulated time
	simulatedTime := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	adapter.simulatedTime = &simulatedTime

	// First Cycle call
	adapter.Cycle(&simulatedTime)

	if adapter.value != 1.0 {
		t.Fatalf("Expected value 1.0, got %f", adapter.value)
	}

	// Move simulated time forward by 15 minutes
	simulatedTime = simulatedTime.Add(15 * time.Minute)
	adapter.Cycle(&simulatedTime)

	if adapter.value != 2.0 {
		t.Fatalf("Expected value 2.0, got %f", adapter.value)
	}
}

func TestOutput(t *testing.T) {
	// Create a logger for testing
	logger := zerolog.New(nil)

	// Example test data
	testData := [][]string{
		{"date", "value"},
		{"2023-10-10 12:00:00", "1.0"},
		{"2023-10-10 12:15:00", "2.0"},
	}

	conf := &conf{
		StartDate:  "2023-10-10 12:00:00",
		TimeFormat: "2006-01-02 15:04:05",
		CsvPath:    "test.csv",
		DataCol:    1,
		DateCol:    0,
	}

	adapter := &Adapter{
		conf:   conf,
		logger: &logger,
		data:   testData,
	}

	// Set the initial simulated time
	simulatedTime := time.Date(2023, 10, 10, 12, 0, 0, 0, time.UTC)
	adapter.simulatedTime = &simulatedTime

	// Ensure the output is correct after the first Cycle call
	adapter.Cycle(&simulatedTime)
	expectedOutput := map[string]any{"value": 1.0}
	output := adapter.Output()

	if output["value"] != expectedOutput["value"] {
		t.Fatalf("Expected output %v, got %v", expectedOutput, output)
	}
}
