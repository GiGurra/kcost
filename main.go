package main

import (
	"encoding/json"
	"fmt"
	"github.com/gigurra/kcost/pkg/model"
	"log/slog"
	"os"
	"os/exec"
)

func main() {
	// Estimate the cost of a Kubernetes cluster
	// We start by just fetching all pods and nodes

	// Check that kubectl is installed
	_, err := exec.LookPath("kubectl")
	if err != nil {
		slog.Error("kubectl not found")
		os.Exit(1)
	}

	nodes, err := getNodes()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting nodes: %v\n", err))
		os.Exit(1)
	}

	for _, node := range nodes {
		slog.Info(fmt.Sprintf("Node: %s { spot=%v, region=%v, zone=%s }\n", node.Name(), node.IsSpotNode(), node.Region(), node.Zone()))
	}

	pods, err := getPods()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting pods: %v\n", err))
		os.Exit(1)
	}

	for _, pod := range pods {
		slog.Info(fmt.Sprintf("Pod: %s { node=%s, cpu=%s, memory=%s }\n", pod.Name(), pod.NodeName(), pod.CPURequest(), pod.MemoryRequest()))
	}
}

func getNodes() ([]model.Node, error) {
	// Fetch all nodes in the cluster
	// use the kubectl command to get the nodes
	cmd := exec.Command("kubectl", "get", "nodes", "-o", "json")
	// Run the command and get the output
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %s, %s", err, out)
	}

	// Parse the output as JSON
	var listing model.K8sListing[model.Node]
	err = json.Unmarshal(out, &listing)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON output: %s", err)
	}

	return listing.Items, nil
}

func getPods() ([]model.Pod, error) {
	// Fetch all pods in the cluster
	// use the kubectl command to get the pods
	cmd := exec.Command("kubectl", "get", "pods", "-o", "json")
	// Run the command and get the output
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %s, %s", err, out)
	}

	// Parse the output as JSON
	var listing model.K8sListing[model.Pod]
	err = json.Unmarshal(out, &listing)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON output: %s", err)
	}
	return listing.Items, nil
}
