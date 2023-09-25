package kubectl

import (
	"encoding/json"
	"fmt"
	"github.com/gigurra/kcost/pkg/model"
	"os/exec"
	"strings"
)

func getListing[T any](kind string, namespace string) ([]T, error) {
	// Fetch all nodes in the cluster
	// use the kubectl command to get the nodes
	args := []string{"get", kind, "-o", "json"}
	if namespace != "" {
		args = append(args, "-n", namespace)
	}
	cmd := exec.Command("kubectl", args...)
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
	return getListing[model.Node]("nodes", "")
}

func GetPods(namespace string) ([]model.Pod, error) {
	return getListing[model.Pod]("pods", namespace)
}

//
//func GetNamespace() (string, error) {
//	cmd := exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={..namespace}")
//	bytes, err := cmd.Output()
//	if err != nil {
//		return "", fmt.Errorf("error running kubectl command: %s, %s", err, bytes)
//	}
//	return string(bytes), nil
//}

func GetAllNamespaces() ([]string, error) {
	cmd := exec.Command("kubectl", "get", "namespaces", "-o", "jsonpath={..metadata.name}")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running kubectl command: %s, %s", err, bytes)
	}
	// Split the output into a slice of strings
	return strings.Split(string(bytes), " "), nil
}
