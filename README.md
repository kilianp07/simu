
# Simu: Simulation Tool

Simu is a powerful simulation tool designed to model and simulate complex systems using a modular approach. It provides a set of packages for creating, configuring, and running simulations of various domains.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Running a Simulation](#running-a-simulation)
- [Contributing](#contributing)
## Overview

Simu is a versatile simulation tool that allows you to model and simulate dynamic systems with ease. It provides the following key components:

- **Core Package**: The `core` package serves as the central component of the simulation framework, managing configuration, adapters, links, and the simulation execution flow.

- **Adapters Packages**: Simu includes various adapter packages for different domains, such as `battery`, `poc`, `pv`, and `sum`, to model and simulate specific system behaviors.

- **Links Package**: The `links` package enables the creation of links between adapters, facilitating data exchange between simulation components.

## Installation

To get started with Simu, follow these installation steps:

1. Clone the Simu repository to your local machine:

   ```bash
   git clone https://github.com/kilianp07/simu.git
   ```

2. Navigate to the Simu project directory:

   ```bash
   cd simu
   ```

3. Install the necessary dependencies:

   ```bash
   make build
   ```

## Getting Started

To create and run simulations using Simu, you can follow these steps:

1. Define your simulation configuration by creating a JSON file, specifying simulation parameters and adapter configurations.

2. Launch the simulation using the `simu` command-line tool, providing the path to your configuration file:

   ```bash
   ./build/simu --conf your_simulation_config.json
   ```

3. Observe the simulation as it runs and generates results.

Thank you for the additional information. You can update the README to include details about the CSV format for defining links. Here's an updated section in the README to explain how to define links using the CSV file:

## Configuration

The configuration of your simulation is essential. You can customize various aspects of your simulation by editing the configuration file. The configuration includes settings for simulation time, adapter configurations, and links between adapters.

### Example Configuration:

```json
{
  "start": "2023-10-01T00:00:00Z",
  "end": "2023-10-01T23:59:59Z",
  "debug": false,
  "period_ms": 1000,
  "timestep_ms": 100,
  "links_path": "simulation_links.csv",  // Specify the path to your CSV file containing links
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
}
```

### Defining Links

To establish connections between adapters and define data exchange within your simulation, you can use a CSV file. The CSV file should have the following header format:

```
source, key, target, key
```

- `source`: The name of the source adapter from which data will be transferred.
- `key`: The key or identifier for the data to be transferred from the source adapter.
- `target`: The name of the target adapter to which data will be transferred.
- `key`: The key or identifier for the data to be transferred to the target adapter.

Here's an example CSV file for defining links:

```
PV_System_1, output_data, Battery_1, input_data
Battery_1, output_data, Sum_Operator_1, input_data
```

In the above example, links are defined between adapters `PV_System_1`, `Battery_1`, and `Sum_Operator_1`, specifying the source and target adapters along with the respective keys for data transfer.

## Running a Simulation

You can run your simulation by executing the `simu` command with your configuration file. The simulation will start and progress according to the defined parameters.

```bash
./build/simu --conf your_simulation_config.json
```



## Contributing

We welcome contributions to the Simu project. If you'd like to contribute, please follow these guidelines:

1. Fork the project.
2. Create a new branch for your feature or bug fix.
3. Make your changes and write tests if applicable.
4. Submit a pull request.
