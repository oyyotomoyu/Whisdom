SHELL := /bin/bash

.PHONY: install dev build test lint clean server client docs

install:
	$(MAKE) -C server install
	$(MAKE) -C client install

# Runs the Go backend and the React dev server together; Ctrl+C stops both.
dev:
	@trap 'kill 0' EXIT INT TERM; \
	$(MAKE) -C server run & \
	$(MAKE) -C client dev & \
	wait

build:
	$(MAKE) -C client build
	$(MAKE) -C server build

test:
	$(MAKE) -C server test
	$(MAKE) -C client test

lint:
	$(MAKE) -C server lint
	$(MAKE) -C client lint

clean:
	$(MAKE) -C server clean
	$(MAKE) -C client clean

server:
	$(MAKE) -C server run

client:
	$(MAKE) -C client dev

docs:
	@echo "No documentation tooling configured yet."
