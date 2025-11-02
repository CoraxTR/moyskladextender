package main

import (
	"fmt"
	"mstorefgo/internal/config"
	"mstorefgo/internal/moyskladapi"
	"mstorefgo/internal/xlsxbuilder"
	"mstorefgo/usecases"
)

func main() {
	cfg := config.LoadConfig()

	msRateLimiter := moyskladapi.NewRatelimiter(cfg.RequestCap, cfg.TimeSpan)
	msProcessor := moyskladapi.NewMoySkladProcessor(msRateLimiter, &cfg.Moyskladapiconfig, &cfg.RefGoconfig)
	xlsxbuilder := xlsxbuilder.NewXlsxBuilder(*cfg)

	count, storage, err := usecases.PrepareUploadableOrders(msProcessor)
	if err != nil {
		panic(err)
	}
	usecases.BuildUploadXlsx(xlsxbuilder, *storage)
	usecases.ChangeStatusToShiped(msProcessor, *storage)
	usecases.ShipOrders(msProcessor, *storage)

	fmt.Println(count)
}
