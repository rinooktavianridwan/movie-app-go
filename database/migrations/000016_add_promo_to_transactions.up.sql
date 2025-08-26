ALTER TABLE transactions 
ADD COLUMN promo_id INTEGER REFERENCES promos(id) ON DELETE SET NULL,
ADD COLUMN discount_amount NUMERIC(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
ADD COLUMN original_amount NUMERIC(12,2);

CREATE INDEX idx_transactions_promo_id ON transactions(promo_id);

COMMENT ON COLUMN transactions.promo_id IS 'Reference to applied promo (nullable)';
COMMENT ON COLUMN transactions.discount_amount IS 'Amount discounted from original total';
COMMENT ON COLUMN transactions.original_amount IS 'Original amount before discount';