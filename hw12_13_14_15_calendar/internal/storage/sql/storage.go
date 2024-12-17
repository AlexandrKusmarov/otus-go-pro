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

func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
