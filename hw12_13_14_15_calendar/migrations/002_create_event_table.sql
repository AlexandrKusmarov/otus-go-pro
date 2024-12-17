-- +goose Up
CREATE TABLE IF NOT EXISTS public.events
(
    id                  bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title               character varying,
    event_date_time     time with time zone,
    event_end_date_time time with time zone,
    description         text,
    user_id             bigint REFERENCES public.users (id),
    notify_before_event time with time zone
);