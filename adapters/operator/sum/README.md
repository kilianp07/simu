# Sum Adapter

The `sum` package provides functionality for performing weighted sum calculations based on a configuration. This package is designed for scenarios where you need to calculate a weighted sum of values.

## Table of Contents
- [Types](#types)
- [Methods](#methods)
- [Configuration](#configuration)

## Types

### `Adapter`

The `Adapter` type is the central component of the `sum` package. It facilitates weighted sum calculations and offers methods for configuration and result retrieval.

```go
type Adapter struct {
    conf     *conf           // Configuration for the sum operation.
    logger   *zerolog.Logger // Logger for logging operations.
    confPath string          // Path to the configuration file.
    result   float64         // Result of the sum operation.
}
```

## Methods

### `New`

```go
func New(confpath string, simulatedTime *time.Time, logger *zerolog.Logger) *Adapter
```

Creates a new `Adapter` instance and initializes it with the provided configuration path, simulated time, and logger.

### `Configure`

```go
func (a *Adapter) Configure() error
```

Reads the configuration from the specified path, initializes the `Adapter` with the provided values, and prepares it for calculations.

### `Cycle`

```go
func (a *Adapter) Cycle(simulatedTime *time.Time)
```

Performs a cycle of the sum operation, calculating the weighted sum based on the current configuration.

### `Output`

```go
func (a *Adapter) Output() map[string]float64
```

Returns the result of the sum operation as a map containing the calculated sum.

### `Input`

```go
func (a *Adapter) Input(value float64, key string)
```

Handles input values for the sum operation. This method allows you to update the values of specific members in the configuration.

## Configuration

The `sum` package relies on a JSON configuration file to specify its parameters. Below is an example configuration structure:

```json
{
  "members": [
    {
      "value": 10.5,
      "key": "member1",
      "weight": 0.2
    },
    {
      "value": 8.0,
      "key": "member2",
      "weight": 0.5
    }
  ]
}
```

- `members`: An array of member objects, each representing a value and weight.
  - `value`: The initial value of the member.
  - `key`: A unique identifier for the member.
  - `weight`: The weight assigned to the member in the weighted sum calculation.

This JSON configuration allows you to define the members and their initial values and weights for the sum operation.

