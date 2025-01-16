package sender

import (
	"encoding/json"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/kaf"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
)

type App struct { // TODO
	storage common.StorageInterface
}

type Logger interface { // TODO
}

type Storage interface { // TODO

}

func NewSender(log logger.Log, storage common.StorageInterface, conf *config.SenderConfig) (*App, error) {
	// Инициализация rmq
	//client, err := rmq.NewRMQClient(conf.RMQ.URI)

	//if err != nil {
	//	log.Error("Ошибка создания подключения к Rabbit MQ, Error: ", err)
	//	return nil, err // Важно вернуть nil, если произошла ошибка
	//}
	//defer client.Close()

	// Создаем Kafka подписчика

	kafkaConsumer, err := kaf.NewKafkaConsumer(conf.KafkaConf.Broker, conf.KafkaConf.Consumer.GroupID, conf.KafkaConf.Consumer.Topic)
	if err != nil {
		log.Error("Failed to create Kafka consumer: ", err)
	}
	defer kafkaConsumer.Close()

	// Определяем обработчик сообщений
	//handler := func(msgBody []byte) {
	//	log.Info("Получено сообщение:", string(msgBody))
	//	// Пример: десериализация и сохранение в storage
	//	var notification scheduler.Notification
	//	if err := json.Unmarshal(msgBody, &notification); err != nil {
	//		log.Error("Ошибка десериализации сообщения:", err)
	//		return
	//	}
	//
	//	// Отправляем сообщение из Kafka пользователю.
	//	log.Info("Сообщение обработано и отправлено пользователю")
	//}
	//
	//// Подписываемся на очередь с обработчиком
	//if err := client.Consume(conf.Binding.QueueName, conf.Binding.ExchangeName, handler); err != nil {
	//	log.Error("Ошибка подписки на очередь:", err)
	//	return nil, err // Важно вернуть nil, если произошла ошибка
	//}

	//Обработчик кафка
	kafkaConsumer.Consume(func(value []byte) {
		// Обработка сообщения
		log.Info("Получено сообщение:", string(value))
		// Пример: десериализация и сохранение в storage
		var notification scheduler.Notification
		if err := json.Unmarshal(value, &notification); err != nil {
			log.Error("Ошибка десериализации сообщения:", err)
			return
		}

		// Отправляем сообщение из Kafka пользователю.
		log.Info("Сообщение обработано и отправлено пользователю")
	})

	fmt.Println(log, storage) // антилинтер
	return &App{storage: storage}, err
}
