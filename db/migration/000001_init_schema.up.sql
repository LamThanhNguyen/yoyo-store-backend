CREATE TABLE "items" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "inventory_level" integer NOT NULL,
  "price" bigint NOT NULL,
  "description" varchar DEFAULT '',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transaction_statuses" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "amount" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "last_four" varchar NOT NULL,
  "bank_return_code" varchar NOT NULL,
  "transaction_status_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE transactions
  ADD CONSTRAINT fk_transactions_transaction_status_id
  FOREIGN KEY (transaction_status_id)
  REFERENCES transaction_statuses(id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;