//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/config"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"net/http"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"github.com/stretchr/testify/assert"
)

// cleanupTestData удаляет тестовые данные из базы данных
func cleanupTestData(ctx context.Context, storage common.StorageInterface, eventIDs []int64) {
	for _, eventID := range eventIDs {
		_ = storage.DeleteNotificationByEventId(ctx, eventID)
		err := storage.DeleteEvent(ctx, eventID)
		if err != nil {
			fmt.Printf("Ошибка при удалении тестовых данных с ID %d: %v\n", eventID, err)
		}
	}
}

// SetupTestEnvironment настраивает окружение для интеграционных тестов, включая подключение к базе данных и инициализацию HTTP-сервера.
func SetupTestEnvironment(t *testing.T, configFile string) (context.Context, func(), common.StorageInterface, string) {
	config := config.NewConfig(configFile)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	t.Cleanup(cancel) // Убедитесь, что контекст будет отменен после завершения теста

	// Формируем DSN для подключения к базе данных
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Database.Username, config.Database.Password, config.Database.Database, config.Database.Host, config.Database.Port)

	// Инициализация хранилища
	storage, err := common.NewStorage(ctx, config.Database.IsInMemoryStorage, config.Database.DriverName, dsn)
	if err != nil {
		t.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Убедитесь, что сервер запущен, прежде чем продолжить
	serverAddr := fmt.Sprintf("%s:%s", config.ServerConf.Host, config.ServerConf.Port)

	// Возвращаем контекст, функцию для очистки и хранилище
	return ctx, func() {
			storage.Close(ctx) // Закрываем соединение после завершения теста
		}, storage,
		serverAddr
}

// TestIntegrationCreateEvent выполняет интеграционный тест для создания события
func TestIntegrationCreateEvent(t *testing.T) {
	ctx, cleanup, storage, serverAddr := SetupTestEnvironment(t, "../configs/config.yaml")
	defer cleanup()

	//serverAddr := fmt.Sprintf("%s:%s", "localhost", "8888")

	// Создание события для тестирования
	newEvent := event.Event{
		Title:            "Тестовое событие",
		EventDateTime:    time.Now().Add(-24 * time.Hour), // Событие вчера
		EventEndDateTime: time.Now().Add(-23 * time.Hour), // Событие заканчивается через час
		Description:      "",
		UserID:           1,
	}

	// Кодирование данных в JSON
	eventJSON, err := json.Marshal(newEvent)
	if err != nil {
		t.Fatalf("Ошибка при кодировании события: %v", err)
	}

	// Создание нового запроса
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/event/create", serverAddr), bytes.NewBuffer(eventJSON))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статус-кода
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Проверка, что событие было создано
	var createdEvent event.Event
	err = json.NewDecoder(resp.Body).Decode(&createdEvent)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверка, что созданное событие соответствует отправленному
	assert.Equal(t, newEvent.Title, createdEvent.Title)

	// Очистка данных из базы данных, если это необходимо
	defer cleanupTestData(ctx, storage, []int64{createdEvent.ID}) // Удалите тестовые данные после завершения теста
}

// TestIntegrationCreateEventWithoutDescription выполняет интеграционный тест для создания события без описания,
// что должно запустить процесс отправки события в кафка и обновления поля description на "Добавлено в очередь"
func TestIntegrationCreateEventWithoutDescription(t *testing.T) {
	ctx, cleanup, storage, serverAddr := SetupTestEnvironment(t, "../configs/config.yaml")
	defer cleanup()

	// Создание события для тестирования
	newEvent := event.Event{
		Title:            "Тестовое событие",
		EventDateTime:    time.Now(),
		EventEndDateTime: time.Now().Add(-23 * time.Hour), // Событие заканчивается через час
		Description:      "",
		UserID:           1,
	}

	// Кодирование данных в JSON
	eventJSON, err := json.Marshal(newEvent)
	if err != nil {
		t.Fatalf("Ошибка при кодировании события: %v", err)
	}

	// Создание нового запроса
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/event/create", serverAddr), bytes.NewBuffer(eventJSON))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статус-кода
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Проверка, что событие было создано
	var createdEvent event.Event
	err = json.NewDecoder(resp.Body).Decode(&createdEvent)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверка, что созданное событие соответствует отправленному
	assert.Equal(t, newEvent.Title, createdEvent.Title)

	// Спим 7 секунд, чтобы сервис по отправке сообщения в кафка успел его обработать.
	time.Sleep(7 * time.Second)

	getEvent, err := storage.GetEvent(ctx, createdEvent.ID)
	// Проверка, что созданное событие было отправлено в кафка
	assert.Equal(t, "Добавлено в очередь", getEvent.Description)

	// Очистка данных из базы данных, если это необходимо
	defer cleanupTestData(ctx, storage, []int64{createdEvent.ID}) // Удалите тестовые данные после завершения теста
}

