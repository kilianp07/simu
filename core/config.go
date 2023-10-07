package core

import "time"

type Conf struct {
	Start     time.Time  `json:"start"`
	End       time.Time  `json:"end"`
	Debug     bool       `json:"debug"`
	Period    uint       `json:"period_ms"`
	Timestep  uint       `json:"timestep_ms"`
	LinksPath string     `json:"links_path"`
	Adapters  []SimBlock `json:"adapters"`
}

type SimBlock struct {
	Adapter  string `json:"adapter"`
	Name     string `json:"name"`
	ConfPath string `json:"confPath"`
}
