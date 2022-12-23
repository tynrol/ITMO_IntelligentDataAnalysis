package main

import (
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/config"
	app "github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/internal"
	"log"
	"os"
)

func main() {
	var (
		conf config.Config
		quit = make(chan os.Signal, 1)
	)

	if err := conf.Parse(); err != nil {
		log.Fatalf("error parsing config: %s", err)
	}

	application := app.NewApp(&conf)

	go application.Start()

	<-quit

	application.Close()
}
