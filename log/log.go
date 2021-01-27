package log

import (
	"log"

	"go.uber.org/zap"
)

/*
Zap is a lightweight verbose logging tool that needs limited
configurattion.
- Transfering all the logging output to utilize zap
- Need to figure ou to to initialize logger before all tests
	-Maybe a setup function?
- Can potentially move to a custom logging solution in the future
*/

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
