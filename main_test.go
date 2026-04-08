package main_test

import (
	"github.com/adamchenEpm/ym3-go/internal/config"
	"testing"
)

func Test_NewConfig(t *testing.T) {

	cfg := config.NewConfig()
	//t.Assert(cfg != nil)

	t.Logf("Config.name: %v,  code :%v", cfg.Name, cfg.Code)
}
