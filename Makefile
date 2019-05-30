VER=${shell git describe --tags}

.PHONY: docker-build
docker-build:
	docker build --build-arg GOPROXY=${GOPROXY} -t iwittkau/minio-analytics .

.PHONY: docker-build-exporter
docker-build-exporter:
	docker build -f Dockerfile.exporter --build-arg GOPROXY=${GOPROXY} -t iwittkau/minio-analytics-exporter .

.PHONY: docker-push
docker-push:
	docker tag iwittkau/minio-analytics iwittkau/minio-analytics:${VER}
	docker tag iwittkau/minio-analytics iwittkau/minio-analytics:latest
	docker push iwittkau/minio-analytics:${VER}
	docker push iwittkau/minio-analytics:latest

.PHONY: docker-run-minio-server
docker-run-minio-server:
	docker run --rm -i --tty minio/minio:latest server /data