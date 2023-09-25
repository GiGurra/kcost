package model

import (
	"fmt"
	"log/slog"
	"strings"
)

// We are only interested in the pod name, its node and its resource requests
type Pod struct {
	// Name and labels are found in the metadata section
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	// The node name and resource requests are found in the spec section
	Spec struct {
		NodeName   string `json:"nodeName"`
		Containers []struct {
			Resources struct {
				Requests struct {
					CPU    string `json:"cpu"`
					Memory string `json:"memory"`
				} `json:"requests"`
			} `json:"resources"`
		} `json:"containers"`
	} `json:"spec"`
}

func (p Pod) Name() string {
	return p.Metadata.Name
}

func (p Pod) NodeName() string {
	return p.Spec.NodeName
}

func (p Pod) CPURequestCores() float64 {
	if len(p.Spec.Containers) == 0 {
		slog.Warn(fmt.Sprintf("Pod '%s' has no containers\n", p.Name()))
		return 0.0
	}

	if len(p.Spec.Containers) > 1 {
		slog.Warn(fmt.Sprintf("Not supported! Pod '%s' has multiple containers. Using first only\n", p.Name()))
	}

	str := p.Spec.Containers[0].Resources.Requests.CPU

	// Parse str into a float64
	// str is a string like "100m" or "1"
	// We want to convert it to a float64 like 0.1 or 1.0
	// if the string ends with "m" then we divide by 1000
	// otherwise treat as a float64
	if len(str) > 0 && str[len(str)-1] == 'm' {
		str = str[:len(str)-1]
		str = "0." + str
	} else {
		str = str + ".0"
	}

	// parse into float64
	var f float64
	_, err := fmt.Sscanf(str, "%f", &f)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error parsing CPU request '%s': %v\n", str, err))
	}

	return f
}

func parseMem(str string, divisor float64) float64 {

	// parse into float64
	var f float64
	_, err := fmt.Sscanf(str, "%f", &f)
	if err != nil {
		slog.Warn(fmt.Sprintf("Error parsing memory request '%s': %v\n", str, err))
		return 0.0
	}

	return f / divisor
}

func (p Pod) MemoryRequestGB() float64 {
	if len(p.Spec.Containers) == 0 {
		slog.Warn(fmt.Sprintf("Pod '%s' has no containers\n", p.Name()))
		return 0.0
	}

	if len(p.Spec.Containers) > 1 {
		slog.Warn(fmt.Sprintf("Not supported! Pod '%s' has multiple containers. Using first only\n", p.Name()))
	}

	str := p.Spec.Containers[0].Resources.Requests.Memory
	if strings.HasSuffix(str, "Gi") {
		return parseMem(str[:len(str)-2], 1.0)
	} else if strings.HasSuffix(str, "Mi") {
		return parseMem(str[:len(str)-2], 1024.0)
	} else if strings.HasSuffix(str, "Ki") {
		return parseMem(str[:len(str)-2], 1024.0*1024.0)
	} else {
		return parseMem(str[:len(str)-2], 1024.0*1024.0*1024.0)
	}
}
