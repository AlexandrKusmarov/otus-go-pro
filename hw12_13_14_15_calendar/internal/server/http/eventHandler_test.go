package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockApp - это мок для приложения, которое реализует интерфейс, используемый в Server.
type MockApp struct {
	mock.Mock
}

func (m *MockApp) CreateEvent(ctx context.Context, event *event.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApp) UpdateEvent(ctx context.Context, event *event.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockApp) DeleteEvent(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApp) GetAllEvents(ctx context.Context) ([]event.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]event.Event), args.Error(1)
}

func (m *MockApp) GetEvent(ctx context.Context, id int64) (*event.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*event.Event), args.Error(1)
}

func TestCreateEventHandler(t *testing.T) {
	mockApp := new(MockApp)
	server := &Server{app: mockApp}
	router := mux.NewRouter()
	router.HandleFunc("/events", server.createEventHandler).Methods("POST")

	eventToCreate := event.Event{Title: "New Event"}
	mockApp.On("CreateEvent", mock.Anything, &eventToCreate).Return(nil)

	body, _ := json.Marshal(eventToCreate)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockApp.AssertExpectations(t)
}

func TestUpdateEventHandler(t *testing.T) {
	mockApp := new(MockApp)
	server := &Server{app: mockApp}
	router := mux.NewRouter()
	router.HandleFunc("/events/{eventId}", server.updateEventHandler).Methods("PUT")

	eventToUpdate := event.Event{ID: 1, Title: "Updated Event"}
	mockApp.On("UpdateEvent", mock.Anything, &eventToUpdate).Return(nil)

	body, _ := json.Marshal(eventToUpdate)
	req := httptest.NewRequest("PUT", "/events/1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockApp.AssertExpectations(t)
}

func TestDeleteEventHandler(t *testing.T) {
	mockApp := new(MockApp)
	server := &Server{app: mockApp}
	router := mux.NewRouter()
	router.HandleFunc("/events/{eventId}", server.deleteEventHandler).Methods("DELETE")

	mockApp.On("DeleteEvent", mock.Anything, int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/events/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockApp.AssertExpectations(t)
}

func TestGetAllEventsHandler(t *testing.T) {
	mockApp := new(MockApp)
	server := &Server{app: mockApp}
	router := mux.NewRouter()
	router.HandleFunc("/events", server.getAllEventsHandler).Methods("GET")

	events := []event.Event{{ID: 1, Title: "Event 1"}, {ID: 2, Title: "Event 2"}}
	mockApp.On("GetAllEvents", mock.Anything).Return(events, nil)

	req := httptest.NewRequest("GET", "/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Проверка наличия полей в JSON
	var responseArray []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseArray)
	assert.NoError(t, err)

	// Проверка наличия полей ID и Title в каждом объекте
	for _, eventMap := range responseArray {
		assert.Contains(t, eventMap, "ID")
		assert.Contains(t, eventMap, "Title")
	}

	// Проверка значений
	assert.Equal(t, float64(1), responseArray[0]["ID"]) // Значение ID будет float64
	assert.Equal(t, "Event 1", responseArray[0]["Title"])
	assert.Equal(t, float64(2), responseArray[1]["ID"])
	assert.Equal(t, "Event 2", responseArray[1]["Title"])
	mockApp.AssertExpectations(t)
}

func TestGetEventHandler(t *testing.T) {
	mockApp := new(MockApp)
	server := &Server{app: mockApp}
	router := mux.NewRouter()
	router.HandleFunc("/events/{eventId}", server.getEventHandler).Methods("GET")

	eventToGet := &event.Event{ID: 1, Title: "Event 1"}
	mockApp.On("GetEvent", mock.Anything, int64(1)).Return(eventToGet, nil)

	req := httptest.NewRequest("GET", "/events/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Проверка конкретных полей из JSON
	var responseMap map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseMap)
	assert.NoError(t, err)

	// Проверка наличия полей ID и Title
	assert.Contains(t, responseMap, "ID")
	assert.Contains(t, responseMap, "Title")

	// Проверка значений
	assert.Equal(t, float64(1), responseMap["ID"]) // Значение ID будет float64
	assert.Equal(t, "Event 1", responseMap["Title"])
	mockApp.AssertExpectations(t)
}
