package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Code    string `json:"code"`
}

func NewConfig() *Config {
	c := &Config{}

	// exe
	exe, err := os.Executable()
	filename := filepath.Join(filepath.Dir(exe), "data", "config.json")
	err = c.unmarshal(filename)
	if err == nil {
		return c
	}

	//lcoal
	filename = filepath.Join("data", "config.json")
	err = c.unmarshal(filename)
	if err == nil {
		return c
	}

	panic(err)
}

/*
 * 取配置文件
 */
func (c *Config) unmarshal(filename string) error {

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	return nil
}
