CREATE TABLE public.Users
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    age INT,
    gender VARCHAR(100) NOT NULL,
    country_id VARCHAR(4) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP
);
CREATE INDEX id_user_create_at_pagination ON public.Users (created_at, id);
CREATE INDEX id_user_age_create_at_pagination ON public.Users (age, id);