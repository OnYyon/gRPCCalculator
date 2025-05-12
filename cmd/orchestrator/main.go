package main

import (
	"github.com/OnYyon/gRPCCalculator/internal/app"
	"github.com/OnYyon/gRPCCalculator/internal/config"
)

func main() {
	cfg, err := config.Load("./internal/config/config.yaml")
	if err != nil {
		panic("don`t have config")
	}
	app := app.New(cfg)
	app.Run()
}
