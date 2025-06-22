CREATE TABLE "tokens" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "token_hash" bytea NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "expiry" timestamptz NOT NULL
);

CREATE TABLE "sessions" (
  "token" varchar PRIMARY KEY,
  "data" bytea NOT NULL,
  "expiry" timestamptz NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);