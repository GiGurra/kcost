package main

import (
	"fmt"
	"github.com/gigurra/kcost/pkg/kubectl"
	"github.com/gigurra/kcost/pkg/log"
	"github.com/gigurra/kcost/pkg/model"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"slices"
)

func main() {
	_, err := exec.LookPath("kubectl")
	if err != nil {
		log.ErrLn("kubectl not found")
		os.Exit(1)
	}

	config, err := model.NewConfigFromFile("config.yaml")
	if err != nil {
		log.ErrLn(fmt.Sprintf("error reading config: %v", err))
		os.Exit(1)
	}

	log.OutLn("-- Config --")
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		log.ErrLn(fmt.Sprintf("error marshalling config: %v", err))
		os.Exit(1)
	}
	log.OutLn(fmt.Sprintf("%s", string(configYaml)))

	nodes, err := kubectl.GetNodes()
	if err != nil {
		log.ErrLn(fmt.Sprintf("error getting nodes: %v", err))
		os.Exit(1)
	}

	log.OutLn("-- Nodes --")
	for _, node := range nodes {
		log.OutLn(fmt.Sprintf(" * Node %s: spot=%v, region=%v, zone=%s", node.Name(), node.IsSpotNode(), node.Region(), node.Zone()))
	}
	log.OutLn("")

	log.OutLn("-- Namespaces included --")

	includedNamespaces, err := getIncludedNamespaces(config)
	if err != nil {
		log.ErrLn(fmt.Sprintf("error getting included namespaces: %v", err))
		os.Exit(1)
	}

	for _, ns := range includedNamespaces {
		log.OutLn(fmt.Sprintf(" * %s", ns))
	}

	log.OutLn("")

	log.OutLn("-- Results --")
	total := 0.0
	for _, ns := range includedNamespaces {
		total += namespacePrice(config, ns, nodes)
	}

	log.OutLn(fmt.Sprintf(""))
	log.OutLn(fmt.Sprintf("-->> TOTAL PRICE: %f", total))

}

func getIncludedNamespaces(config model.Config) ([]string, error) {

	allNamespaces, err := kubectl.GetAllNamespaces()

	if err != nil {
		log.ErrLn(fmt.Sprintf("error getting all namespaces: %v", err))
		os.Exit(1)
	}

	result := []string{}
	for _, ns := range allNamespaces {
		if !slices.Contains(config.Namespaces.Excluded, ns) {
			result = append(result, ns)
		}
	}

	return result, err
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
		log.ErrLn(fmt.Sprintf("Error getting pods: %v", err))
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

	log.OutLn(fmt.Sprintf(" <> Price for namespace %s:", namespace))
	for _, pod := range pods {
		node := nodeLkup[pod.NodeName()]
		log.OutLn(fmt.Sprintf("  + pod %s: spot=%v, cpu=%f, memory=%f => price=%f", pod.Name(), node.IsSpotNode(), pod.CPURequestCores(), pod.MemoryRequestGB(), podPriceLkup[pod.Name()]))
	}
	priceSum := 0.0
	for _, pod := range pods {
		priceSum += podPriceLkup[pod.Name()]
	}
	log.OutLn(fmt.Sprintf("  = %f", priceSum))

	return priceSum
}
