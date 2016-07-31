package config

import (
	"encoding/json"
	"io"
)

type Config struct {
	Address     string `json:"address"`     // Interface and port to listen requests.
	Target      string `json:"target"`      // Address of target server.
	Pattern     string `json:"pattern"`     // Regexp pattern.
	Replacement string `json:"replacement"` // Replacement string.
}

func ReadConfig(r io.Reader) (Config, error) {
	conf := Config{}

	err := json.NewDecoder(r).Decode(&conf)

	return conf, err
}
