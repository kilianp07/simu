package core

import (
	"time"

	"github.com/kilianp07/simu/adapters"
	"github.com/kilianp07/simu/core/links"
	"github.com/kilianp07/simu/utils"
	"github.com/rs/zerolog"
)

type Runner struct {
	conf          *Conf
	adapters      map[string]adapters.Adapter
	links         *links.Links
	simulatedTime *time.Time
	logger        *zerolog.Logger
	lastCycle     time.Time
}

func (r *Runner) readConfig(confPath string) (*Conf, error) {

	conf := &Conf{}

	if err := utils.ReadJsonFile(confPath, conf); err != nil {
		return nil, err
	}

	if conf.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return conf, nil
}

func (r *Runner) instanciate() {
	for _, adapter := range r.conf.Adapters {
		if a := adapters.New(adapter.Adapter, adapter.ConfPath, r.simulatedTime, r.logger); a != nil {
			r.adapters[adapter.Name] = *a
		}
	}
}

func (r *Runner) configureAdapters() error {
	for _, adapter := range r.adapters {
		if err := adapter.Configure(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) configureLinks() error {
	var (
		l   *links.Links
		err error
	)
	if l, err = links.New(r.conf.LinksPath, r.adapters); err != nil {
		return err
	}

	r.links = l

	return nil
}

func (r *Runner) run() {
	for {
		// Wait for the next cycle
		if time.Since(r.lastCycle) < time.Duration(r.conf.Timestep)*time.Millisecond {
			continue
		}

		// At each cycle, we increment the simulated time
		t := r.simulatedTime.Add(time.Duration(r.conf.Period) * time.Millisecond)
		r.simulatedTime = &t

		// Update I/O
		r.links.Update()

		for _, adapter := range r.adapters {
			adapter.Cycle(r.simulatedTime)
		}

		r.lastCycle = time.Now()
	}
}
