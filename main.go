package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("Node: %s { spot=%v, region=%v, zone=%s }\n", node.getName(), node.isSpotNode(), node.region(), node.zone())
	}
}

func getNodes() ([]Node, error) {
	// Fetch all nodes in the cluster
	// use the kubectl command to get the nodes
	cmd := exec.Command("kubectl", "get", "nodes", "-o", "json")
	// Run the command and get the output
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %s, %s", err, out)
	}

	// Parse the output as JSON
	var listing K8sListing[Node]
	err = json.Unmarshal(out, &listing)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON output: %s", err)
	}

	return listing.Items, nil
}

type K8sListing[T any] struct {
	// Items is the list of nodes
	Items []T `json:"items"`
}

// We are only interested in the node name and labels
type Node struct {
	// Name and labels are found in the metadata section
	Metadata struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func (n Node) getName() string {
	return n.Metadata.Name
}

func (n Node) getLabels() map[string]string {
	return n.Metadata.Labels
}

func (n Node) isSpotNode() bool {
	value, ok := n.Metadata.Labels["cloud.google.com/gke-spot"]
	return ok && value == "true"
}

func (n Node) region() string {
	return n.Metadata.Labels["failure-domain.beta.kubernetes.io/region"]
}

func (n Node) zone() string {
	return n.Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
}

func (n Node) tpe() string {
	return n.Metadata.Labels["beta.kubernetes.io/instance-type"]
}
