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

	prices, err := model.NewPriceTableFromFile("prices.yaml")
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading prices: %v\n", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Prices: %+v\n", prices))

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

	nodeLkup := make(map[string]model.Node)
	for _, node := range nodes {
		nodeLkup[node.Name()] = node
	}

	podPriceLkup := make(map[string]float64)
	for _, pod := range pods {
		podPriceLkup[pod.Name()] = podPrice(pod, nodeLkup[pod.NodeName()], prices.GKE.Autopilot)
	}

	slices.SortFunc(pods, func(a, b model.Pod) int {
		return int(podPriceLkup[a.Name()] - podPriceLkup[b.Name()])
	})

	slog.Info("-------------------------")
	slog.Info("-------------------------")
	slog.Info("      PRICES PER POD	 ")
	slog.Info("-------------------------")
	for _, pod := range pods {
		node := nodeLkup[pod.NodeName()]
		slog.Info(fmt.Sprintf("Pod %s: spot=%v, cpu=%f, memory=%f => price=%f\n", pod.Name(), node.IsSpotNode(), pod.CPURequestCores(), pod.MemoryRequestGB(), podPriceLkup[pod.Name()]))
	}

	slog.Info("-------------------------")
	slog.Info("-------------------------")
	slog.Info("      PRICE FOR NS       ")
	slog.Info("-------------------------")
	priceSum := 0.0
	for _, pod := range pods {
		priceSum += podPriceLkup[pod.Name()]
	}
	slog.Info(fmt.Sprintf("Total price: %f\n", priceSum))
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
