SHELL := /bin/bash

run:
	go run main.go

build :
	go build -ldflags "-X main.build=local"

# Building containers

VERSION := 1.0

all: service

service:
	docker build \
		-f deployments/docker/Dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.



# Running from within k8s/kind
KIND_CLUSTER := service-cluster
kind-up:
	kind create cluster \
		--image kindest/node:v1.27.3@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72 \
		--name $(KIND_CLUSTER) \
		--config deployments/k8s/kind/kind-config.yaml
	#kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
