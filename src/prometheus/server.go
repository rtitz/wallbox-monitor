package prometheus

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Prometheus metrics server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}