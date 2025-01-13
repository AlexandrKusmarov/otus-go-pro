package grpc_local

import (
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/grpc_local/mapper"
	eventpb "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/grpc_local/pb"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/storage/common"
	"golang.org/x/net/context"
)

type GrpcServer struct {
	eventpb.UnimplementedEventServiceServer
	storage common.StorageInterface
}

func NewGrpcServer(storage *common.StorageInterface) *GrpcServer {
	return &GrpcServer{
		storage: *storage,
	}
}

func (s *GrpcServer) CreateEvent(ctx context.Context, req *eventpb.CreateEventRequest) (*eventpb.CreateEventResponse, error) {
	eventForDb := mapper.EventGrpcToEventStorage(req.Event)

	err := s.storage.CreateEvent(ctx, eventForDb)
	if err != nil {
		return nil, err
	}
	return &eventpb.CreateEventResponse{Event: req.Event}, nil
}

func (s *GrpcServer) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	eventFromDb, err := s.storage.GetEvent(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	eventResp := &eventpb.GetEventResponse{
		Event: mapper.EventStorageToEventGrpc(eventFromDb),
	}

	return eventResp, nil
}

func (s *GrpcServer) UpdateEvent(ctx context.Context, req *eventpb.UpdateEventRequest) (*eventpb.UpdateEventResponse, error) {
	eventForDb := mapper.EventGrpcToEventStorage(req.Event)

	err := s.storage.UpdateEvent(ctx, eventForDb)
	if err != nil {
		return nil, err
	}
	return &eventpb.UpdateEventResponse{Event: req.Event}, nil
}

func (s *GrpcServer) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	err := s.storage.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &eventpb.DeleteEventResponse{}, nil
}

func (s *GrpcServer) GetAllEvents(ctx context.Context, req *eventpb.GetAllEventsRequest) (*eventpb.GetAllEventsResponse, error) {
	// Получаем все события из хранилища
	eventsFromDb, err := s.storage.GetAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	// Преобразуем события из формата хранилища в формат gRPC
	var events []*eventpb.Event
	for _, event := range eventsFromDb {
		events = append(events, mapper.EventStorageToEventGrpc(&event))
	}

	// Формируем ответ
	eventResp := &eventpb.GetAllEventsResponse{
		Events: events,
	}

	return eventResp, nil
}

func (s *GrpcServer) GetAllEventsForDay(ctx context.Context, req *eventpb.GetAllEventsForDayRequest) (*eventpb.GetAllEventsForDayResponse, error) {
	// Получаем все события за указанный день из хранилища
	eventsFromDb, err := s.storage.GetAllEventsForDay(ctx, req.Day.AsTime())
	if err != nil {
		return nil, err
	}

	// Преобразуем события из формата хранилища в формат gRPC
	var events []*eventpb.Event
	for _, event := range eventsFromDb {
		events = append(events, mapper.EventStorageToEventGrpc(&event))
	}

	// Формируем ответ
	eventResp := &eventpb.GetAllEventsForDayResponse{
		Events: events,
	}

	return eventResp, nil
}

func (s *GrpcServer) GetAllEventsForWeek(ctx context.Context, req *eventpb.GetAllEventsForWeekRequest) (*eventpb.GetAllEventsForWeekResponse, error) {
	// Получаем все события за указанную неделю из хранилища
	eventsFromDb, err := s.storage.GetAllEventsForWeek(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}

	// Преобразуем события из формата хранилища в формат gRPC
	var events []*eventpb.Event
	for _, event := range eventsFromDb {
		events = append(events, mapper.EventStorageToEventGrpc(&event))
	}

	// Формируем ответ
	eventResp := &eventpb.GetAllEventsForWeekResponse{
		Events: events,
	}

	return eventResp, nil
}

func (s *GrpcServer) GetAllEventsForMonth(ctx context.Context, req *eventpb.GetAllEventsForMonthRequest) (*eventpb.GetAllEventsForMonthResponse, error) {
	// Получаем все события за указанный месяц из хранилища
	eventsFromDb, err := s.storage.GetAllEventsForMonth(ctx, req.Start.AsTime())
	if err != nil {
		return nil, err
	}

	// Преобразуем события из формата хранилища в формат gRPC
	var events []*eventpb.Event
	for _, event := range eventsFromDb {
		events = append(events, mapper.EventStorageToEventGrpc(&event))
	}

	// Формируем ответ
	eventResp := &eventpb.GetAllEventsForMonthResponse{
		Events: events,
	}

	return eventResp, nil
}
