package internalhttp

import (
	"encoding/json"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) createEventHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var eventVar event.Event
	if err := decoder.Decode(&eventVar); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := s.app.CreateEvent(r.Context(), &eventVar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventVar)
}

func (s *Server) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID события из URL
	vars := mux.Vars(r)
	idStr := vars["eventId"] // Извлекаем параметр eventId из URL
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Декодируем тело запроса
	decoder := json.NewDecoder(r.Body)
	var eventVar event.Event
	if err := decoder.Decode(&eventVar); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Обновляем событие
	eventVar.ID = id // Устанавливаем ID события
	err = s.app.UpdateEvent(r.Context(), &eventVar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventVar)
}

func (s *Server) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID события из URL
	vars := mux.Vars(r)
	idStr := vars["eventId"] // Извлекаем параметр eventId из URL
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	// Удаляем событие
	err = s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Успешное удаление
}

func (s *Server) getAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := s.app.GetAllEvents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) getEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["eventId"] // Извлекаем параметр eventId из URL
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	eventVar, err := s.app.GetEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if eventVar == nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eventVar)
}
