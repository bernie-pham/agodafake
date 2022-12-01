CREATE TABLE "Reservations" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "start_date" timestamptz NOT NULL DEFAULT (now()),
  "end_date" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "room_id" bigserial NOT NULL
);

CREATE TABLE "Rooms" (
  "id" bigserial PRIMARY KEY,
  "room_name" varchar NOT NULL,
  "price_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Prices" (
  "id" bigserial PRIMARY KEY,
  "price_value" float8 NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "Room_restriction" (
  "id" bigserial PRIMARY KEY,
  "start_date" timestamptz NOT NULL DEFAULT (now()),
  "end_date" timestamptz NOT NULL DEFAULT (now()),
  "room_id" bigserial NOT NULL,
  "reservation_id" bigserial NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at_at" timestamptz NOT NULL DEFAULT (now()),
  "restriction_id" bigserial NOT NULL
);

CREATE TABLE "Restrictions" (
  "id" bigserial PRIMARY KEY,
  "restriction_name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "Reservations" ADD FOREIGN KEY ("room_id") REFERENCES "Rooms" ("id");

ALTER TABLE "Rooms" ADD FOREIGN KEY ("price_id") REFERENCES "Prices" ("id");

ALTER TABLE "Room_restriction" ADD FOREIGN KEY ("room_id") REFERENCES "Rooms" ("id");

ALTER TABLE "Room_restriction" ADD FOREIGN KEY ("reservation_id") REFERENCES "Reservations" ("id");

ALTER TABLE "Room_restriction" ADD FOREIGN KEY ("restriction_id") REFERENCES "Restrictions" ("id");
