package extender

import (
	v1 "k8s.io/api/core/v1"
)

// ExtenderArgs matches kube-scheduler extender request body for Filter verb.
// See: https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/#scheduler-extensions

type ExtenderArgs struct {
	Pod   *v1.Pod   `json:"pod"`
	Nodes *v1.NodeList `json:"nodes,omitempty"`
	// NodeNames is used if Nodes is nil
	NodeNames []string `json:"nodeNames,omitempty"`
}

// ExtenderFilterResult is the response for Filter verb.
// Filtered nodes will be considered for scheduling; Failed nodes carry failure messages.

type ExtenderFilterResult struct {
	Nodes     *v1.NodeList   `json:"nodes,omitempty"`
	NodeNames []string       `json:"nodeNames,omitempty"`
	FailedNodes map[string]string `json:"failedNodes,omitempty"`
	Error     string         `json:"error,omitempty"`
}
