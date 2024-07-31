package main

import (
	"github.com/genson1808/balancer/config"
	"log"
	"net/http"
	"strconv"

	"github.com/genson1808/balancer/metrics"
	"github.com/genson1808/balancer/proxy"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = cfg.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	metrics.RegisterPrometheusMetrics()
	metrics.InitPrometheus()

	router := mux.NewRouter()
	router.Use(metrics.RecordRequestLatency)
	router.Use(metrics.RecordRequestCount)

	router.Handle("/metrics", promhttp.Handler())
	for _, l := range cfg.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, l.BalanceMode)
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		// start health check
		if cfg.HealthCheck {
			httpProxy.HealthCheck(cfg.HealthCheckInterval)
		}
		router.Handle(l.Pattern, httpProxy)
	}
	if cfg.MaxAllowed > 0 {
		router.Use(maxAllowedMiddleware(cfg.MaxAllowed))
	}
	svr := http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}

	// print config detail
	cfg.Print()

	// listen and serve
	if cfg.Schema == "http" {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if cfg.Schema == "https" {
		err := svr.ListenAndServeTLS(cfg.SSLCertificate, cfg.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}
