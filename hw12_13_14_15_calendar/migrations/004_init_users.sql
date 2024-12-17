-- +goose Up
INSERT INTO public.users (user_name) SELECT 'admin'
WHERE NOT EXISTS (
    SELECT 1 FROM public.users
);