package kaf

import (
	"context"
	"fmt"
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

// NewKafkaConsumer создает нового потребителя для Kafka
//func NewKafkaConsumer(broker string, groupID string, topic string) (*KafkaClient, error) {
//	// Создаем потребителя
//	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers": broker,
//		"group.id":          groupID,
//		"auto.offset.reset": "earliest",
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	// Подписываемся на топик
//	err = consumer.Subscribe(topic, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return &KafkaClient{
//		consumer: consumer,
//	}, nil
//}

// Consume начинает прослушивание топика
//func (c *KafkaClient) Consume(handler func([]byte)) {
//	go func() {
//		for {
//			msg, err := c.consumer.ReadMessage(-1)
//			if err == nil {
//				handler(msg.Value)
//			} else {
//				log.Printf("Error while consuming message: %v", err)
//			}
//		}
//	}()
//}
