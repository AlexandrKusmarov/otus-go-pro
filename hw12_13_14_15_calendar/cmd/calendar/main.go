package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/grpc_local"
	eventpb "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/grpc_local/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/configs"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/app"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	sqlstorage "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string
var wg sync.WaitGroup

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	//if flag.Arg(0) == "version" {
	//	printVersion()
	//	return
	//}

	config := configs.NewConfig(configFile)
	logg := logger.New(config.Logger.IsEnabled, config.Logger.Level, "")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// Проверяем, какое хранилище будет использовано
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host,
		config.Database.Port)

	// Вызываем файл миграции
	if !config.Database.IsInMemoryStorage {
		// Запустим миграцию
		err := sqlstorage.MigrateData(ctx, config.Database)
		if err != nil {
			os.Exit(1) //nolint
		}
	}

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

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config.ServerConf.Host, config.ServerConf.Port, calendar, logg)

	wg.Add(2)

	// Инициализация gRPC-сервера.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	eventService := grpc_local.NewGrpcServer(&storage)
	eventpb.RegisterEventServiceServer(s, eventService)
	reflection.Register(s)

	// Горутина для остановки gRPC сервера.
	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(ctx, time.Second*3)
		defer cancel()

		s.GracefulStop() // Корректная остановка gRPC сервера
		log.Println("gRPC server stopped gracefully")
	}()

	log.Println("gRPC server is running on port :50051")

	// Запуск gRPC сервера в горутине.
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

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
		os.Exit(1)
	}

	// Блокировка основного потока, чтобы программа не завершилась немедленно.
	//select {}

}
