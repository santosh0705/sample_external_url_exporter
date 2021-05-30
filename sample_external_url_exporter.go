package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace        = "sample"
	subsystem        = "external_url"
	metricsEndpoint  = "/metrics"
	listeningAddress = ":8080"
)

var (
	externalURIs []string = []string{"https://httpstat.us/503", "https://httpstat.us/200"}
)

type Exporter struct {
	URIs   []string
	mutex  sync.Mutex
	client *http.Client

	up           *prometheus.GaugeVec
	responseTime *prometheus.GaugeVec
}

func NewExporter(uris []string) *Exporter {
	return &Exporter{
		URIs: uris,
		up: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "up",
			Help:      "Could the url be reached",
		},
			[]string{"url"},
		),
		responseTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "response_ms",
			Help:      "Request response time in ms",
		},
			[]string{"url"},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.up.Describe(ch)
	e.responseTime.Describe(ch)
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	for _, uri := range e.URIs {
		startTime := time.Now()
		resp, err := http.Get(uri)
		responseTime := float64(time.Since(startTime).Milliseconds())
		if err != nil {
			log.Printf(err.Error())
			e.up.WithLabelValues(uri).Set(0)
		} else {
			if resp.StatusCode != 200 {
				e.up.WithLabelValues(uri).Set(0)
			} else {
				e.up.WithLabelValues(uri).Set(1)
			}
			resp.Body.Close()
		}
		e.responseTime.WithLabelValues(uri).Set(responseTime)
	}

	e.up.Collect(ch)
	e.responseTime.Collect(ch)

	return nil
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Printf("Error request urls: %s", err)
	}
	return
}

func main() {
	exporter := NewExporter(externalURIs)
	prometheus.MustRegister(exporter)

	log.Printf("Starting %s exporter", namespace)

	http.Handle(metricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(listeningAddress, nil))
}
