REGISTRY ?= narahari92
IMG ?= vaali
TAG ?= 1.0

build-amd:
	docker build --tag ${REGISTRY}/${IMG}-amd:${TAG} --build-arg amd64 .

build-arm:
	docker build --tag ${REGISTRY}/${IMG}-arm:${TAG} --build-arg arm64 .

push-amd:
	docker push ${REGISTRY}/${IMG}-amd:${TAG}

push-arm:
	docker push ${REGISTRY}/${IMG}-arm:${TAG}
