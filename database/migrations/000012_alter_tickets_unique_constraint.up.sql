DROP INDEX IF EXISTS idx_tickets_schedule_seat;

CREATE UNIQUE INDEX idx_tickets_schedule_seat_active 
ON tickets(schedule_id, seat_number) 
WHERE deleted_at IS NULL AND status IN ('pending', 'active', 'used');