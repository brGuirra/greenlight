#!make
include .env

migrate_up:
	migrate -path=./migrations -database ${DATABASE_URL} -verbose up

migrate_down:
	migrate -path=./migrations -database ${DATABASE_URL} -verbose down

server:
	go run ./cmd/api


.PHONY: migrate_up migrate_down server 
