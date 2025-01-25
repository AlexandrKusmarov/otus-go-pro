package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/senderConfig.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.NewConfigSender(configFile)
	logg := logger.New(config.Logger.IsEnabled, config.Logger.Level, "")
	// Проверяем, какое хранилище будет использовано
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host,
		config.Database.Port)

	storage, err := common.NewStorage(ctx, config.Database.IsInMemoryStorage, config.Database.DriverName, dsn)
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		// logg.Error("Error initializing storage:", err)
		return
	}
	defer func() {
		if err := storage.Close(ctx); err != nil {
			fmt.Printf("Error closing storage: %v\n", err)
		}
	}()

	_, err = sender.NewSender(logg, storage, &config, ctx, cancel)
	if err != nil {
		logg.Error("failed to create senderApp: ", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	return
}
