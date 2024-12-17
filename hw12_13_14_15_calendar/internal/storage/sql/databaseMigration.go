package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/configs"
	"github.com/pressly/goose/v3"
	"log"
)

func MigrateData(ctx context.Context, config configs.DatabaseConf) {

	// Подключение к серверу базы данных (без указания конкретной базы данных)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка и создание базы данных, если она не существует
	_, err = db.Exec("CREATE DATABASE " + config.Database)
	if err != nil && err.Error() != "pq: database \""+config.Database+"\" already exists" {
		log.Fatalf("Ошибка при создании базы данных: %v", err)
	}

	// Теперь подключаемся к созданной базе данных
	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.Database, config.Sslmode)

	db, err = sql.Open(config.DriverName, dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Вызываем файл миграции
	if !config.IsInMemoryStorage {
		// Запустим миграцию
		if err := goose.Up(db, config.MigrationPath); err != nil {
			log.Fatalf("Ошибка при выполнении миграции: %v", err)
		}
		fmt.Println("Миграция выполнена успешно.")
	}
}
