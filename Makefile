SHELL := /bin/bash

# ===============================================
# Testing running system

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2040
# openssl rsa -pubout -in private.pem -out public.pem

admin:
	go run app/tooling/admin/main.go

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

build :
	go build -ldflags "-X main.build=local"

tidy:
	go mod tidy
	go mod vendor

# Building containers

VERSION := 1.0

all: sales-api

sales-api:
	docker build --no-cache \
		-f deployments/docker/sales-api.Dockerfile \
		-t sales-api-amd64:$(VERSION) \
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
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-apply:
	kustomize build deployments/k8s/kind/sales-pod | kubectl apply -f -

kind-load:
	cd deployments/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)


kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-sales:
	kubectl get pods -o wide --watch --namespace=sales-system

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 --namespace=sales-system | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment sales-pod --namespace=sales-system


kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply


kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=sales