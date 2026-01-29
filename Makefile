ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: server start lint vet fix with_logs

LOG_FILE?=$(shell pwd)/${LOG_DIR}/$(shell date +%Y-%m-%d).log

FILTER?="./..."

# Generic logging wrapper that can be used by any target. Redirecting stderr to stdout to make use of pipe and tee.
# look https://en.wikipedia.org/wiki/Tee_(command), T-split output to stdout and a file
with_logs:
	@if [ -z "$(TARGET)" ]; then \
		echo "Error: TARGET must be specified. Usage: make with_logs TARGET=your-command"; \
		exit 1; \
	fi
	@echo "Running $(TARGET) with logging to $(LOG_FILE)..."
	@mkdir -p $(dir $(LOG_FILE))
	@FORCE_COLOR=1 $(MAKE) $(TARGET) 2>&1 | tee -a $(LOG_FILE)

server:
	@air -c air.toml

# Start server with hot reload and logging
start:
	@$(MAKE) with_logs TARGET=server

lint:
	golangci-lint run --fix $(FILTER) && go fmt -x -n $(FILTER)

vet:
	go vet $(FILTER)

#Only used when upgrading go versions
fix:
	go fix -x $(FILTER)
