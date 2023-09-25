package main

import (
	"fmt"
	"github.com/gigurra/kcost/pkg/kubectl"
	"log/slog"
	"os"
	"os/exec"
)

func main() {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		slog.Error("kubectl not found")
		os.Exit(1)
	}

	nodes, err := kubectl.GetNodes()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting nodes: %v\n", err))
		os.Exit(1)
	}

	for _, node := range nodes {
		slog.Info(fmt.Sprintf("Node: %s { spot=%v, region=%v, zone=%s }\n", node.Name(), node.IsSpotNode(), node.Region(), node.Zone()))
	}

	pods, err := kubectl.GetPods()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting pods: %v\n", err))
		os.Exit(1)
	}

	for _, pod := range pods {
		slog.Info(fmt.Sprintf("Pod: %s { node=%s, cpu=%s, memory=%s }\n", pod.Name(), pod.NodeName(), pod.CPURequest(), pod.MemoryRequest()))
	}
}
