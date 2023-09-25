package model

type K8sListing[T any] struct {
	// Items is the list of nodes
	Items []T `json:"items"`
}
