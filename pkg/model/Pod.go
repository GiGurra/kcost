package model

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

func (p Pod) CPURequest() string {
	if len(p.Spec.Containers) == 0 {
		return ""
	}
	return p.Spec.Containers[0].Resources.Requests.CPU
}

func (p Pod) MemoryRequest() string {
	if len(p.Spec.Containers) == 0 {
		return ""
	}
	return p.Spec.Containers[0].Resources.Requests.Memory
}
