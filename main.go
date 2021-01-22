package main

import (
	"github.com/bdkiran/nolan/cmd"
	"github.com/bdkiran/nolan/log"
)

func main() {
	log.Initialize()
	cmd.Execute()
}
