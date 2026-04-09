package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Code    string      `json:"code"`
	DB      DBConfig    `json:"db"` // 数据库配置子结构
	Redis   RedisConfig `json:"redis"`
}

type DBConfig struct {
	Type     string `json:"type"` // 数据库类型，如 "mysql", "postgres"
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"` // 数据库名
	// 可选：额外参数，如 charset, sslmode 等
	Params map[string]string `json:"params,omitempty"`
}

type RedisConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	DialTimeout  int    `json:"dial_timeout"`  // 秒
	ReadTimeout  int    `json:"read_timeout"`  // 秒
	WriteTimeout int    `json:"write_timeout"` // 秒
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
