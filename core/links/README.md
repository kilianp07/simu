# Links Package

The `links` package provides functionality for creating and managing links between different adapters. It allows you to establish connections between adapters based on a CSV file, enabling data transfer between them within a simulation environment.

## Table of Contents
- [Types](#types)
- [Creating Links](#creating-links)
- [Updating Links](#updating-links)

## Types

### `Links`

The `Links` type represents a collection of links between adapters. It contains an array of `Link` objects and maintains a reference to the available adapters.

```go
type Links struct {
    Links    []Link `json:"links"`
    adapters map[string]adapters.Adapter
}
```

### `Link`

The `Link` type defines a link between two adapters, specifying the source and target adapters and the associated keys.

```go
type Link struct {
    source data
    target data
}
```

### `data`

The `data` type represents information about an adapter, including its name and associated key.

```go
type data struct {
    adapterName string
    key         string
}
```

## Creating Links

### `New`

```go
func New(csvPath string, adapters map[string]adapters.Adapter) (*Links, error)
```

The `New` function creates links between adapters based on a CSV file. It initializes and returns a `Links` instance.

#### Parameters:

- `csvPath`: The path to the CSV file containing link information. The CSV file must have four columns: source_adapter, source_key, target_adapter, and target_key. These columns specify the source and target adapters and the corresponding keys used for data transfer.

- `adapters`: A map of adapter instances indexed by their names.

#### Returns:

A `Links` instance representing the created links, or an error if there is an issue with the CSV file or its structure.

This function reads the CSV file and creates links between adapters based on the information provided in the file.

## Updating Links

### `Update`

```go
func (l *Links) Update()
```

The `Update` method updates values from source adapters to target adapters based on the established links. It retrieves data from the source adapters and transfers it to the target adapters using the specified keys.

This method is used to synchronize data between connected adapters within the simulation environment.

