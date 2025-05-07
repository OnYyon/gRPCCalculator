package main

import (
	app "github.com/OnYyon/gRPCCalculator/internal/app/rest"
)

func main() {
	// TODO: загрузка конфигуряция из .env
	// TODO: сделать логирование

	app.StartOrchestrator()
}
