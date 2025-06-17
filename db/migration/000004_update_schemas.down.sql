ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_customer_id;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_id;

ALTER TABLE items DROP COLUMN IF EXISTS image;
ALTER TABLE items DROP COLUMN IF EXISTS is_recurring;
ALTER TABLE items DROP COLUMN IF EXISTS plan_id;

ALTER TABLE transactions DROP COLUMN IF EXISTS expiry_month;
ALTER TABLE transactions DROP COLUMN IF EXISTS expiry_year;
ALTER TABLE transactions DROP COLUMN IF EXISTS payment_intent;
ALTER TABLE transactions DROP COLUMN IF EXISTS payment_method;