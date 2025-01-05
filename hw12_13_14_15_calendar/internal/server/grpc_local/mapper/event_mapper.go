package mapper

import (
	eventpb "github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/internal/server/grpc_local/pb"
	"github.com/AlexandrKusmarov/otus-go-pro/hw12_13_14_15_calendar/model/event"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func EventGrpcToEventStorage(grpcEvent *eventpb.Event) *event.Event {
	return &event.Event{
		ID:                grpcEvent.Id,
		Title:             grpcEvent.Title,
		Description:       grpcEvent.Description,
		UserID:            grpcEvent.UserId,
		EventDateTime:     grpcEvent.EventDateTime.AsTime(),
		EventEndDateTime:  grpcEvent.EventEndDateTime.AsTime(),
		NotifyBeforeEvent: grpcEvent.NotifyBeforeEvent.AsTime(),
	}
}

func EventStorageToEventGrpc(eventStor *event.Event) *eventpb.Event {
	return &eventpb.Event{
		Id:                eventStor.ID,
		Title:             eventStor.Title,
		EventDateTime:     &timestamp.Timestamp{Seconds: eventStor.EventDateTime.Unix()},
		EventEndDateTime:  &timestamp.Timestamp{Seconds: eventStor.EventEndDateTime.Unix()},
		Description:       eventStor.Description,
		UserId:            eventStor.UserID,
		NotifyBeforeEvent: &timestamp.Timestamp{Seconds: eventStor.NotifyBeforeEvent.Unix()},
	}
}
