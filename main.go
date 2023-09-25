package main

import (
	"fmt"
	"github.com/gigurra/kcost/pkg/kubectl"
	"github.com/gigurra/kcost/pkg/model"
	"log/slog"
	"os"
	"os/exec"
	"slices"
)

func main() {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		slog.Error("kubectl not found")
		os.Exit(1)
	}

	config, err := model.NewConfigFromFile("config.yaml")
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading config: %v\n", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Config: %+v\n", config))

	nodes, err := kubectl.GetNodes()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting nodes: %v\n", err))
		os.Exit(1)
	}

	for _, node := range nodes {
		slog.Info(fmt.Sprintf("Node %s: spot=%v, region=%v, zone=%s", node.Name(), node.IsSpotNode(), node.Region(), node.Zone()))
	}

	slog.Info(fmt.Sprintf(""))

	allNamespaces, err := kubectl.GetAllNamespaces()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting all namespaces: %v\n", err))
		os.Exit(1)
	}

	total := 0.0
	for _, ns := range allNamespaces {
		if !slices.Contains(config.Namespaces.Excluded, ns) {
			nsPrice := namespacePrice(config, ns, nodes)
			total += nsPrice
		}
	}

	slog.Info(fmt.Sprintf(""))
	slog.Info(fmt.Sprintf("-->> TOTAL PRICE: %f\n", total))

}

func podPrice(
	pod model.Pod,
	node model.Node,
	prices model.GkePrice,
) float64 {
	if node.IsSpotNode() {
		return pod.CPURequestCores()*prices.Spot.CPU + pod.MemoryRequestGB()*prices.Spot.RAM
	} else {
		return pod.CPURequestCores()*prices.Regular.CPU + pod.MemoryRequestGB()*prices.Regular.RAM
	}
}

func namespacePrice(
	config model.Config,
	namespace string,
	nodes []model.Node,
) float64 {

	pods, err := kubectl.GetPods(namespace)
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting pods: %v\n", err))
		os.Exit(1)
	}

	nodeLkup := make(map[string]model.Node)
	for _, node := range nodes {
		nodeLkup[node.Name()] = node
	}

	podPriceLkup := make(map[string]float64)
	for _, pod := range pods {
		podPriceLkup[pod.Name()] = podPrice(pod, nodeLkup[pod.NodeName()], config.Prices.GKE.Autopilot)
	}

	slices.SortFunc(pods, func(a, b model.Pod) int {
		return int(podPriceLkup[a.Name()] - podPriceLkup[b.Name()])
	})

	slog.Info(fmt.Sprintf("-----------PRICE FOR NAMESPACE %s------------", namespace))
	for _, pod := range pods {
		node := nodeLkup[pod.NodeName()]
		slog.Info(fmt.Sprintf(" + pod %s: spot=%v, cpu=%f, memory=%f => price=%f\n", pod.Name(), node.IsSpotNode(), pod.CPURequestCores(), pod.MemoryRequestGB(), podPriceLkup[pod.Name()]))
	}
	priceSum := 0.0
	for _, pod := range pods {
		priceSum += podPriceLkup[pod.Name()]
	}
	slog.Info(fmt.Sprintf(" = %f\n", priceSum))

	return priceSum
}
