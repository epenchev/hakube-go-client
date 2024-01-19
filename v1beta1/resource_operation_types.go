package v1beta1

type ResourceOperationType string

const (
	// Max retry attempts in case of operation failure
	OperationMaxRetry = 5

	// How many seconds to wait before each operation retry.
	OperationRetryInterval = 10

	// Max retry threshold time in seconds (how long a failed operation could retry).
	OperationRetryThreshold = 300
)

// Referring to https://clusterlabs.org/pacemaker/doc/2.1/Pacemaker_Explained/html/resources.html#operation-properties
const (
	// Start the resource.
	OperationStart ResourceOperationType = "start"

	// Stop the resource.
	OperationStop ResourceOperationType = "stop"

	// Perform monitor on the resource.
	OperationMonitor ResourceOperationType = "monitor"

	// Notify is used to inform the resource in advance for the next operation.
	OperationNotify ResourceOperationType = "notify"

	// Promote a resource in a Primary/Secondary scenario.
	OperationPromote ResourceOperationType = "promote"

	// Demote a resource in a Primary/Secondary scenario.
	OperationDemote ResourceOperationType = "demote"
)

type ResourceActionType string

const (
	// Stop the resource and start it again (possibly on a different node).
	ActionRestart ResourceActionType = "restart"

	// Fence the node on which the resource failed.
	ActionFence ResourceActionType = "fence"

	// Ignore action, nothing to do (pretend the resource did not fail if used in onFail).
	ActionIgnore ResourceActionType = "ignore"

	// Don’t perform any further operations on the resource.
	ActionBlock ResourceActionType = "block"

	// Stop the resource and do not start it elsewhere.
	ActionStop ResourceActionType = "stop"

	// Move all resources away from the node on which the resource failed.
	ActionStandby ResourceActionType = "standby"

	// Demote the resource instance, without a full restart.
	// This is valid only for promote actions, and for monitor actions with a nonzero interval and for a multi-state resource,
	// for any other action, a configuration error will be logged, and the default behaviour will be used.
	// Still unsupported.
	ActionDemote ResourceActionType = "demote"
)

// Operations are actions the resource manager can perform on a resource by executing a resource agent program.
type ResourceOperationSpec struct {
	// The operation to perform.
	// +kubebuilder:validation:Enum={start,stop,monitor,notify,promote,demote}
	// +required
	// +kubebuilder:default=notify
	Operation ResourceOperationType `json:"operation,omitempty"`

	// How long (in seconds) to wait before declaring the operation has failed.
	// +optional
	// +kubebuilder:default=5
	Timeout uint32 `json:"timeout,omitempty"`

	// How frequently (in seconds) to perform the operation.
	// A value of 0 means never.
	// A positive value defines a recurring action, which is typically used with monitor.
	// +optional
	// +kubebuilder:default=0
	Interval uint32 `json:"interval,omitempty"`

	// The action to take if this operation fails.
	// +kubebuilder:validation:Enum={restart,fence,ignore,suicide,block,stop,standby}
	// +optional
	// +kubebuilder:default=restart
	OnFail ResourceActionType `json:"onfail,omitempty"`

	// Disable the operation
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
}
