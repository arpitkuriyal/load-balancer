package logger

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
		panic("falide to initialize logger" + err.Error())
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
