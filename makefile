PG_PASS ?= mypassword
PG_DB ?= wivvus

run: start-postgres
run:
	go run cmd/api/main.go	

start-postgres:
ifeq ($(podman ps | grep postgres | wc -l), 0)
	podman run --replace -d --name postgres -e POSTGRES_PASSWORD=${PG_PASS} -e POSTGRES_DB=${PG_DB} -p 5432:5432 -v pgdata:/var/lib/postgresql postgres
else
	@echo "postgres is already running"
endif