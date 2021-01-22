package log

import (
	"log"

	"go.uber.org/zap"
)

//Initialize sets up a zap global logger
//Allow for this to be configurable for development
//and production use for max performance.
func Initialize() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Unable to Start logger")
	}
	defer logger.Sync()
	//sugar := logger.Sugar()
	zap.ReplaceGlobals(logger)
}
