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

	// 可执行文件目录下
	exe, err := os.Executable()
	filename := filepath.Join(filepath.Dir(exe), "data", "config.json")
	err = c.unmarshal(filename)
	if err == nil {
		return c
	}

	// 当前目录下
	filename = filepath.Join("data", "config.json")
	err = c.unmarshal(filename)
	if err == nil {
		return c
	}

	panic(err)
}

/*
 * 读取配置文件转成配置类
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
