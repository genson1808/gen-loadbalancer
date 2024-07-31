package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var collectorContainer []prometheus.Collector

//InitPrometheus ... initalize prometheus
func InitPrometheus() {
	prometheus.MustRegister(collectorContainer...)
}

//PushRegister ... Push collectores to prometheus before inializing
func PushRegister(c ...prometheus.Collector) {
	collectorContainer = append(collectorContainer, c...)
}

func RegisterPrometheusMetrics() {
	PushRegister(totalRequest)
	PushRegister(latency)
	PushRegister(statusTraceFailure)
	PushRegister(statusTraceSuccess)
}

func RecordRequestLatency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start).Milliseconds()
		latency.WithLabelValues(
			r.Method,
			r.URL.Path,
		).Observe(float64(elapsed))
	})
}

func RecordRequestCount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		totalRequest.Inc()
	})
}

func IncSuccess(ip string) {
	statusTraceSuccess.WithLabelValues(ip).Inc()
}

func IncFailure(ip string) {
	statusTraceFailure.WithLabelValues(ip).Inc()
}

var latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "Balancer",
		Name:       "latency_duration_seconds",
		Help:       "Load balancer latency distribution",
		Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"method", "path"},
)

var statusTraceSuccess = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "Balancer",
		Name:      "status_success_count",
		Help:      "Number of request success",
	},
	[]string{"IP"},
)

var statusTraceFailure = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "Balancer",
		Name:      "status_failure_count",
		Help:      "Number of request faild",
	},
	[]string{"IP"},
)

var totalRequest = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "Balancer",
		Name:      "request_count",
		Help:      "Number of request",
	},
)
