package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"log"

	"github.com/pressly/goose/v3"
)

func MigrateData(_ context.Context, config config.DatabaseConf) error {
	// Подключение к серверу базы данных (без указания конкретной базы данных)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer db.Close()

	// Проверка и создание базы данных, если она не существует
	_, err = db.Exec("CREATE DATABASE " + config.Database)
	if err != nil && err.Error() != "pq: database \""+config.Database+"\" already exists" {
		log.Printf("Ошибка при создании базы данных: %v", err)
		return err
	}

	// Теперь подключаемся к созданной базе данных
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.Sslmode)

	db, err = sql.Open(config.DriverName, dsn)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return err
	}
	defer db.Close()

	// Вызываем файл миграции
	if !config.IsInMemoryStorage {
		// Запустим миграцию
		if err := goose.Up(db, config.MigrationPath); err != nil {
			log.Printf("Ошибка при выполнении миграции1111: %v", err)
			log.Print("Путь к файлу: " + config.MigrationPath)
			return err
		}
		fmt.Println("Миграция выполнена успешно.")
	}
	return nil
}
