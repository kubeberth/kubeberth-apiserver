IMG ?= kubeberth/kubeberth-apiserver:v1alpha1

.PHONY: build
build:
	go build -o bin/kubeberth-apiserver main.go

.PHONY: run
run:
	go run main.go

.PHONY: deploy
deploy:
	kubectl apply -f manifest.yaml

.PHONY: undeploy
undeploy:
	kubectl delete -f manifest.yaml

.PHONY: redeploy
redeploy: undeploy deploy

.PHONY: docker-build
docker-build:
	docker build --no-cache -t ${IMG} .

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: docker-buildx
docker-buildx:
	docker buildx build --no-cache --platform linux/amd64,linux/arm64 -t ${IMG} --push .
