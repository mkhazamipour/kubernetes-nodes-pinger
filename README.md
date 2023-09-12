# Kubernetes nodes Pinger with Prometheus Metrics
This package provides functionality for pinging Kubernetes nodes from within pods in a Kubernetes cluster  to check the reachability of nodes and collects Prometheus metrics about the success and failure of these pings.
## Prerequisites

Before running this application, ensure you have the following prerequisites:

- Go installed on your machine.
- Access to a Kubernetes cluster (as this application interacts with Kubernetes nodes and pods).

## Installation
You can deploy it with deploy/deploy.yaml
Before running `kubectl create -f deploy/deploy.yaml` make sure you have replaced your kubeconfig in configmap and set or remove ingress accordingly
## Prometheus Metrics
This package provides two Prometheus counters to track the success and failure of ping operations:

- ping_success_total: Total number of successful pings.

- ping_failure_total: Total number of failed pings.

These metrics are labeled with the following information:

- Hostname: The hostname of the target node.
- NodeIP: The IP address of the target node.
- Pod: The name of the pod from which the ping was performed.
- PodOnNode: The name of the node where the pod is running.


You can scrape these metrics using a Prometheus server and visualize them using Grafana or any other monitoring tool of your choice.

## Disclaimer
This is part of my Golang learning