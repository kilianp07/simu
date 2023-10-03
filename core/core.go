package core

import (
	"log"
	"os"

	"github.com/kilianp07/simu/logger"
)

func Launch(confpath string) {
	var (
		conf *Conf
		err  error
		r    = &Runner{
			logger: logger.Get(),
		}
	)

	if conf, err = r.readConfig(confpath); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	r.conf = conf
	r.simulatedTime = &conf.Start

	r.instanciate()

	if err = r.configureAdapters(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	r.run()
}
