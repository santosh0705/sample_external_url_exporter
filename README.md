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

Coming soon...

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
