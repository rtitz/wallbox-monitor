package prometheus

import (
	"log"
	"net/http"
	"wallbox-monitor/variables"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Prometheus metrics server starting on port %s", port)
	if variables.ForceIpv4 {
		log.Printf("Forcing IPv4 binding on 0.0.0.0:%s", port)
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
	} else {
		log.Printf("Binding on all available interfaces on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
