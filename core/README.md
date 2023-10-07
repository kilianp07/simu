# Core Package

The `core` package serves as the central component of the simulation framework. It manages the configuration, adapters, links, and the simulation execution flow.

## Table of Contents
- [Types](#types)
- [Launching the Simulation](#launching-the-simulation)
- [Simulation Execution](#simulation-execution)
- [Creating a Configuration](#creating-a-configuration)

## Types

### `Conf`

The `Conf` type represents the configuration for the simulation. It includes parameters such as start and end times, debug mode, time intervals, links configuration, and a list of simulation blocks (adapters).

```go
type Conf struct {
    Start     time.Time  `json:"start"`
    End       time.Time  `json:"end"`
    Debug     bool       `json:"debug"`
    Period    uint       `json:"period_ms"`
    Timestep  uint       `json:"timestep_ms"`
    LinksPath string     `json:"links_path"`
    Adapters  []SimBlock `json:"adapters"`
}
```

### `SimBlock`

The `SimBlock` type represents a simulation block, including the adapter name, block name, and the configuration file path.

```go
type SimBlock struct {
    Adapter  string `json:"adapter"`
    Name     string `json:"name"`
    ConfPath string `json:"confPath"`
}
```

## Launching the Simulation

### `Launch`

```go
func Launch(confpath string)
```

The `Launch` function initiates the simulation by reading the configuration, configuring adapters and links, and running the simulation loop.

#### Parameters:

- `confpath`: The path to the configuration file.

This function sets up and starts the simulation based on the provided configuration file. It manages the simulation lifecycle, including configuration, initialization, and continuous execution.

## Simulation Execution

### `Runner`

The `Runner` type is an internal structure used for managing the simulation execution flow. It includes methods for reading the configuration, instantiating adapters, configuring adapters and links, and running the simulation.

```go
type Runner struct {
    conf          *Conf
    adapters      map[string]adapters.Adapter
    links         *links.Links
    simulatedTime *time.Time
    logger        *zerolog.Logger
    lastCycle     time.Time
}
```

### `readConfig`

```go
func (r *Runner) readConfig(confPath string) (*Conf, error)
```

The `readConfig` method reads the configuration file specified by `confPath` and returns a `Conf` instance representing the simulation configuration.

### `instanciate`

```go
func (r *Runner) instanciate()
```

The `instanciate` method initializes and instantiates adapter instances based on the configuration.

### `configureAdapters`

```go
func (r *Runner) configureAdapters() error
```

The `configureAdapters` method configures the adapter instances, ensuring that they are ready for simulation.

### `configureLinks`

```go
func (r *Runner) configureLinks() error
```

The `configureLinks` method configures the links between adapters, establishing connections for data transfer.

### `run`

```go
func (r *Runner) run()
```

The `run` method is responsible for the core execution of the simulation. It continuously updates the simulation state, handles I/O operations between adapters, and advances the simulated time.

This method is the heart of the simulation engine, ensuring that data flows between adapters and that the simulation progresses according to the configured time intervals.



## Creating a Configuration

To create a configuration for your simulation, you need to define the parameters for the simulation blocks (adapters) and the overall simulation settings. Here are the steps to create a configuration:

1. **Define Simulation Blocks (Adapters):**

   Start by listing the simulation blocks (adapters) you want to include in your simulation. Each block should be specified with its name, the corresponding adapter type, and the path to its configuration file (if applicable). For example:
```json
   "adapters": [
       {
           "adapter": "generic/pv",
           "name": "PV_System_1",
           "confPath": "pv_system_1_config.json"
       },
       {
           "adapter": "generic/battery",
           "name": "Battery_1",
           "confPath": "battery_config.json"
       },
       {
           "adapter": "operator/sum",
           "name": "Sum_Operator_1",
           "confPath": "sum_operator_config.json"
       }
   ]
   ```

   In the above example, three simulation blocks are defined, including a photovoltaic system, a battery, and a sum operator.

2. **Set Simulation Parameters:**

   Configure the overall simulation parameters, such as the start and end times, debug mode, time intervals, links configuration, and any other settings specific to your simulation. Here's an example configuration:

   ```json
   {
       "start": "2023-10-01T00:00:00Z",
       "end": "2023-10-01T23:59:59Z",
       "debug": false,
       "period_ms": 1000,
       "timestep_ms": 100,
       "links_path": "simulation_links.csv",
       "adapters": [
           ...
       ]
   }
   ```

   - `start` and `end`: Define the start and end times for the simulation.
   - `debug`: Set to `true` to enable debug mode for detailed logging.
   - `period_ms` and `timestep_ms`: Specify the period and timestep in milliseconds for the simulation.
   - `links_path`: Provide the path to the CSV file containing link information.

3. **Save Configuration:**

   Save the configuration as a JSON file. For example, you can save it as `simulation_config.json`.

4. **Launch the Simulation:**

   To launch the simulation with your custom configuration, call the bin, passing the path to the configuration file as an argument:

   ```bash
    simu --conf simulation_config.json
   ```

By following these steps, you can create a custom configuration for your simulation, specifying the adapters to use and their settings. This allows you to tailor the simulation to your specific needs and experiment with different configurations.
