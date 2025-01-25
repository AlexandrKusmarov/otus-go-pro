package sqlstorage

import (
	"context"
	"fmt"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/scheduler"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт PostgreSQL-драйвера
)

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) GetAllEventsForDay(ctx context.Context, day time.Time) ([]event.Event, error) {
	var events []event.Event
	startOfDay := day.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	//query := `SELECT * FROM public.events WHERE event_date_time >= $1 AND event_date_time < $2`
	query := `SELECT * FROM public.events WHERE event_date_time >= $1 AND event_date_time < $2`
	if err := s.db.SelectContext(ctx, &events, query, startOfDay, endOfDay); err != nil {
		return nil, fmt.Errorf("failed to get events for day: %w", err)
	}
	return events, nil
}

func (s *Storage) GetAllEventsForWeek(ctx context.Context, startDayOfWeek time.Time) ([]event.Event, error) {
	var events []event.Event
	startOfWeek := startDayOfWeek.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	query := `SELECT * FROM public.events 
              WHERE event_date_time >= $1 AND event_date_time < $2`
	if err := s.db.SelectContext(ctx, &events, query, startOfWeek, endOfWeek); err != nil {
		return nil, fmt.Errorf("failed to get events for week: %w", err)
	}
	return events, nil
}

func (s *Storage) GetAllEventsForMonth(ctx context.Context, startDayOfMonth time.Time) ([]event.Event, error) {
	var events []event.Event
	startOfMonth := time.Date(startDayOfMonth.Year(), startDayOfMonth.Month(), 1, 0, 0, 0, 0, startDayOfMonth.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // Переход к следующему месяцу

	query := `SELECT * FROM public.events 
              WHERE event_date_time >= $1 AND event_date_time < $2`
	if err := s.db.SelectContext(ctx, &events, query, startOfMonth, endOfMonth); err != nil {
		return nil, fmt.Errorf("failed to get events for month: %w", err)
	}
	return events, nil
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, driverName string, dsn string) error {
	var err error
	s.db, err = sqlx.ConnectContext(ctx, driverName, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

// Notification представляет структуру уведомления.
type Notification struct {
	ID            int64  `db:"id"`
	EventID       int64  `db:"event_id"`
	Title         string `db:"title"`
	EventDateTime string `db:"event_date_time"` // Используйте time.Time для работы с временем
	UserID        int64  `db:"user_id"`
}

// CreateEvent добавляет новое событие в базу данных.
func (s *Storage) CreateEvent(ctx context.Context, event *event.Event) error {
	query := `INSERT INTO public.events 
    (title, event_date_time, event_end_date_time, description, user_id, notify_before_event) 
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	_, err := s.db.NamedExecContext(ctx, query, event)

	return err
}

// GetEvent возвращает событие по его ID.
func (s *Storage) GetEvent(ctx context.Context, id int64) (*event.Event, error) {
	var event event.Event
	query := `SELECT * FROM public.events WHERE id = $1`

	if err := s.db.GetContext(ctx, &event, query, id); err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	return &event, nil
}

// UpdateEvent обновляет существующее событие.
func (s *Storage) UpdateEvent(ctx context.Context, event *event.Event) error {
	query := `UPDATE public.events SET 
              title = $1, 
              event_date_time = $2, 
              event_end_date_time = $3, 
              description = $4, 
              user_id = $5, 
              notify_before_event = $6 
              WHERE id = $7`

	_, err := s.db.ExecContext(ctx, query,
		event.Title,
		event.EventDateTime,
		event.EventEndDateTime,
		event.Description,
		event.UserID,
		event.NotifyBeforeEvent,
		event.ID) // Передаем ID события для условия WHERE
	return err
}

// DeleteEvent удаляет событие по его ID.
func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	query := `DELETE FROM public.events WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// GetAllEvents возвращает все события.
func (s *Storage) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	var events []event.Event
	query := `SELECT * FROM public.events`

	if err := s.db.SelectContext(ctx, &events, query); err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

// CreateNotification добавляет новое уведомление в базу данных.
func (s *Storage) CreateNotification(ctx context.Context, notification scheduler.Notification) (scheduler.Notification, error) {
	query := `INSERT INTO public.notification (event_id, title, event_date_time, user_id) 
              VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(ctx, query, notification.EventID, notification.Title, notification.EventDateTime, notification.UserID)
	if err != nil {
		return scheduler.Notification{}, err
	}

	return notification, nil
}

// GetNotification возвращает уведомление по его ID.
func (s *Storage) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	var notification Notification
	query := `SELECT * FROM public.notification WHERE id = $1`

	if err := s.db.GetContext(ctx, &notification, query, id); err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return &notification, nil
}

// UpdateNotification обновляет существующее уведомление.
func (s *Storage) UpdateNotification(ctx context.Context, notification *Notification) error {
	query := `UPDATE public.notification SET 
              event_id = $1, 
              title = $2, 
              event_date_time = $3, 
              user_id = $4 
              WHERE id = $5`

	_, err := s.db.ExecContext(ctx, query,
		notification.EventID,
		notification.Title,
		notification.EventDateTime,
		notification.UserID,
		notification.ID) // Передаем ID уведомления для условия WHERE
	return err
}

// DeleteNotification удаляет уведомление по его ID.
func (s *Storage) DeleteNotification(ctx context.Context, id int64) error {
	query := `DELETE FROM public.notification WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// GetAllNotifications возвращает все уведомления.
func (s *Storage) GetAllNotifications(ctx context.Context) ([]Notification, error) {
	var notifications []Notification
	query := `SELECT * FROM public.notification`

	if err := s.db.SelectContext(ctx, &notifications, query); err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifications, nil
}
