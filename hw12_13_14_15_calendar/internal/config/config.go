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
