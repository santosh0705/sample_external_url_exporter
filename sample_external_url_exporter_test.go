package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	wantMetric1 = "sample_external_url_up"
	wantMetric2 = "sample_external_url_response_ms"
	wantDesc1   = `Desc{fqName: "sample_external_url_up", help: "Could the url be reached", constLabels: {}, variableLabels: [url]}`
	wantDesc2   = `Desc{fqName: "sample_external_url_response_ms", help: "Request response time in ms", constLabels: {}, variableLabels: [url]}`
)

func TestMetrics(t *testing.T) {
	// Create server.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	server := httptest.NewServer(handler)

	// Create exporter and read metrics.
	exporter := NewExporter([]string{server.URL})
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		exporter.Collect(ch)
	}()

	// Loop through the registered metrics and compare with wanted.
	for i := 1; i <= 2; i++ {
		m := <-ch
		if m == nil {
			t.Error("expected metric but got nil")
		} else {
			desc := m.Desc().String()
			if !(desc == wantDesc1 || desc == wantDesc2) {
				t.Errorf("expected '%s' or '%s', got '%s'", wantDesc1, wantDesc2, desc)
			}
		}
	}
	// Fail if additional metrices are registered.
	extraMetrics := 0
	for <-ch != nil {
		extraMetrics++
	}
	if extraMetrics > 0 {
		t.Errorf("expected closed channel, got %d extra metrics", extraMetrics)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	// Create two web servers, one with good and one bad responses.
	handlerGood := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	handlerBad := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	})
	serverGood := httptest.NewServer(handlerGood)
	serverBad := httptest.NewServer(handlerBad)

	exporter := NewExporter([]string{serverGood.URL, serverBad.URL})
	prometheus.MustRegister(exporter)

	req := httptest.NewRequest("GET", "/metrics", nil)

	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	// Tests response status.
	gotStatus := rr.Code
	gotBody := rr.Body.String()
	wantGoodStatus := http.StatusOK

	if gotStatus != wantGoodStatus {
		t.Errorf("Returned wrong status code: got %v want %v", gotStatus, wantGoodStatus)
	}

	// Test response body.
	// Prometheus /metrics endpoint will contain our metrics.
	wantMetric := wantMetric1 + `{url="` + serverGood.URL + `"} 1`
	if !strings.Contains(gotBody, wantMetric) {
		t.Errorf("Response does not match, expected: '%s' in' %v'", wantMetric, gotBody)
	}
	wantMetric = wantMetric1 + `{url="` + serverBad.URL + `"} 0`
	if !strings.Contains(gotBody, wantMetric) {
		t.Errorf("Response does not match, expected: '%s' in' %v'", wantMetric, gotBody)
	}
	wantMetric = wantMetric2 + `{url="` + serverGood.URL + `"} `
	if !strings.Contains(gotBody, wantMetric) {
		t.Errorf("Response does not match, expected: '%s' in' %v'", wantMetric, gotBody)
	}
	wantMetric = wantMetric2 + `{url="` + serverBad.URL + `"} `
	if !strings.Contains(gotBody, wantMetric) {
		t.Errorf("Response does not match, expected: '%s' in' %v'", wantMetric, gotBody)
	}
}
