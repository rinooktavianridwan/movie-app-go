CREATE TABLE promos (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    discount_type VARCHAR(50) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value NUMERIC(12,2) NOT NULL CHECK (discount_value > 0),
    min_tickets INTEGER DEFAULT 1 CHECK (min_tickets >= 1),
    max_discount NUMERIC(12,2) CHECK (max_discount IS NULL OR max_discount > 0),
    usage_limit INTEGER CHECK (usage_limit IS NULL OR usage_limit > 0),
    usage_count INTEGER DEFAULT 0 CHECK (usage_count >= 0),
    is_active BOOLEAN DEFAULT true,
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL CHECK (valid_until > valid_from),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_promos_code ON promos(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_promos_active ON promos(is_active, valid_from, valid_until) WHERE deleted_at IS NULL;
CREATE INDEX idx_promos_dates ON promos(valid_from, valid_until);