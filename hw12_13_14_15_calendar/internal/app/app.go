package app

import (
	"context"
	"fmt"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App {
	fmt.Println(logger, storage) // антилинтер
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	fmt.Println(ctx, id, title) // антилинтер
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
