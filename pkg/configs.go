package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	kubernetesClient *kubernetes.Clientset
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Read kubeconfig file and return kubernetes clientset
func (c *KubernetesClient) loadClient() *kubernetes.Clientset {
	if _, err := os.Stat("/app/" + os.Getenv("KUBECONFIG_LOCATION")); err == nil {
		config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG_LOCATION"))
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		c.kubernetesClient = clientset
		return c.kubernetesClient
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		c.kubernetesClient = clientset
		return c.kubernetesClient
	}

}
