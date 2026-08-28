SHELL := /bin/bash

.PHONY: install dev build test lint clean server client docs

install:
	$(MAKE) -C server install
	$(MAKE) -C UI install

# Runs the Go backend and the React dev server together; Ctrl+C stops both.
dev:
	@trap 'kill 0' EXIT INT TERM; \
	$(MAKE) -C server run & \
	$(MAKE) -C UI dev & \
	wait

build:
	$(MAKE) -C UI build
	$(MAKE) -C server build

test:
	$(MAKE) -C server test
	$(MAKE) -C UI test

lint:
	$(MAKE) -C server lint
	$(MAKE) -C UI lint

clean:
	$(MAKE) -C server clean
	$(MAKE) -C UI clean

server:
	$(MAKE) -C server run

client:
	$(MAKE) -C UI dev

docs:
	@echo "No documentation tooling configured yet."
