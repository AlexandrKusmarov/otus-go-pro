-- +goose Up
CREATE TABLE IF NOT EXISTS public.users (
                                            id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                            user_name character varying NOT NULL
);