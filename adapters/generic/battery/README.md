# Battery Adapter

The `battery` package provides functionality for simulating battery behavior, including state of charge (SOC), state of health (SOH), and power management. It is intended to be used in simulation environments where you need to model the behavior of a battery.

## Table of Contents
- [Types](#Types)
- [Methods](#Methods)
- [Configuration](#Configuration)
- [Modbus Interface](./Interface.md)

## Types
### `Adapter`

The `Adapter` type is the core component of the `battery` package. It simulates the behavior of a battery and allows you to interact with it.

```go
type Adapter struct {
	server        *modbus.ModbusServer
	conf          *conf
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
	initialeEnergy_wh  float64
}
```
### `Methods`
- New
```go
New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter
```
Creates a new Adapter instance and initializes it with the provided configuration path, simulated time, and logger.
- Configure
```go 
Configure() error
````
Reads the configuration from the specified path, initializes the battery parameters, and starts the Modbus server for communication.

- Cycle:
```go
Cycle(simulatedTime *time.Time)
````
Performs a cycle of battery simulation, updating SOC, power output, and other parameters.
- Input Registers: 
```go
HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error)
````
Handles Modbus input registers requests for reading SOC, SOH, and capacity values.
- Coils Registers: 
```go
HandleCoils(req *modbus.CoilsRequest) (res []bool, err error)
````
Handles Modbus coils requests (not implemented).

- Discrete Inputs : 
```go
HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error)
```
Handles Modbus discrete inputs requests (not implemented).

- Holding Registers:
```go
HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error)
```
Handles Modbus holding registers requests for setting battery setpoints.

- Ouput:
```go
Output() map[string]float64
```
Returns a map containing the current power output of the battery.
- Input:
```go
Input(value float64, key string)
```
Handles input values for the battery (not implemented).

## Configuration
The battery package relies on a JSON configuration file to set its parameters. Here's an example configuration structure:

```json

{
  "soc": 50,
  "soh": 90,
  "capacity_wh": 1000.0,
  "p_charge_w": 500.0,
  "p_discharge_w": 500.0,
  "host": "localhost:502",
  "attenuation": 0.1
}
```
- `soc`: The initial state of charge (SOC) of the battery (0-100).
- `soh`: The initial state of health (SOH) of the battery (0-100).
- `capacity_wh`: The total capacity of the battery in watt-hours (Wh).
- `p_charge_w`: The maximum power that can be charged into the battery in watts (W).
- `p_discharge_w`: The maximum power that can be discharged from the battery in watts (W).
- `host`: The Modbus server host address.
- `attenuation`: The attenuation factor for the low-pass filter used in power output filtering.