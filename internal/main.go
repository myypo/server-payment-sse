package main

import (
	"payment-sse/internal/config"
	"payment-sse/internal/controller"
	"payment-sse/internal/controller/service"
	dom "payment-sse/internal/domain"
	"payment-sse/internal/protocol/http"
	"payment-sse/internal/repo/postgres/order"
	postgres "payment-sse/internal/repo/postgres/provider"
	pgTx "payment-sse/internal/repo/postgres/transaction"

	"go.uber.org/zap"
)

func main() {
	conf := must(config.NewConfig())

	log := (func() *zap.Logger {
		switch conf.Dom.LogLevel {
		case dom.Debug:
			return must(zap.NewDevelopment())
		case dom.Info:
			return must(zap.NewProduction())
		}
		return must(zap.NewProduction())
	})()

	pg := (func() *postgres.PostgresProvider {
		p := must(postgres.NewPostgresProvider(conf.PG))
		if err := p.MigrateUp(); err != nil {
			panic(err)
		}
		return p
	})()

	pgOrd := pgOrd.NewOrderPostgres(pg.Conn())
	pgTx := pgTx.NewPGTransaction(pg.Conn())

	streamEve := service.NewEventStreamer(pgOrd, log, conf.Dom)
	ordCont := controller.NewOrderController(pgTx, pgOrd, streamEve)

	if err := must(http.NewHttp(conf.Http, conf.Dom, log, ordCont)).Listen(); err != nil {
		panic(err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
