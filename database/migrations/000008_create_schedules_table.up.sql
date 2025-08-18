CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    studio_id INTEGER NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    date DATE NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk performa query
CREATE INDEX idx_schedules_movie_id ON schedules(movie_id);
CREATE INDEX idx_schedules_studio_id ON schedules(studio_id);
CREATE INDEX idx_schedules_date ON schedules(date);
CREATE INDEX idx_schedules_start_time ON schedules(start_time);

-- Constraint untuk memastikan tidak ada jadwal bentrok di studio yang sama
CREATE UNIQUE INDEX idx_schedules_studio_time ON schedules(studio_id, start_time, end_time) WHERE deleted_at IS NULL;