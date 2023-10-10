# CSV Reader Adapter

The `csvreader` package provides functionality for reading CSV data and using it as a data source for simulation. It is useful when you want to simulate a system based on data stored in CSV files, such as time-series data.

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)

## Overview

The `csvreader` adapter allows you to read CSV files containing time-series data and use that data as input for your simulation. This can be particularly useful when you need to simulate systems based on historical or recorded data.

## Installation

To use the `csvreader` adapter in your Simu project, follow these steps:

1. Import the `csvreader` package into your project:

   ```go
   import "github.com/kilianp07/simu/adapters/csvreader"
   ```

2. Instantiate the `csvreader` adapter with your configuration, simulated time, and logger.

   ```go
   adapter := csvreader.New(confpath, simulatedTime, logger)
   ```

3. Configure the adapter to read CSV data and set up the simulation parameters.

   ```go
   err := adapter.Configure()
   if err != nil {
       // Handle configuration error
   }
   ```

4. Use the adapter's data as input for your simulation.

## Usage

The `csvreader` adapter reads CSV data and provides it as input for your simulation. It uses time-based data columns to update values at specified time intervals.

```go
// Example usage:
for {
    // Simulated time is updated here
    adapter.Cycle(simulatedTime)

    // Get the current value from the adapter's output
    value := adapter.Output()["value"]
    
    // Use 'value' in your simulation
}
```

## Configuration

You can configure the `csvreader` adapter by providing a configuration file. The configuration file specifies details such as the CSV file path, date column, data column, and time format.

### Example Configuration:

```json
{
  "start_date": "2023-10-01T00:00:00Z",
  "time_format": "2006-01-02T15:04:05Z",
  "csv_path": "your_data.csv",
  "data_column": 1,
  "date_column": 0
}
```

- `start_date`: The start date for the simulation.
- `time_format`: The format of date and time in the CSV file.
- `csv_path`: The path to the CSV data file.
- `data_column`: The column index in the CSV file containing the data values.
- `date_column`: The column index in the CSV file containing date and time information.