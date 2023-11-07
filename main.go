package main

import (

	// "github.com/guptarohit/asciigraph"
	cmd "github.com/zakisk/dock-stats/internal/cli"
	"github.com/zakisk/dock-stats/pkg/logger"
)

func main() {
	log := logger.NewLog()
	err := cmd.Execute()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
}
