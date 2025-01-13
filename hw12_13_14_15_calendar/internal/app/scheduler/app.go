package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/rmq"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
	"time"
)

type AppScheduler struct { // TODO
	storage common.StorageInterface
}

type Storage interface { // TODO

}

func NewScheduler(log logger.Log, storage common.StorageInterface, ctx context.Context, conf *config.SchedulerConfig) (*AppScheduler, error) {
	fmt.Println(log, storage) // антилинтер

	// Инициализация rmq
	client, err := rmq.NewRMQClient(conf.RMQ.URI)

	if err != nil {
		log.Error("Ошибка создания подключения к Rabbit MQ, Error: ", err)
	}
	defer client.Close()

	startEventPublisher(log, ctx, storage, *client, conf)

	return &AppScheduler{storage: storage}, nil
}

// Функция для запуска горутины
func startEventPublisher(log logger.Log, ctx context.Context, storage common.StorageInterface, client rmq.RMQClient, conf *config.SchedulerConfig) {
	ticker := time.NewTicker(time.Minute) // Тикер, который срабатывает раз в час
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			allEvents, err := storage.GetAllEventsForDay(ctx, time.Now())
			if err != nil {
				log.Error("Error:", err)
				continue
			}

			for _, event := range allEvents {
				// Публикуем только ивенты, где поле description = nil или пустая строка
				if event.Description == "" || event.Description == "nil" {
					notification, err := storage.CreateNotification(ctx, scheduler.Notification{
						EventID:       event.ID,
						Title:         event.Title,
						EventDateTime: event.EventDateTime,
						UserID:        event.UserID,
					})

					jsonData, err := json.Marshal(notification)
					if err != nil {
						log.Error("Error marshalling notification to JSON:", err)
						continue
					}

					// Публикуем JSON notification
					log.Info("Публикуем JSON:", string(jsonData))
					client.Publish(conf.Binding.QueueName, jsonData)

					// Обновляем описание события
					event.Description = "Добавлено в очередь"
					err = storage.UpdateEvent(ctx, &event)
					if err != nil {
						log.Info("Error updating event description:", err)
					}
				}
			}
		case <-ctx.Done():
			log.Info("Горутина завершена.")
			return
		}
	}
}
