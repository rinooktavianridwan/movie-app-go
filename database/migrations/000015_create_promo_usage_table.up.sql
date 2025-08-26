CREATE TABLE promo_usages (
    id SERIAL PRIMARY KEY,
    promo_id INTEGER NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    discount_amount NUMERIC(12,2) NOT NULL CHECK (discount_amount >= 0),
    used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_promo_usage_promo_id ON promo_usages(promo_id);
CREATE INDEX idx_promo_usage_user_id ON promo_usages(user_id);
CREATE INDEX idx_promo_usage_transaction_id ON promo_usages(transaction_id);
CREATE INDEX idx_promo_usage_used_at ON promo_usages(used_at);

CREATE UNIQUE INDEX idx_promo_usage_transaction ON promo_usages(transaction_id);