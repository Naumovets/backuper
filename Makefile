# Makefile for Backuper Project

# Variables
IMAGE_NAME?=daniilnaumovets/backuper

# Build docker image
build:
	@echo "Building docker image $(IMAGE_NAME):$(VERSION)..."
	docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .
	@echo "Docker build completed: $(IMAGE_NAME):$(VERSION)"

# Start containers
up:
	docker-compose up -d --build

# Stop containers
down:
	docker-compose down

# Remove logs and backups
clean:
	rm -rf ./logs/* ./backups/*

help:
	@echo "Available commands:"
	@echo "  make up            - Start the backuper and database containers"
	@echo "	 make build			- Build images"
	@echo "  make down          - Stop the containers"
	@echo "  make clean         - Clear logs and backups"
