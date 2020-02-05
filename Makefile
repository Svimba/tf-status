SHELL:=/bin/bash
PWD := $(shell pwd)
TAG :="willco/tf-status"


.PHONY: build
build: ## Build tungsten-fabric-operator executable file in local go env
	echo "Building tf-status-proxy bin"
	go build -o build/_output/bin/tf-status-proxy src/main/main.go
	echo "Building container"
	docker build -t $(TAG) -f build/Dockerfile .

push:
	docker push $(TAG)


