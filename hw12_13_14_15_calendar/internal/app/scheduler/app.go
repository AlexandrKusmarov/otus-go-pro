package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/kaf"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
	"time"
)

type AppScheduler struct { // TODO
	storage common.StorageInterface
	cancel  context.CancelFunc
}

type Storage interface { // TODO

}

func NewScheduler(log logger.Log, storage common.StorageInterface, ctx context.Context, conf *config.SchedulerConfig) (*AppScheduler, error) {
	fmt.Println(log, storage) // антилинтер

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Инициализация rmq
	//client, err := rmq.NewRMQClient(conf.RMQ.URI)

	//Инициализация kafka
	kaf.CreateTopic(conf.KafkaConf.Broker, conf.KafkaConf.Producer.Topic)

	kafkaClient := kaf.NewKafkaClient(conf.KafkaConf.Broker, conf.KafkaConf.Producer.Topic)

	//if err != nil {
	//	log.Error("Ошибка создания подключения к Rabbit MQ, Error: ", err)
	//}
	//defer client.Close()

	// Удаляет события, которые хранятся больше 1 года
	go startDailyCleaner(log, ctx, storage)

	go startEventPublisher(log, ctx, storage, *kafkaClient, conf)

	return &AppScheduler{storage: storage, cancel: cancel}, nil
}

// Функция для запуска горутины
func startEventPublisher(log logger.Log, ctx context.Context, storage common.StorageInterface, kafkaClient kaf.KafkaClient, conf *config.SchedulerConfig) {
	ticker := time.NewTicker(5 * time.Second) // Тикер, который срабатывает раз в час
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

					err = kafkaClient.Publish("10", jsonData)

					if err != nil {
						log.Error("Ошибка публикации в kafka", err)
						fmt.Printf("Ошибка публикации в kafka %v", err)
						continue
					}

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
			kafkaClient.Close()
			return
		}
	}
}

func eventCleaner(log logger.Log, ctx context.Context, storage common.StorageInterface) {
	log.Info("Сервис по очистке старых событий запущен")

	events, err := storage.GetAllEventsForDay(ctx, time.Now().AddDate(-1, 0, 0)) // Получаем события, которые были созданы больше года назад
	if err != nil {
		log.Error("Ошибка получения событий:", err)
		return
	}

	// Удаляем старые события
	for _, event := range events {
		err := storage.DeleteEvent(ctx, event.ID) // Удаляем событие по ID
		if err != nil {
			log.Error("Ошибка удаления события:", err)
		} else {
			log.Info("Событие удалено:", event.ID)
		}
	}
}
func startDailyCleaner(log logger.Log, ctx context.Context, storage common.StorageInterface) {
	ticker := time.NewTicker(24 * time.Hour) // Создаем тикер с интервалом 24 часа
	defer ticker.Stop()

	// Запускаем сразу при старте
	eventCleaner(log, ctx, storage)

	for {
		select {
		case <-ticker.C: // Ждем каждую 24 часа
			eventCleaner(log, ctx, storage)
		case <-ctx.Done(): // Проверяем контекст на завершение
			log.Info("Остановка планировщика очистки событий")
			return
		}
	}
}

// Метод для остановки планировщика
func (s *AppScheduler) Stop() {
	s.cancel() // Вызываем функцию отмены контекста
}
