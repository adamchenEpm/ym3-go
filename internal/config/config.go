package config

import (
	"embed"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewConfig() *Config {
	c := &Config{}

	c, err := c.getConfigOuter()
	if err != nil {
		c, err = c.getConfigInner()
		if err != nil {
			panic(err)
		}
	}

	return c
}

/*
 * 从外部获取配置文件
 */
func (c *Config) getConfigOuter() (*Config, error) {

	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exePath := filepath.Dir(exe)
	exePath = "E:/app1/ym3-go"
	configPath := filepath.Join(exePath, "data", "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
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

	return c, nil
}

/*
 * 从内部获取配置文件
 */
func (c *Config) getConfigInner() (*Config, error) {

	var configFile embed.FS
	data, err := configFile.ReadFile("data/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
