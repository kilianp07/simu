package links

import (
	"errors"

	"github.com/kilianp07/simu/adapters"
	"github.com/kilianp07/simu/logger"
	"github.com/kilianp07/simu/utils"
)

type Links struct {
	Links    []Link `json:"links"`
	adapters map[string]adapters.Adapter
}

type Link struct {
	source data
	target data
}

type data struct {
	adapterName string
	key         string
}

// New creates links between adapters based on a csv file
func New(csvPath string, adapters map[string]adapters.Adapter) (*Links, error) {
	var (
		err     error
		csvdata [][]string
		links   *Links
	)

	csvdata, err = utils.ReadCsvFile(csvPath)
	if err != nil {
		return nil, err
	}

	if len(csvdata[0]) != 4 {
		return nil, errors.New("csv file must have 4 columns: source_adapter, source_key, target_adapter, target_key")
	}

	links = &Links{
		adapters: adapters,
	}

	for i, row := range csvdata {

		// Skip headers
		if i == 0 {
			continue
		}
		links.Links = append(links.Links, Link{
			source: data{
				adapterName: row[0],
				key:         row[1],
			},
			target: data{
				adapterName: row[2],
				key:         row[3],
			},
		})

	}

	return links, nil

}

// Update sets values from source adapters to target adapters
func (l *Links) Update() {
	var (
		source adapters.Adapter
		target adapters.Adapter
		ok     bool

		data any
	)

	for _, link := range l.Links {
		source, ok = l.adapters[link.source.adapterName]
		if !ok {
			logger.Get().Error().Msgf("Output adapter %s not found", link.source.adapterName)
			continue
		}
		target, ok = l.adapters[link.target.adapterName]
		if !ok {
			logger.Get().Error().Msgf("Input adapter %s not found", link.target.adapterName)
			continue
		}

		data = source.Output()[link.source.key]
		target.Input(data, link.target.key)
	}
}
