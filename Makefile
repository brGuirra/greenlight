#!make
include .env

migrate_up:
	migrate -path ./migrations -database ${DATABASE_URL} -verbose up

migrate_down:
	migrate -path ./migrations -database ${DATABASE_URL} -verbose down

server:
	go run ./cmd/api -db-dsn=${DATABASE_URL} -smtp-host=${SMPT_HOST} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD} -smtp-sender=${SMTP_SENDER} -cors-trusted-origins=${CORS_TRUSTED_ORIGINS}


.PHONY: migrate_up migrate_down server
