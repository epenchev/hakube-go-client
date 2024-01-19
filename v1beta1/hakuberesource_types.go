package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceConditionType string

// Resource Condition Types
const (
	// Is the condition type used when Resource has started.
	ConditionTypeStarted ResourceConditionType = "started"

	// Is the condition type used when Resource has been stopped.
	ConditionTypeStopped ResourceConditionType = "stopped"

	// Is the condition type used when Resource has been fenced/remediated.
	ConditionTypeFenced ResourceConditionType = "fenced"

	// Is the condition type used when Resource has failed.
	// Resource can be set to failed when a monitor operation fails for example.
	ConditionTypeFailed ResourceConditionType = "failed"

	// Is the condition type used when Resource has been promoted, set to Promoted/Primary role.
	// When resource has multiple replicas and one of the resplicas is been promoted to be the Primary.
	ConditionTypePromoted ResourceConditionType = "promoted"

	// Is the condition type used when Resource has been demoted (removed Primary/Promoted role) and in set to Demoted role.
	ConditionTypeDemoted ResourceConditionType = "demoted"
)

// HAKubeResourceSpec defines the 'blueprint' or in other words configuration for a resource/service,
// that is managed from the high availability resource/service manager.
type HAKubeResourceSpec struct {
	// Name is the unique name of the resource.
	// +required
	Name string `json:"name"`

	// Entry point command for the resource agent.
	// +required
	Command []string `json:"command"`

	// An additional metadata for the resource.
	// +optional
	Metadata *Metadata `json:"metadata,omitempty"`

	// How many resource replicas/clones should be started.
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:default=1
	Replicas uint32 `json:"replicas,omitempty"`

	// Operations are actions the resource manager performs on a resource by executing a resource agent program.
	// Resource agents support a common set of operations such as start, stop and monitor, and may implement others.
	// +optional
	// +kubebuilder:default:{}
	// +listType=map
	// +listMapKey=operation
	Operations []ResourceOperationSpec `json:"operations,omitempty"`

	// Workers are K8s worker nodes.
	// Specify which nodes are eligible to execute the resource,
	// or just add a constrain/limit for workers who will be running the resource.
	// +optional
	// +kubebuilder:default:{}
	Workers []string `json:"workers,omitempty"`

	// Additional resources options specific to the given resource.
	// +optional
	Options map[string]string `json:"options,omitempty"`

	// nodeAffinity defines terms that constrain the selection of the worker node a resource can be scheduled to run on.
	// +optional
	NodeAffinity *ResourceNodeAffinity `json:"nodeAffinity,omitempty"`

	// Disable/stop the resource.
	// +optional
	// +kubebuilder:default=false
	Disabled *bool `json:"disabled,omitempty"`

	// Reserved but not implemented.
	// When stopping or starting resource instance,
	// notify all other replicas in advance and when the action was successful as well.
	// +optional
	// Notify *bool `json:"notify,omitempty"`
}

// HAKubeResourceStatus defines the observed state of HAKubeResourceSpec.
type HAKubeResourceStatus struct {
	// UID is the generated unique identifier for the resource status.
	// It is intended to distinguish between occurrences of similar entities in the controller's status list,
	// for example when resource has multiple replicas running on different nodes.
	// +kubebuilder:default=""
	UID string `json:"uid"`

	// Resource associated with the status.
	// +kubebuilder:default=""
	Resource string `json:"resource"`

	// conditions represent the observations of a resource current state
	// Known .status.conditions.type are: "ConditionResourceFailed"
	// +kubebuilder:default:{}
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Node associated with the worker Pod running the Resource.
	// A replicated/cloned Resource can be running on multiple Nodes.
	// +kubebuilder:default=""
	Node string `json:"node"`

	// Attempts to restart resource when failing.
	// +kubebuilder:default=0
	FailCount uint32 `json:"failCount"`
}

// A composite type representing HA resource with spec and status fields to be used from resource APIs for the resource management logic.
type HAKubeResource struct {
	Spec   HAKubeResourceSpec `json:"spec,omitempty"`
	Status []string           `json:"status,omitempty"`
}

// Default sets default values for the HAKubeResource resource operations attributes.
func (r *HAKubeResourceSpec) Default() {
	if r.Replicas < 1 {
		r.Replicas = 1
	}

	if r.Operations != nil {
		return
	}
	r.Operations = []ResourceOperationSpec{
		ResourceOperationSpec{
			Operation: OperationStart,
			OnFail:    ActionRestart,
			Timeout:   10,
		},
		ResourceOperationSpec{
			Operation: OperationStop,
			OnFail:    ActionBlock,
			Timeout:   10,
		},
		ResourceOperationSpec{
			Operation: OperationMonitor,
			OnFail:    ActionRestart,
			Timeout:   10,
			Interval:  10,
		},
	}
}
