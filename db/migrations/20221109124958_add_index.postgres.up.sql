CREATE UNIQUE INDEX "Users_email_idx" ON "Users"(email);

CREATE UNIQUE INDEX "Rr_start_end_date_idx" ON "Room_restriction" ("start_date", "end_date");

CREATE INDEX "Rr_roomID_idx" ON "Room_restriction" ("room_id");

CREATE UNIQUE INDEX "Rr_reservationID_idx" ON "Room_restriction" ("reservation_id");

CREATE INDEX "Rs_email_idx" ON "Reservations" ("email");

CREATE INDEX "Rs_last_name_idx" ON "Reservations" ("last_name");



