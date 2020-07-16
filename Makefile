APP ?= kube-ctl-svc

# Default env is always dev.
ENV ?= dev
BUILD_TAG ?= local
DOCKER_TAG ?= latest
DOCKER  ?= $(shell which docker)
COMPOSE ?= $(shell which docker-compose)
GIT ?= $(COMPOSE) -f build-tools.yml run --rm git
GO ?= $(COMPOSE) -f build-tools.yml run --rm go
NEXUS_USER ?= gitlab
NEXUS_PASSWORD ?=
COMMITSHA ?= DEV
COMMITTIME ?= NONE

build: ## Build inside docker
	mkdir -p app
	$(COMPOSE) build app
	$(COMPOSE) run app

push:
	curl -v -u $(NEXUS_USER):$(NEXUS_PASSWORD) --upload-file app/kube-svc-ctl https://nexus.flotech.co/repository/tools/kube-svc-ctl/kube-svc-ctl

