package utils

import "go.uber.org/zap"

var Log *zap.Logger

func InitLogger(env string) {
	var err error

	if env == "prod" {
		Log, err = zap.NewProduction()
	} else {
		Log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}
}
