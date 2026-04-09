package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Code    string `json:"code"`
}

var (
	defaultConfig *Config
	once          sync.Once
)

func Get() *Config {
	once.Do(func() {
		var err error
		defaultConfig, err = loadConfig()
		if err != nil {
			panic(fmt.Errorf("加载配置失败: %w", err))
		}
	})
	return defaultConfig
}

// loadConfig 内部加载逻辑，尝试多个路径
func loadConfig() (*Config, error) {
	c := &Config{}
	paths := []string{
		filepath.Join(executableDir(), "data", "config.json"),
		filepath.Join("data", "config.json"),
		filepath.Join("config.json"),
	}
	for _, p := range paths {
		if err := c.LoadFromFile(p); err == nil {
			return c, nil
		}
	}
	return nil, fmt.Errorf("未找到配置文件")
}

// LoadFromFile 从指定文件加载配置到当前结构体
func (c *Config) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}
