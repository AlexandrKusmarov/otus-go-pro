package kaf

import (
	"context"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/logger"
	"github.com/segmentio/kafka-go"
	"log"
)

// KafkaClient представляет клиент для работы с Kafka
type KafkaClient struct {
	writer *kafka.Writer
}

func NewKafkaClient(brokerAddress string, topic string) *KafkaClient {
	return &KafkaClient{
		kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddress},
			Topic:   topic,
		}),
	}
}

// Создание топика (если его еще нет)
func CreateTopic(brokerAddress string, topic string) {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatalf("Ошибка подключения к Kafka: %v", err)
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		log.Fatalf("Ошибка создания топика: %v", err)
	}
	fmt.Println("Топик успешно создан!")
}

func (c *KafkaClient) Publish(key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	return c.writer.WriteMessages(context.Background(), msg)
}

func (c *KafkaClient) Close() error {
	return c.writer.Close()
}

// Чтение из кафка
func ConsumeMessage(log logger.Log, ctx context.Context, reader *kafka.Reader) (kafka.Message, error) {
	fmt.Println("Ожидание сообщений...")
	message, err := reader.ReadMessage(ctx)
	if err != nil {
		log.Error("Ошибка чтения сообщения: ", err)
		return kafka.Message{}, err
	}
	log.Info("Получено сообщение: ", string(message.Value))
	return message, nil
}
