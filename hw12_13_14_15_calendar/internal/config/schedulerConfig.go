package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type SchedulerConfig struct {
	Logger    LoggerConf   `yaml:"logger"`
	Database  DatabaseConf `yaml:"database"`
	RMQ       RMQ          `yaml:"rmq"`
	KafkaConf KafkaConf    `yaml:"kafka"`
	Binding   Binding      `yaml:"binding"`
}

func NewConfigScheduler(pathConfigFile string) SchedulerConfig {
	// Открываем файл
	file, err := os.Open(pathConfigFile)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	var config SchedulerConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		file.Close()
		log.Fatalf("Ошибка декодирования YAML: %v", err) //nolint
	}

	return config
}
