package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Enabled bool   `json:"enabled"`
}

func NewConfig() *Config {
	c := &Config{}
	c.decode()
	return c
}

func (c *Config) decode() *Config {

	file, err := os.Open("data/config.json")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		panic(err)
	}

	return c
}
