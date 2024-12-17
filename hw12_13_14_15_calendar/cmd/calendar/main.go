package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/configs"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/common"
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := configs.NewConfig(configFile)
	logg := logger.New(config.Logger.IsEnabled, config.Logger.Level, "")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Проверяем, какое хранилище будет использовано
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host, config.Database.Port)

	// Вызываем файл миграции
	if !config.Database.IsInMemoryStorage {
		// Запустим миграцию
		sqlstorage.MigrateData(ctx, config.Database)
	}

	storage, err := common.NewStorage(ctx, config.Database.IsInMemoryStorage, config.Database.DriverName, dsn)
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		//logg.Error("Error initializing storage:", err)
		return
	}
	defer func() {
		if err := storage.Close(ctx); err != nil {
			fmt.Printf("Error closing storage: %v\n", err)
		}
	}()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config.ServerConf.Host, config.ServerConf.Port, calendar, logg)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
