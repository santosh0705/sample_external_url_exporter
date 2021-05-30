FROM golang:alpine3.13 as builder

COPY . /code

RUN cd /code; \
    go build -ldflags="-s -w"

FROM alpine:3.13

COPY --from=builder /code/sample_external_url_exporter /bin/sample_external_url_exporter

ENTRYPOINT ["/bin/sample_external_url_exporter"]
EXPOSE     8080
