package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Nodes struct {
	IpAddr   string
	HostName string
}

// get all kubernetes nodes IP address and hostname from kube api
func (c *KubernetesClient) GetNodes() []Nodes {
	var kubeNodes []Nodes
	nodes, err := c.loadClient().CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, v := range nodes.Items {
		kubeNodes = append(kubeNodes, Nodes{IpAddr: v.Status.Addresses[0].Address, HostName: v.Status.Addresses[1].Address})
	}
	return kubeNodes
}

// func ping all kubeNodes IP addresses from each pod
func (c *KubernetesClient) PingNodes() {
	podname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return
	}
	podOnNode, err := c.loadClient().CoreV1().Pods(os.Getenv("NAMESPACE_NAME")).Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error getting pod: %v\n", err)
	}

	nodes := c.GetNodes()
	for _, ip := range nodes {
		pinger, err := probing.NewPinger(ip.IpAddr)
		if err != nil {
			panic(err)
		}
		pinger.Count = 4
		pinger.Timeout = 10 * time.Second
		pinger.SetPrivileged(true)
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			panic(err)
		}
		stats := pinger.Statistics()
		if stats.PacketLoss == 0 {
			fmt.Printf("Node %s is reachable from Pod %s on Node %s\n", ip.HostName, podname, podOnNode.Spec.NodeName)
		} else {
			fmt.Printf("Node %s is unreachable from Pod %s on Node %s\n", ip.HostName, podname, podOnNode.Spec.NodeName)
		}

	}
}
