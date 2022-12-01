server:
	go run cmd/web/main.go cmd/web/routes.go cmd/web/middleware.go cmd/web/mail_worker.go

dockerdb_init:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

create_db:
	docker exec -it postgres12 createdb --username=root --owner=root booking_room

drop_db:
	docker exec -it postgres12 dropdb booking_room

# generate migration file
# soda g sql db_init -p db/migrations

migrateup:
	soda migrate up -d -p db/migrations

migratedown:
	soda migrate down 2 -d -p db/migrations 

resetdb:
	soda reset -p db/migrations

.PHONY: server dockerdb_init create_db migrateup drop_db migratedown