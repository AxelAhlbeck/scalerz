CREATE TABLE question (
    id SERIAL PRIMARY KEY,
    question character varying(255) NOT NULL UNIQUE,
    answer character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);