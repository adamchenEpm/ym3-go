package config_test

import (
	"github.com/adamchenEpm/go-ym3/internal/config"
	"testing"
)

func Test_NewConfig(t *testing.T) {

	cfg := config.NewConfig()
	t.Assert(cfg != nil)

	t.Logf("Config.name: %v", cfg.Name)
}
