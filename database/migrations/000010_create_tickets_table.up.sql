CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id),
    schedule_id INTEGER NOT NULL REFERENCES schedules(id),
    seat_number INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);