package model

type Node struct {
	// Name and labels are found in the metadata section
	Metadata struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

func (n Node) Name() string {
	return n.Metadata.Name
}

func (n Node) GetLabels() map[string]string {
	return n.Metadata.Labels
}

func (n Node) IsSpotNode() bool {
	value, ok := n.Metadata.Labels["cloud.google.com/gke-spot"]
	return ok && value == "true"
}

func (n Node) Region() string {
	return n.Metadata.Labels["failure-domain.beta.kubernetes.io/region"]
}

func (n Node) Zone() string {
	return n.Metadata.Labels["failure-domain.beta.kubernetes.io/zone"]
}

func (n Node) Tpe() string {
	return n.Metadata.Labels["beta.kubernetes.io/instance-type"]
}
