CREATE TABLE IF NOT EXISTS bank."accounts" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar UNIQUE NOT NULL,
  "password" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS bank."account_numbers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "number" varchar UNIQUE NOT NULL,
  "balance" integer DEFAULT 0,
  "account_id" UUID,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS bank."transfers" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "from" UUID NOT NULL,
  "to" UUID NOT NULL,
  "amount" integer,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

CREATE TABLE IF NOT EXISTS bank."histories" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "balance" integer,
  "amount" integer,
  "account_number_id" UUID,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_deleted" bool DEFAULT false
);

ALTER TABLE IF EXISTS bank."account_numbers" ADD FOREIGN KEY ("account_id") REFERENCES bank."accounts" ("id");

ALTER TABLE IF EXISTS bank."transfers" ADD FOREIGN KEY ("from") REFERENCES bank."account_numbers" ("id");

ALTER TABLE IF EXISTS bank."transfers" ADD FOREIGN KEY ("to") REFERENCES bank."account_numbers" ("id");

ALTER TABLE IF EXISTS bank."histories" ADD FOREIGN KEY ("account_number_id") REFERENCES bank."account_numbers" ("id");
