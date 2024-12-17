package configs

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf   `yaml:"logger"`
	Database   DatabaseConf `yaml:"database"`
	ServerConf ServerConf   `yaml:"server"`
	// TODO
}

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
		log.Fatalf("Ошибка декодирования YAML: %v", err)
	}

	return config
}

// TODO
