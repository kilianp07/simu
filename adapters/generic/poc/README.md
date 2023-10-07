## Package `poc`

The `poc` package provides functionality for a generic point of control (POC) system. It allows interaction with a Modbus server and handles power control for a specific key named `p_w`.

## Table of Contents
- [Types](#Types)
- [Methods](#Methods)
- [Configuration](#Configuration)
- [Modbus Interface](#modbus-registers)

### Types

#### `Adapter`

The `Adapter` type is the core component of the `poc` package. It manages the communication with a Modbus server and handles power control through the `p_w` key.

```go
type Adapter struct {
	server   *modbus.ModbusServer
	conf     *conf
	logger   *zerolog.Logger
	confPath string

	p_w int32
}
```

##### Methods

- `New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter`: Creates a new `Adapter` instance and initializes it with the provided configuration path, simulated time, and logger.
- `Configure() error`: Reads the configuration from the specified path, initializes the Modbus server, and starts it for communication.
- `Cycle(simulatedTime *time.Time)`: Placeholder method for a cycle of POC simulation.
- `HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error)`: Handles Modbus input registers requests for reading the `p_w` power value.
- `HandleCoils(req *modbus.CoilsRequest) (res []bool, err error)`: Handles Modbus coils requests (not implemented).
- `HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error)`: Handles Modbus discrete inputs requests (not implemented).
- `HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error)`: Handles Modbus holding registers requests (not implemented).
- `Output() map[string]float64`: Provides a map containing the current power output of the POC (not implemented).
- `Input(value float64, key string)`: Handles input values for the POC, specifically for the `p_w` key.

### Configuration

The `poc` package uses a JSON configuration file to set its parameters. Here's an example configuration structure:

```json
{
  "host": "localhost:502"
}
```

- `host`: The Modbus server host address.

### Example

Here's an example of how to use the `poc` package:

```go
import (
	"time"
	"github.com/rs/zerolog"
	"github.com/kilianp07/simu/poc"
)

func main() {
	// Create a logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: nil}).With().Logger()

	// Create a simulated time
	simulatedTime := time.Now()

	// Create a new POC Adapter
	adapter := poc.New("config.json", &simulatedTime, &logger)

	// Configure the POC Adapter
	err := adapter.Configure()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to configure POC")
		return
	}

	// Perform a cycle of POC simulation (if applicable)
	// adapter.Cycle(&simulatedTime)

	// Handle input for the `p_w` key
	adapter.Input(100.0, "p_w")
}
```

## Modbus Registers


| Register Address | Description              | Data Type | Unit       |
|------------------|--------------------------|-----------|------------|
| 0                | Power Output (p_w)       | int32     | Watts (W)  |


