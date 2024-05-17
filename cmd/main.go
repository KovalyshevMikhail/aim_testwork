package main

import (
	"fmt"
	"log"
	"sync"

	"aim_testwork/internal/conf"
	"aim_testwork/internal/data"
	"aim_testwork/internal/logger"
	"aim_testwork/internal/sender"

	"github.com/joho/godotenv"
)

func main() {
	// get all environments
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// get all prepared (config + logger)
	cfg := conf.Get()
	log := logger.Get()

	log.Info().
		Str("Config", "load - ok").
		Str("Settings", fmt.Sprintf("%v", cfg)).
		Send()

	// init synchronization
	wg := sync.WaitGroup{}
	pwg := &wg

	log.Info().Msg("application started")

	if cfg.App.FactsGorCount > 1 {
		// if more than 1 needs computation, start it

		for i := 0; i < cfg.App.FactsGorCount; i++ {
			log.Info().Int("Computation", i).Msg("started computation")
			dtos := make(chan data.FactFormDTO, 5)

			pwg.Add(1)
			go data.Gen(cfg.App.FactsGenCount, dtos, pwg)

			pwg.Add(1)
			go sender.SendToServer(dtos, pwg)
		}
	} else {
		// else start only one

		dtos := make(chan data.FactFormDTO, 5)

		pwg.Add(1)
		go data.Gen(cfg.App.FactsGenCount, dtos, pwg)

		pwg.Add(1)
		go sender.SendToServer(dtos, pwg)
	}

	pwg.Wait()

	log.Info().Msg("application exit")
}