// TestIntegrationCreateMultipleEvents проверяет создание нескольких событий
func TestIntegrationCreateMultipleEvents(t *testing.T) {
	ctx, cleanup, storage, serverAddr := SetupTestEnvironment(t, "../configs/config.yaml")
	defer cleanup()

	// Создание событий для тестирования
	eventDates := []time.Time{
		time.Now(),                    // Сегодня
		time.Now().AddDate(0, 0, -5),  // 5 дней назад
		time.Now().AddDate(0, 0, -20), // 20 дней назад
	}

	var createdEventIDs []int64

	// Выполнение запроса
	client := &http.Client{}

	for _, eventDate := range eventDates {
		newEvent := event.Event{
			Title:            "Тестовое событие",
			EventDateTime:    eventDate,
			EventEndDateTime: eventDate.Add(1 * time.Hour), // Событие заканчивается через час
			Description:      "",
			UserID:           1,
		}

		// Кодирование данных в JSON
		eventJSON, err := json.Marshal(newEvent)
		if err != nil {
			t.Fatalf("Ошибка при кодировании события: %v", err)
		}

		// Создание нового запроса
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/event/create", serverAddr), bytes.NewBuffer(eventJSON))
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Ошибка при выполнении запроса: %v", err)
		}
		defer resp.Body.Close()

		// Проверка статус-кода
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Проверка, что событие было создано
		var createdEvent event.Event
		err = json.NewDecoder(resp.Body).Decode(&createdEvent)
		if err != nil {
			t.Fatalf("Ошибка при декодировании ответа: %v", err)
		}

		// Проверка, что созданное событие соответствует отправленному
		assert.Equal(t, newEvent.Title, createdEvent.Title)

		// Сохраняем ID созданного события для последующей очистки
		createdEventIDs = append(createdEventIDs, createdEvent.ID)

		// Проверка, что событие было создано
		getEvent, err := storage.GetEvent(ctx, createdEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, newEvent.Title, getEvent.Title)
	}

	timeNow := time.Now().Format("2006-01-02")

	// Проверка получения событий за день, неделю и месяц
	for _, duration := range []struct {
		period string
		value  string
	}{
		{"day", timeNow},   // Сегодня
		{"week", timeNow},  // На прошлой неделе
		{"month", timeNow}, // В прошлом месяце
	} {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/event/events/%s/%s", serverAddr, duration.period, duration.value), nil)
		if err != nil {
			t.Fatalf("Ошибка при создании запроса: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Ошибка при выполнении запроса: %v", err)
		}
		defer resp.Body.Close()

		// Проверка статус-кода
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Проверка, что события были получены
		var events []event.Event
		err = json.NewDecoder(resp.Body).Decode(&events)
		if err != nil {
			t.Fatalf("Ошибка при декодировании ответа: %v", err)
		}

		// Проверка, что количество событий соответствует ожидаемому
		assert.Greater(t, len(events), 0, "Должно быть хотя бы одно событие")
	}

	// Очистка данных из базы данных
	defer cleanupTestData(ctx, storage, createdEventIDs) // Удаление тестовых данных после завершения теста
}
