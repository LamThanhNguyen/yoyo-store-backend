CREATE TABLE "statuses" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);


CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "item_id" bigint NOT NULL,
  "transaction_id" bigint NOT NULL,
  "status_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE orders
  ADD CONSTRAINT fk_orders_item_id
  FOREIGN KEY (item_id)
  REFERENCES items(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE orders
  ADD CONSTRAINT fk_orders_transaction_id
  FOREIGN KEY (transaction_id)
  REFERENCES transactions(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE orders
  ADD CONSTRAINT fk_orders_status_id
  FOREIGN KEY (status_id)
  REFERENCES statuses(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;