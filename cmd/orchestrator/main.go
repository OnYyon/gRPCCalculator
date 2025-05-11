package main

import (
	"github.com/OnYyon/gRPCCalculator/internal/app"
	"github.com/OnYyon/gRPCCalculator/internal/config"
)

func main() {
	// TODO: сделать логирование

	// NOTE: поменял рабочий каталог
	// if err := os.Chdir("../../"); err != nil {
	// 	panic(err)
	// }
	cfg, err := config.Load("./internal/config/config.yaml")
	if err != nil {
		panic("don`t have config")
	}
	app := app.New(cfg)
	app.Run()
}
