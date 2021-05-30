# Sample External URL Exporter for Prometheus

Exports status of external web service URLs and access response time via HTTP for Prometheus consumption.

# Building and running

## Build binary

With working golang environment it can be built with `go build` or `go get` on the project directory. With `go build` the binary will be saved in the project directory itself, with `go get` it will be saved in the `$GOPATH/bin`.
```
go build
```

## Run

The binary doesn't accept any parameter. The service bind to the port `8080` and the metrics are exposed at `/metrics` endpoint.
```
./sample_external_url_exporter
``` 

# Using Docker

## Build image

With working docker environment you can built the docker image with `docker build .`. Docker image building is multi stage image builder, it uses intermediate container for compiling binary. You don't need a working golang environment to build docker image.
```
docker build -t sample-external-url-exporter .
```

## Tag and push image

Tag the image with your private or public registry repository and push it to the remote docker registry so it can be used in other system or k8s environment. You may have to login to the remote registry first to push the image.
```
docker image tag sample-external-url-exporter santosh0705/sample-external-url-exporter
docker image push santosh0705/sample-external-url-exporter
```

## Run

```
docker run -d -p 8080:8080 santosh0705/sample-external-url-exporter
```

# Using Kubernetes

Coming soon...

# Monitoring using Prometheus and Grafana

Coming soon...

## Collectors

Sample external URL metrics:

```
# HELP sample_external_url_response_ms Request response time in ms
# TYPE sample_external_url_response_ms gauge
# HELP sample_external_url_up Could the url be reached
# TYPE sample_external_url_up gauge
```

## Author

The exporter is created by [santosh0705](https://github.com/santosh0705).
