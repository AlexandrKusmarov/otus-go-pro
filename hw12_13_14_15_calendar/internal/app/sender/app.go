package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/kaf"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
)

type App struct { // TODO
	storage common.StorageInterface
}

type Logger interface { // TODO
}

type Storage interface { // TODO

}

func NewSender(log logger.Log, storage common.StorageInterface, conf *config.SenderConfig, ctx context.Context,
	cancel context.CancelFunc) (*App, error) {

	app := &App{storage: storage}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel() // Отменяем контекст при получении сигнала
	}()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{conf.KafkaConf.Broker},
		Topic:   conf.KafkaConf.Consumer.Topic,
		GroupID: conf.KafkaConf.Consumer.GroupID,
	})

	// Запускаем бесконечный цикл для прослушивания сообщений
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Завершение работы горутины прослушивания сообщений")
				return // Выходим из горутины при отмене контекста
			default:
				message, err := kaf.ConsumeMessage(log, ctx, reader)
				if err != nil {
					log.Error("Ошибка при получении сообщения:", err)
					continue // Пропускаем итерацию в случае ошибки
				}

				var notification scheduler.Notification
				if err = json.Unmarshal(message.Value, &notification); err != nil {
					log.Error("Ошибка десериализации сообщения:", err)
					continue // Пропускаем итерацию в случае ошибки
				}

				// Имитация отправки сообщения пользователю
				fmt.Println("Сообщение отправлено пользователю: ", string(message.Value))
				log.Info("Сообщение отправлено пользователю: " + string(message.Value))
			}
		}
	}()

	return app, nil
}
