CREATE TABLE IF NOT EXISTS "accounts" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar UNIQUE NOT NULL,
  "password" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS "accounts_numbers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "numbers" varchar UNIQUE NOT NULL,
  "balance" integer,
  "account_id" UUID,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS "transfers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "from" UUID NOT NULL,
  "to" UUID NOT NULL,
  "amount" integer,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS "histories" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "balance" integer,
  "amount" integer,
  "account_id" UUID,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

ALTER TABLE IF EXISTS "accounts_numbers" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE IF EXISTS "transfers" ADD FOREIGN KEY ("from") REFERENCES "accounts" ("id");

ALTER TABLE IF EXISTS "transfers" ADD FOREIGN KEY ("to") REFERENCES "accounts" ("id");

ALTER TABLE IF EXISTS "histories" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
