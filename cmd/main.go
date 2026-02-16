package main

import (
	"bot/internal/bot"
	"bot/internal/database"
	"bot/internal/database/models"
	"bot/internal/logger"
	"bot/internal/settings"
	"bot/pkg/config"
	"context"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger.New()
	var cfg settings.Config
	config.MustLoad(&cfg, settings.DefaultConfig, logger)

	database := database.MustConnect(cfg.Database, logger)
	database.AutoMigrate(&models.File{})

	bot := bot.New(cfg.Bot, database, logger)
	go bot.Run(ctx)

	logger.Info("Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
