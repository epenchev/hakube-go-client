package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Amount of time which K8s allows running Node to be unresponsive before marking it unhealthy.
	// This is the default setting as per the documentation for kube-controller-manager
	// https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/
	NodeMonitorGracePeriodSeconds = 40
	// SafeTimeToAssumeNodeRebootedSeconds is the time after which the healthy self node remediation
	// agents will assume the unhealthy node has been rebooted, and it is safe to recover affected workloads.
	// As quoted from the selfnoderemediationconfig_types.go .
	SafeTimeToAssumeNodeRebootedSeconds = 180
)

// HAKubeControllerSpec defines the desired state of a HAKubeController
type HAKubeControllerSpec struct {
	// +optional
	Metadata *Metadata `json:"metadata,omitempty"`

	// NodeSelector targeting specific worker Nodes.
	// Will deploy the worker agents on the target Nodes.
	// +optional
	// +kubebuilder:default={node-role.kubernetes.io/worker: saphana}
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Global cluster fencing is on by default.
	// Also failed nodes running resources will be fenced by the fencing/remediation operator.
	// +optional
	// +kubebuilder:default=true
	Fencing *bool `json:"fencing,omitempty"`

	// Whether or not the HAKubeController should be stopped.
	// When this is true all resources are stopped as well.
	// +optional
	// +kubebuilder:default=false
	Shutdown *bool `json:"shutdown,omitempty"`

	// Whether to create additional node attributes storage object.
	// If a resource with multiple replicas is configured,
	// an node attributes object storage is created unless if explicitly disabled.
	// +optional
	// +kubebuilder:default=true
	AttributeStorage *bool `json:"attributeStorage,omitempty"`

	// One or more resources the high availability resource manager is responsible for.
	// +optional
	// +kubebuilder:default:{}
	// +listType=map
	// +listMapKey=name
	Resources []HAKubeResourceSpec `json:"resources"`

	// Defines the start sequence of the resources previously defined in the 'Resources' section,
	// If there is a dependency constraints for a resource to start before other resources, otherwise leave empty.
	// +optional
	// +kubebuilder:default:{}
	StartSequence []string `json:"startSequence"`
}

// HAKubeStatus defines the observed state of HAKube cluster
type HAKubeControllerStatus struct {
	// Current state of HAKube resources.
	// +optional
	// +kubebuilder:default:{}
	// +listType=map
	// +listMapKey=uid
	Resources []HAKubeResourceStatus `json:"resources"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=hkb
//+kubebuilder:subresource:status

// HAKubeController represents the 'blueprint' or in other words the configuration
// of the high availability resource/service manager controller together with the resources/services managed.
type HAKubeController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HAKubeControllerSpec   `json:"spec,omitempty"`
	Status HAKubeControllerStatus `json:"status,omitempty"`
}

// Metadata contains metadata for HAKubeController resources
type Metadata struct {
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Default defines several key default values for a HAKubeController configuration.
func (hactrl *HAKubeController) Default() {
	for ridx, _ := range hactrl.Spec.Resources {
		resrSpec := &hactrl.Spec.Resources[ridx]
		resrSpec.Default()
	}

	// allocate startup sequnce list
	if hactrl.Spec.StartSequence == nil {
		hactrl.Spec.StartSequence = []string{}
	}
	// allocate resources list
	if hactrl.Spec.Resources == nil {
		hactrl.Spec.Resources = []HAKubeResourceSpec{}
	}
	// allocate status
	if hactrl.Status.Resources == nil {
		hactrl.Status.Resources = []HAKubeResourceStatus{}
	}
}

// Gets labels from a Metadata pointer, if Metadata
// hasn't been set return nil.
func (meta *Metadata) GetLabelsOrNil() map[string]string {
	if meta == nil {
		return nil
	}
	return meta.Labels
}

// Gets annotations from a Metadata pointer, if Metadata
// hasn't been set return nil.
func (meta *Metadata) GetAnnotationsOrNil() map[string]string {
	if meta == nil {
		return nil
	}
	return meta.Annotations
}

//+kubebuilder:object:root=true

// HAKubeControllerList contains a list of HAKube controllers
type HAKubeControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HAKubeController `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HAKubeController{}, &HAKubeControllerList{})
}
