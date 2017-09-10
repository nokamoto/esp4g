package utils

import (
	"go.uber.org/zap"
	"log"
)

var Logger *zap.SugaredLogger

func init() {
	production, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	Logger = production.Sugar()
}
