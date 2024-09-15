package main

import (
	_ "embed"

	"github.com/clicklord/lms/bootstrap"
	"github.com/clicklord/lms/config"
	"github.com/clicklord/lms/log"
	"github.com/clicklord/lms/rrcache"
)

//go:embed "data/movie-icon.png"
var defaultIcon []byte

func main() {
	logger := log.New()
	logger.SetAsDefault()

	cfg, err := bootstrap.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}
	cache := &config.FFprobeCache{
		Cache: rrcache.New(64 << 20),
	}
	if err := cache.Load(cfg.FFprobeCachePath); err != nil {
		log.Print(err)
	}

	mainWindow := bootstrap.LoadMainWindow(cfg)
	dmsServer := bootstrap.LoadDMS(cfg, cache, defaultIcon, logger)

	if err := dmsServer.Init(); err != nil {
		log.Fatalf("error initing dms server: %v", err)
	}
	go func() {
		if err := dmsServer.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	mainWindow.ShowAndRun()

	err = dmsServer.Close()
	if err != nil {
		log.Fatal(err)
	}
	if err := cache.Save(cfg.FFprobeCachePath); err != nil {
		log.Print(err)
	}
}
