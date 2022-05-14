IMG ?= kubeberth/berth-apiserver:v1alpha1

.PHONY: build
build:
	go build -o bin/berth-apiserver main.go

.PHONY: run
run:
	go run main.go

.PHONY: deploy
deploy:
	kubectl apply -f manifest.yaml

.PHONY: undeploy
undeploy:
	kubectl delete -f manifest.yaml

.PHONY: docker-build
docker-build:
	docker build . -t ${IMG}

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: docker-buildx
docker-buildx:
	docker buildx build --platform linux/amd64,linux/arm64 -t ${IMG} --push .
