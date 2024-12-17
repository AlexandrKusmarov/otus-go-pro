-- +goose Up
CREATE TABLE IF NOT EXISTS public.notification
(
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_id        bigint REFERENCES public.events(id),
    title           character varying,
    event_date_time time with time zone,
    user_id         bigint REFERENCES public.users (id)
);