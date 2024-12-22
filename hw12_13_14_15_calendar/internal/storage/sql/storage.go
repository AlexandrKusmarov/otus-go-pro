package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт PostgreSQL-драйвера
)

type Storage struct {
	db *sqlx.DB
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

type Event struct {
	ID                int64  `db:"id"`
	Title             string `db:"title"`
	EventDateTime     string `db:"event_date_time"` // Используйте time.Time для работы с временем
	EventEndDateTime  string `db:"event_end_date_time"`
	Description       string `db:"description"`
	UserID            int64  `db:"user_id"`
	NotifyBeforeEvent string `db:"notify_before_event"`
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
func (s *Storage) CreateEvent(ctx context.Context, event *Event) error {
	query := `INSERT INTO public.events 
    (title, event_date_time, event_end_date_time, description, user_id, notify_before_event) 
    VALUES (:title, :event_date_time, :event_end_date_time, :description, :user_id, :notify_before_event) RETURNING id`

	return s.db.GetContext(ctx, &event.ID, query, event)
}

// GetEvent возвращает событие по его ID.
func (s *Storage) GetEvent(ctx context.Context, id int64) (*Event, error) {
	var event Event
	query := `SELECT * FROM public.events WHERE id = $1`

	if err := s.db.GetContext(ctx, &event, query, id); err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	return &event, nil
}

// UpdateEvent обновляет существующее событие.
func (s *Storage) UpdateEvent(ctx context.Context, event *Event) error {
	query := `UPDATE public.events SET title = :title, event_date_time = :event_date_time, 
              event_end_date_time = :event_end_date_time, description = :description, 
              user_id = :user_id, notify_before_event = :notify_before_event 
              WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, query, event)
	return err
}

// DeleteEvent удаляет событие по его ID.
func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	query := `DELETE FROM public.events WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

// GetAllEvents возвращает все события.
func (s *Storage) GetAllEvents(ctx context.Context) ([]Event, error) {
	var events []Event
	query := `SELECT * FROM public.events`

	if err := s.db.SelectContext(ctx, &events, query); err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

// CreateNotification добавляет новое уведомление в базу данных.
func (s *Storage) CreateNotification(ctx context.Context, notification *Notification) error {
	query := `INSERT INTO public.notification (event_id, title, event_date_time, user_id) 
              VALUES (:event_id, :title, :event_date_time, :user_id) RETURNING id`

	return s.db.GetContext(ctx, &notification.ID, query, notification)
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
	query := `UPDATE public.notification SET event_id = :event_id, title = :title, 
              event_date_time = :event_date_time, user_id = :user_id 
              WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, query, notification)
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
