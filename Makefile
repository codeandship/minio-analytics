.PHONY: docker-build
docker-build:
	docker build --build-arg GOPROXY=${GOPROXY} -t iwittkau/minio-analytics .

.PHONY: docker-push
docker-push:
	docker push iwittkau/minio-analytics

.PHONY: docker-run-minio-server
docker-run-minio-server:
	docker run --rm -i --tty minio/minio:latest server /data