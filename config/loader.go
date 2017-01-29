package config

import (
	"encoding/json"
	"io/ioutil"
)

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = json.Unmarshal(data, cfg)
	return cfg, err
}
