package main

import (
	"github.com/mkhazamipour/kubernetes-network-functionality-check/pkg"
)

func init() {
	pkg.LoadEnvVariables()
}
func main() {

	kubeClient := pkg.KubernetesClient{}
	kubeClient.PingNodes()
}
