package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Keep attribute forever.
	AttributeLifetimeForever = "forever"

	// Keep attribute until node reboots.
	AttributeLifetimeReboot = "reboot"
)

type NodeAttribute struct {
	// Attribute name.
	// +required
	Name string `json:"name"`

	// Attribute value.
	// +required
	Value string `json:"value"`

	// Node name.
	// +required
	Node string `json:"node"`

	// How long to keep the attribute
	// +required
	// +kubebuilder:default=forever
	Lifetime string `json:"lifetime"`
}

// Defines default values for the HAKubeNodeAttributes object.
func (hkbnattr *HAKubeNodeAttributes) Default() {
	// allocate node attributes.
	if hkbnattr.NodeAttributes == nil {
		hkbnattr.NodeAttributes = []NodeAttribute{}
	}
	if hkbnattr.NodeList == nil {
		hkbnattr.NodeList = []string{}
	}
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=hkbna;hkbnodeattr

// HAKubeNodeAttributes holds additional data in the form of Node attributes objects for the RA to read and write.
type HAKubeNodeAttributes struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Node list targeting specific set of worker Nodes in the cluster using the object to store attributes.
	// +optional
	// +kubebuilder:default:{}
	NodeList []string `json:"nodeList,omitempty"`

	// Special values to be set using attributes.
	// An attribute has a name, and may have a different value for each node (node-attribute).
	// +optional
	// +kubebuilder:default:{}
	NodeAttributes []NodeAttribute `json:"nodeAttributes,omitempty"`
}

//+kubebuilder:object:root=true

// HAKubeNodeAttributes contains a list of HAKubeNodeAttributes objects
type HAKubeNodeAttributesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HAKubeNodeAttributes `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HAKubeNodeAttributes{}, &HAKubeNodeAttributesList{})
}
