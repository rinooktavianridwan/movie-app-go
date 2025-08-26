DROP INDEX IF EXISTS idx_transactions_promo_id;

ALTER TABLE transactions 
DROP COLUMN IF EXISTS original_amount,
DROP COLUMN IF EXISTS discount_amount,
DROP COLUMN IF EXISTS promo_id;