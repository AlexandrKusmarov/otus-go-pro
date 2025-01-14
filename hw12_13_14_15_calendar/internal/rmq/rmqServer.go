package rmq

import (
	"github.com/streadway/amqp"
)

// RMQClient представляет клиент для работы с RabbitMQ
type RMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRMQClient создает новое соединение с RabbitMQ
func NewRMQClient(url string) (*RMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RMQClient{
		connection: conn,
		channel:    ch,
	}, nil
}

// Publish отправляет сообщение в очередь
func (c *RMQClient) Publish(queueName string, exchange string, body []byte) error {
	// Объявляем очередь
	_, err := c.channel.QueueDeclare(
		queueName, // имя очереди
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Отправляем сообщение
	err = c.channel.Publish(
		exchange,  // exchange
		queueName, // routing key (имя очереди)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

// Consume начинает прослушивание очереди
func (c *RMQClient) Consume(queueName string, exchange string, handler func([]byte)) error {
	msgs, err := c.channel.Consume(
		queueName, // имя очереди
		exchange,  // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()

	return nil
}

// Close закрывает соединение и канал
func (c *RMQClient) Close() {
	c.channel.Close()
	c.connection.Close()
}
