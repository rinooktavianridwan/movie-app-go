CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('movie_reminder', 'promo_available', 'booking_confirmation', 'system')),
    is_read BOOLEAN DEFAULT false,
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_is_read ON notifications(is_read) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_type ON notifications(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_created_at ON notifications(created_at) WHERE deleted_at IS NULL;
