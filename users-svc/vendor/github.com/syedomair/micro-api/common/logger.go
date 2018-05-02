package common

import (
	"os"

	gokitlog "github.com/go-kit/kit/log"
)

func GetLogger() gokitlog.Logger {

	var logger gokitlog.Logger
	{
		logger = gokitlog.NewLogfmtLogger(os.Stdout)
		logger = gokitlog.With(logger, "TIME", gokitlog.DefaultTimestamp)
		logger = gokitlog.With(logger, "CALLER", gokitlog.DefaultCaller)
	}
	return logger
}
