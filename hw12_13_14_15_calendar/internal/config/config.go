package config

type LoggerConf struct {
	Level     string `yaml:"level"`
	IsEnabled bool   `yaml:"enabled"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConf struct {
	DriverName        string `yaml:"driver"`
	Host              string `yaml:"host"`
	Port              string `yaml:"port"`
	Username          string `yaml:"user"`
	Password          string `yaml:"password"`
	Database          string `yaml:"dbname"`
	Sslmode           string `yaml:"sslmode"`
	IsInMemoryStorage bool   `yaml:"isInMemoryStorage"`
	MigrationPath     string `yaml:"migrationPath"`
}

type GRPCServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ReConnect struct {
	MaxElapsedTime  string  `yaml:"maxElapsedTime"`
	InitialInterval string  `yaml:"initialInterval"`
	Multiplier      float64 `yaml:"multiplier"`
	MaxInterval     string  `yaml:"maxInterval"`
}

type RMQ struct {
	URI       string    `yaml:"uri"`
	ReConnect ReConnect `yaml:"reConnect"`
}

type Binding struct {
	ExchangeName string `yaml:"exchangeName"`
	ExchangeType string `yaml:"exchangeType"`
	QueueName    string `yaml:"queueName"`
	BindingKey   string `yaml:"bindingKey"` // Message routing rules
}

// Kafka конфигурация
type KafkaConf struct {
	Broker   string        `yaml:"broker"`   // Адрес брокера Kafka
	Consumer KafkaConsumer `yaml:"consumer"` // Конфигурация потребителя
	Producer KafkaProducer `yaml:"producer"` // Конфигурация продюсера
}

type KafkaConsumer struct {
	GroupID         string `yaml:"groupID"`         // ID группы для потребителей
	Topic           string `yaml:"topic"`           // Топик, на который подписывается потребитель
	AutoOffsetReset string `yaml:"autoOffsetReset"` // Параметр для управления смещением
	MaxPollRecords  int    `yaml:"maxPollRecords"`  // Максимальное количество записей для выборки за один раз
	Threads         int    `yaml:"threads"`         // Количество потоков для обработки сообщений
}

type KafkaProducer struct {
	Topic string `yaml:"topic"` // Топик, в который будет отправляться сообщение
}
