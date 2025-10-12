# Makefile to create a Go service project structure
# Usage:
#   make SERVICE=payment-service
# Assumes this Makefile is placed in the parent directory containing "service/"

SERVICE ?= payment-service
SERVICE_DIR := services

.PHONY: all create

all: create

create:
	@echo "Creating Go service: $(SERVICE)"
	cd $(SERVICE_DIR) && \
	mkdir -p $(SERVICE)/cmd && \
	mkdir -p $(SERVICE)/internal/adapters/inbound && \
	mkdir -p $(SERVICE)/internal/adapters/outbound && \
	mkdir -p $(SERVICE)/internal/app && \
	mkdir -p $(SERVICE)/internal/domain && \
	mkdir -p $(SERVICE)/internal/ports/inbound && \
	mkdir -p $(SERVICE)/internal/ports/outbound && \
	touch $(SERVICE)/cmd/main.go
	@echo "Project structure for $(SERVICE) created inside $(SERVICE_DIR)/"
