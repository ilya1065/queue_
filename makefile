.PHONY: up down force1 force0

up:
	migrate -path db/migrations -database "sqlite3://./data/app.db" up

down:
	migrate -path db/migrations -database "sqlite3://./data/app.db" down

force1:
	migrate -path db/migrations -database "sqlite3://./data/app.db" force 1

force0:
	migrate -path db/migrations -database "sqlite3://./data/app.db" force 0