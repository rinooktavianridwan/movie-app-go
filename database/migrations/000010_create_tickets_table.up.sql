CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    schedule_id INTEGER NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
    seat_number INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Index untuk performa query
CREATE INDEX idx_tickets_transaction_id ON tickets(transaction_id);
CREATE INDEX idx_tickets_schedule_id ON tickets(schedule_id);

-- Constraint untuk memastikan tidak ada duplikasi seat pada schedule yang sama
CREATE UNIQUE INDEX idx_tickets_schedule_seat ON tickets(schedule_id, seat_number) WHERE deleted_at IS NULL;