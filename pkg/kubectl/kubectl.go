package kubectl

import (
	"encoding/json"
	"fmt"
	"github.com/gigurra/kcost/pkg/model"
	"os/exec"
)

func getListing[T any](kind string) ([]T, error) {
	// Fetch all nodes in the cluster
	// use the kubectl command to get the nodes
	cmd := exec.Command("kubectl", "get", kind, "-o", "json")
	// Run the command and get the output
	bytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %s, %s", err, bytes)
	}

	// Parse the output as JSON
	var listing model.K8sListing[T]
	err = json.Unmarshal(bytes, &listing)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON output: %s", err)
	}

	return listing.Items, nil
}

func GetNodes() ([]model.Node, error) {
	return getListing[model.Node]("nodes")
}

func GetPods() ([]model.Pod, error) {
	return getListing[model.Pod]("pods")
}
