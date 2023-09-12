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
	http.ListenAndServe(":2112", nil)
	kubeClient := pkg.KubernetesClient{}
	kubeClient.PingNodes()
}
