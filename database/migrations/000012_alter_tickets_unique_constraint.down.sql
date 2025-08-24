DROP INDEX IF EXISTS idx_tickets_schedule_seat_active;

CREATE UNIQUE INDEX idx_tickets_schedule_seat 
ON tickets(schedule_id, seat_number) 
WHERE deleted_at IS NULL;