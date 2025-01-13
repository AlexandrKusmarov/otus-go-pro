package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Logger         LoggerConf     `yaml:"logger"`
	Database       DatabaseConf   `yaml:"database"`
	ServerConf     ServerConf     `yaml:"server"`
	GRPCServerConf GRPCServerConf `yaml:"grpc"`
}

func NewConfig(pathConfigFile string) Config {
	// Открываем файл
	file, err := os.Open(pathConfigFile)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		file.Close()
		log.Fatalf("Ошибка декодирования YAML: %v", err) //nolint
	}

	return config
}
