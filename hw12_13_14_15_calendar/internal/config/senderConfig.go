package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type SenderConfig struct {
	Logger   LoggerConf   `yaml:"logger"`
	Database DatabaseConf `yaml:"database"`
	RMQ      RMQ          `yaml:"rmq"`
	Binding  Binding      `yaml:"binding"`
	// TODO Возможно что-то еще понадобится
}

func NewConfigSender(pathConfigFile string) SenderConfig {
	// Открываем файл
	file, err := os.Open(pathConfigFile)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	var config SenderConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		file.Close()
		log.Fatalf("Ошибка декодирования YAML: %v", err) //nolint
	}

	return config
}
