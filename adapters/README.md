# Adapters Package

The `adapters` package provides a collection of adapter implementations for various simulation components, including PV (photovoltaic) systems, batteries, Point of Connection (POC), sum operators, and dummy sender/receiver components. It allows you to create and manage different types of adapters within your simulation environment.

## Table of Contents
- [Interface](#interface)
- [Adapter Creation](#adapter-creation)

## Interface

### `Adapter`

The `Adapter` interface defines a set of common methods that all adapters within this package must implement. These methods include configuration, cycle, input, and output functions.

```go
type Adapter interface {
    Configure() error
    Cycle(simulatedTime *time.Time)
    Output() map[string]float64
    Input(value float64, key string)
}
```

This interface establishes a standard for interacting with simulation adapters, making it easier to manage and interact with various components.

## Adapter Creation

### `New`

```go
func New(name string, confpath string, simulatedTime *time.Time, logger *zerolog.Logger) Adapter
```

The `New` function is used to create new adapter instances based on the specified name. It initializes and returns an adapter that matches the given name.

#### Parameters:

- `name`: The name of the adapter to create.
- `confpath`: The path to the configuration file for the adapter.
- `simulatedTime`: A pointer to a `time.Time` representing simulated time.
- `logger`: A pointer to a `zerolog.Logger` for logging operations.

#### Returns:

An adapter instance that implements the `Adapter` interface, or `nil` if the specified adapter name is not recognized.

The following adapter types are available:

- `generic/pv`: Photovoltaic system adapter.
- `generic/battery`: Battery adapter.
- `generic/poc`: Point of Connection (POC) adapter.
- `operator/sum`: Sum operator adapter.
- `dummy/sender`: Dummy sender adapter.
- `dummy/receiver`: Dummy receiver adapter.

Please note that the `name` parameter should match one of the supported adapter types, as specified above. If an unsupported name is provided, a warning message will be logged, and `nil` will be returned.

This function provides a convenient way to create adapters based on their names and configuration paths.