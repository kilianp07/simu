## Package `pv`

The `pv` package provides functionality for simulating a photovoltaic (PV) system. It allows interaction with a Modbus server.

### Table of Contents
- [Types](#Types)
- [Configuration](#Configuration)
- [Modbus Interface](#Modbus-Registers)
- [Example](#Example)
### Types

#### `Adapter`

The `Adapter` type is the core component of the `pv` package. It manages the communication with a Modbus server and simulates a PV system's power output based on CSV data.

```go
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
```

##### Methods

- `New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter`: Creates a new `Adapter` instance and initializes it with the provided configuration path, simulated time, and logger.
- `Configure() error`: Reads the configuration from the specified path, initializes the Modbus server, and starts it for communication.
- `Cycle(simulatedTime *time.Time)`: Simulates a cycle of the PV system, updating power output based on CSV data.
- `HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error)`: Handles Modbus input registers requests for reading power output.
- `HandleCoils(req *modbus.CoilsRequest) (res []bool, err error)`: Handles Modbus coils requests (not implemented).
- `HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error)`: Handles Modbus discrete inputs requests (not implemented).
- `HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error)`: Handles Modbus holding registers requests (not implemented).
- `Input(value any, key string)`: Handles input values for the PV system.
- `Output() map[string]any`: Provides a map containing the current power output of the PV system.

### Configuration

The `pv` package relies on a JSON configuration file to set its parameters. Here's an example configuration structure:

```json
{
  "host": "localhost:502"
}
```
- `host`: The Modbus server host address.

### Modbus Registers

Here is a table of Modbus registers used in the `pv` package, along with their data types and units:

| Register Address | Description              | Data Type | Unit       |
|------------------|--------------------------|-----------|------------|
| 0                | Power Output (p_w)       | int32   | Watts (W)  |


This table lists the Modbus register addresses, provides a brief description of each register's purpose, specifies the data type of the data stored in the register, and indicates the unit of measurement for the data.

### Example

Here's an example of how to use the `pv` package:

```go
import (
	"time"
	"github.com/rs/zerolog"
	"github.com/kilianp07/simu/pv"
)

func main() {
	// Create a logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: nil}).With().Logger()

	// Create a simulated time
	simulatedTime := time.Now()

	// Create a new PV Adapter
	adapter := pv.New("config.json", &simulatedTime, &logger)

	// Configure the PV Adapter
	err := adapter.Configure()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to configure PV system")
		return
	}

	// Perform a cycle of PV simulation
	adapter.Cycle(&simulatedTime)

	// Get the current power output
	output := adapter.Output()
	logger.Info().Msgf("Current power output: %f Watts", output["p_w"])
}
```
