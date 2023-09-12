package main

import (
	"net/http"

	"github.com/mkhazamipour/kubernetes-network-functionality-check/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	pkg.LoadEnvVariables()

}
func main() {

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	for {
		pkg.Pinger()
	}

}
