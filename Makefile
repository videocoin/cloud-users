GOOS?=linux
GOARCH?=amd64

NAME=users
GCP_PROJECT=videocoin-network
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)

DBM_MSQLURI=root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local

version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
		go build \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/main.go

build-dev:
	env GO111MODULE=on GOOS=${GOOS} GOARCH=${GOARCH} \
		go build \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/main.go

deps:
	env GO111MODULE=on go mod vendor

docker-build:
	docker build -t gcr.io/${GCP_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	gcloud docker -- push gcr.io/${GCP_PROJECT}/${NAME}:${VERSION}

dbm-status:
	goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" status

dbm-up:
	goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" up

dbm-down:
	goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" down

release: build docker-build docker-push
