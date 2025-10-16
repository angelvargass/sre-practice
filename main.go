package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sre_practice_http_duration",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	requestNumber = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sre_practice_requests_number",
		Help: "The total number of processed requests",
	})
)

func main() {
	prometheus.MustRegister(requestLatency)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", handler)

	http.ListenAndServe(":8080", nil) //nolint
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer requestNumber.Inc()
	start := time.Now()

	duration := time.Since(start).Seconds()
	requestLatency.WithLabelValues(r.URL.Path, r.Method).Observe(duration)
	w.Write([]byte("Hello, Prometheus!")) //nolint
}
