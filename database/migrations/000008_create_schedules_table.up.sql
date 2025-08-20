CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    studio_id INTEGER NOT NULL REFERENCES studios(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    date DATE NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);