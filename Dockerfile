FROM golang:1.12-alpine

ARG GOPROXY

WORKDIR /build

ADD . /build

RUN GOPROXY=$GOPROXY CGO_ENABLED=0 go build -o minio-analytics ./cmd/minio-analytics

FROM scratch

VOLUME [ "/data" ]

COPY --from=0 /build/minio-analytics /minio-analytics

EXPOSE 80

ENTRYPOINT [ "/minio-analytics" ]
