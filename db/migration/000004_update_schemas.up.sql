ALTER TABLE items
    ADD COLUMN image VARCHAR DEFAULT '',
    ADD COLUMN is_recurring BOOLEAN DEFAULT false,
    ADD COLUMN plan_id VARCHAR DEFAULT '';

ALTER TABLE transactions
    ADD COLUMN expiry_month INTEGER DEFAULT 0,
    ADD COLUMN expiry_year INTEGER DEFAULT 0,
    ADD COLUMN payment_intent VARCHAR DEFAULT '',
    ADD COLUMN payment_method VARCHAR DEFAULT '';

ALTER TABLE orders
    ADD COLUMN customer_id BIGINT;

ALTER TABLE orders
    ADD CONSTRAINT fk_orders_customer_id
    FOREIGN KEY (customer_id)
    REFERENCES customers(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE;

INSERT INTO "items" ("name", "inventory_level", "price", "description", "image", "is_recurring", "plan_id")
VALUES
  ('Yoyo', 10, 1000, 'A very nice yoyo', '/static/yoyo4.png', false, ''),
  ('Bronze Plan', 10, 2000, 'Get three yoyo for the price of two every month.', '', true, 'price_1RbGnAQ3QRfxqnBZjUtkxUsT');