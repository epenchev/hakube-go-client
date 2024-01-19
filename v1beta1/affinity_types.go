package v1beta1

// A way to apply resource constrains is to use K8s builtin functionalities like Node Affinity, Inter-Pod Affinity or even NodeSelector.
// Those work well with fully containerised workloads/resources.
// But the case of managing resources on BareMetal ot VM's is a bit different.
// When a scheduling constraint condition is changed resources running on BareMetal or VM's aren't stopped by moving or deleting Pods.
// This requires for a custom solution providing a basic set of constrains.

// A resource status-condition term for selecting a worker node.
type NodeSelectorTerm struct {
	// Target resource to be matched from the resources list defined within the HAKubeController.
	// +required
	Resource string `json:"resource"`

	// Resource status condition that matches a resource clone/replica.
	// +required
	Condition ResourceConditionType `json:"condition"`
}

// ResourceNodeAffinity defines terms that constrain the selection of the worker node that resource will be scheduled to run on.
type ResourceNodeAffinity struct {
	// Selector specifies hard node constraints that must be met.
	// +required
	Selector *NodeSelectorTerm `json:"selector"`

	// If set, reverse will act as anti-affinity rule and the selected node is not used for scheduling the resource.
	// +optional
	// +kubebuilder:default=false
	Reverse *bool `json:"match,omitempty"`
}
