package main

import (
	"github.com/OnYyon/gRPCCalculator/internal/app"
	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/OnYyon/gRPCCalculator/internal/storage/sqlite"
)

func main() {
	// TODO: сделать логирование
	cfg, err := config.Load("./internal/config/config.yaml")
	if err != nil {
		panic("don`t have config")
	}
	sqlite.MustRunNewStorage(cfg)
	app.StartOrchestrator(cfg)
}
