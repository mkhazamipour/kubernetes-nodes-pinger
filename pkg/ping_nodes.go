package pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.MustRegister(pingSuccessCounter)
	prometheus.MustRegister(pingFailureCounter)
}

var (
	k8sclient KubernetesClient
)

type Nodes struct {
	IpAddr   string
	HostName string
}

var (
	pingSuccessCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ping_success_total",
			Help: "Total number of successful pings",
		},
		[]string{"Hostname", "NodeIP", "Pod", "PodOnNode"},
	)

	pingFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ping_failure_total",
			Help: "Total number of failed pings",
		},
		[]string{"Hostname", "NodeIP", "Pod", "PodOnNode"},
	)
)

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

type Hosts struct {
	PodName       string
	PodNodeName   string
	HostName      string
	NodeIpAddress string
}

func HostsToPing() []Hosts {
	var hostlist []Hosts
	var hostname string
	var nodeipaddress string
	client := k8sclient
	nodes := client.GetNodes()
	for _, ip := range nodes {
		nodeipaddress = ip.IpAddr
		hostname = ip.HostName
		podname, err := os.Hostname()
		if err != nil {
			fmt.Printf("Error getting hostname: %v\n", err)
		}
		podnodename, err := client.loadClient().CoreV1().Pods(os.Getenv("NAMESPACE_NAME")).Get(context.TODO(), podname, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting pod: %v\n", err)
		}
		hostlist = append(hostlist, Hosts{PodName: podname, PodNodeName: podnodename.Spec.NodeName, NodeIpAddress: nodeipaddress, HostName: hostname})
	}
	return hostlist
}

func Pinger() {
	hosts := HostsToPing()
	for _, host := range hosts {
		pinger, err := probing.NewPinger(host.NodeIpAddress)
		if err != nil {
			panic(err)
		}
		pinger.Count = 4
		pinger.Timeout = 10 * time.Second
		pinger.Interval = 1 * time.Minute
		pinger.SetPrivileged(true)
		pinger.Run()
		stats := pinger.Statistics()
		if stats.PacketLoss == 0 {
			fmt.Printf("Node %s is reachable from Pod %s on Node %s\n", host.HostName, host.PodName, host.PodNodeName)
			pingSuccessCounter.WithLabelValues(host.HostName, host.NodeIpAddress, host.PodName, host.PodNodeName).Inc()
		} else {
			fmt.Printf("Node %s is unreachable from Pod %s on Node %s\n", host.HostName, host.PodName, host.PodNodeName)
			pingFailureCounter.WithLabelValues(host.HostName, host.NodeIpAddress, host.PodName, host.PodNodeName).Inc()

		}

	}

}
