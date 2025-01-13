package calendar_scheduler

import (
	"context"
	"flag"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/scheduler_config.yaml", "Path to scheduler_config configuration file")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.TODO())
	// ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	config := config.NewConfigScheduler(configFile)
	logg := logger.New(config.Logger.IsEnabled, config.Logger.Level, "")
	// Проверяем, какое хранилище будет использовано
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host,
		config.Database.Port)

	storage, err := common.NewStorage(ctx, config.Database.IsInMemoryStorage, config.Database.DriverName, dsn)
	if err != nil {
		logg.Error("Error initializing storage:", err)
		return
	}
	defer func() {
		if err := storage.Close(ctx); err != nil {
			logg.Error("Error closing storage", err)
		}
	}()

	_, err = scheduler.NewScheduler(logg, storage, ctx, &config)
	if err != nil {
		logg.Error("failed to create schedulerApp: %w", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// <-ctx.Done()

	return
}
