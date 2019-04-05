GOOS?=linux
GOARCH?=amd64

NAME=vc-user
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)

CHARTS=loggly
CHARTS_BUCKET=videocoin-cloud-charts
GCP_PROJECT=videocoin

DBM_MSQLURI=host=127.0.0.1 user=mysql dbname=videocoin sslmode=disable password=mysql

version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
		go build \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/main.go

docker-build:
	docker build -t gcr.io/${GCP_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	gcloud docker -- push gcr.io/${GCP_PROJECT}/${NAME}:${VERSION}

helm-package:
	@echo "Packaging ${NAME}..."
	@helm package --save=false -d helm/ helm/charts/${NAME}

helm-repo-index:
	@echo "Indexing charts repository..."
	@helm repo index helm/repo --url https://${CHARTS_BUCKET}.storage.googleapis.com

helm-repo-sync:
	@echo "Syncing repo..."
	@gsutil -m -h "Cache-Control:public,max-age=0" cp -a public-read helm/repo/* gs://${CHARTS_BUCKET}

helm-repo-update: helm-package helm-repo-index helm-repo-sync

dbm-status:
	goose -dir migrations -table ${NAME} postgres "${DBM_PGURI}" status

dbm-up:
	goose -dir migrations -table ${NAME} postgres "${DBM_PGURI}" up

dbm-down:
	goose -dir migrations -table ${NAME} postgres "${DBM_PGURI}" down

release: build docker-build docker-push
