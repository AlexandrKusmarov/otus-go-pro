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
	Host string `json:"host"`
	Port string `json:"port"`
}

type RMQ struct {
	URI       string `json:"uri"`
	ReConnect struct {
		MaxElapsedTime  string  `json:"maxElapsedTime"`
		InitialInterval string  `json:"initialInterval"`
		Multiplier      float64 `json:"multiplier"`
		MaxInterval     string  `json:"maxInterval"`
	}
}

type Binding struct {
	ExchangeName string `json:"exchangeName"`
	ExchangeType string `json:"exchangeType"`
	QueueName    string `json:"queueName"`
	BindingKey   string `json:"bindingKey"` // Message routing rules
}
