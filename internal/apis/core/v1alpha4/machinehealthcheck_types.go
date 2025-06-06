/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha4

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ANCHOR: MachineHealthCheckSpec

// MachineHealthCheckSpec defines the desired state of MachineHealthCheck.
type MachineHealthCheckSpec struct {
	// clusterName is the name of the Cluster this object belongs to.
	// +kubebuilder:validation:MinLength=1
	ClusterName string `json:"clusterName"`

	// selector is the label selector to match machines whose health will be exercised
	Selector metav1.LabelSelector `json:"selector"`

	// unhealthyConditions contains a list of the conditions that determine
	// whether a node is considered unhealthy.  The conditions are combined in a
	// logical OR, i.e. if any of the conditions is met, the node is unhealthy.
	//
	// +kubebuilder:validation:MinItems=1
	UnhealthyConditions []UnhealthyCondition `json:"unhealthyConditions"`

	// maxUnhealthy specifies the maximum number of unhealthy machines allowed.
	// Any further remediation is only allowed if at most "maxUnhealthy" machines selected by
	// "selector" are not healthy.
	// +optional
	MaxUnhealthy *intstr.IntOrString `json:"maxUnhealthy,omitempty"`

	// unhealthyRange specifies the range of unhealthy machines allowed.
	// Any further remediation is only allowed if the number of machines selected by "selector" as not healthy
	// is within the range of "unhealthyRange". Takes precedence over maxUnhealthy.
	// Eg. "[3-5]" - This means that remediation will be allowed only when:
	// (a) there are at least 3 unhealthy machines (and)
	// (b) there are at most 5 unhealthy machines
	// +optional
	// +kubebuilder:validation:Pattern=^\[[0-9]+-[0-9]+\]$
	UnhealthyRange *string `json:"unhealthyRange,omitempty"`

	// nodeStartupTimeout is the duration after which machines without a node will be considered to
	// have failed and will be remediated.
	// If not set, this value is defaulted to 10 minutes.
	// If you wish to disable this feature, set the value explicitly to 0.
	// +optional
	NodeStartupTimeout *metav1.Duration `json:"nodeStartupTimeout,omitempty"`

	// remediationTemplate is a reference to a remediation template
	// provided by an infrastructure provider.
	//
	// This field is completely optional, when filled, the MachineHealthCheck controller
	// creates a new object from the template referenced and hands off remediation of the machine to
	// a controller that lives outside of Cluster API.
	// +optional
	RemediationTemplate *corev1.ObjectReference `json:"remediationTemplate,omitempty"`
}

// ANCHOR_END: MachineHealthCHeckSpec

// ANCHOR: UnhealthyCondition

// UnhealthyCondition represents a Node condition type and value with a timeout
// specified as a duration.  When the named condition has been in the given
// status for at least the timeout value, a node is considered unhealthy.
type UnhealthyCondition struct {
	// type of Node condition
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Type corev1.NodeConditionType `json:"type"`

	// status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Status corev1.ConditionStatus `json:"status"`

	// timeout is the duration that a node must be in a given status for,
	// after which the node is considered unhealthy.
	// For example, with a value of "1h", the node must match the status
	// for at least 1 hour before being considered unhealthy.
	Timeout metav1.Duration `json:"timeout"`
}

// ANCHOR_END: UnhealthyCondition

// ANCHOR: MachineHealthCheckStatus

// MachineHealthCheckStatus defines the observed state of MachineHealthCheck.
type MachineHealthCheckStatus struct {
	// expectedMachines is the total number of machines counted by this machine health check
	// +kubebuilder:validation:Minimum=0
	ExpectedMachines int32 `json:"expectedMachines,omitempty"`

	// currentHealthy is the total number of healthy machines counted by this machine health check
	// +kubebuilder:validation:Minimum=0
	CurrentHealthy int32 `json:"currentHealthy,omitempty"`

	// remediationsAllowed is the number of further remediations allowed by this machine health check before
	// maxUnhealthy short circuiting will be applied
	// +kubebuilder:validation:Minimum=0
	RemediationsAllowed int32 `json:"remediationsAllowed,omitempty"`

	// observedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// targets shows the current list of machines the machine health check is watching
	// +optional
	Targets []string `json:"targets,omitempty"`

	// conditions defines current service state of the MachineHealthCheck.
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`
}

// ANCHOR_END: MachineHealthCheckStatus

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion
// +kubebuilder:resource:path=machinehealthchecks,shortName=mhc;mhcs,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterName",description="Cluster"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of MachineHealthCheck"
// +kubebuilder:printcolumn:name="MaxUnhealthy",type="string",JSONPath=".spec.maxUnhealthy",description="Maximum number of unhealthy machines allowed"
// +kubebuilder:printcolumn:name="ExpectedMachines",type="integer",JSONPath=".status.expectedMachines",description="Number of machines currently monitored"
// +kubebuilder:printcolumn:name="CurrentHealthy",type="integer",JSONPath=".status.currentHealthy",description="Current observed healthy machines"

// MachineHealthCheck is the Schema for the machinehealthchecks API.
//
// Deprecated: This type will be removed in one of the next releases.
type MachineHealthCheck struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the specification of machine health check policy
	Spec MachineHealthCheckSpec `json:"spec,omitempty"`

	// status is the most recently observed status of MachineHealthCheck resource
	Status MachineHealthCheckStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (m *MachineHealthCheck) GetConditions() Conditions {
	return m.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (m *MachineHealthCheck) SetConditions(conditions Conditions) {
	m.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// MachineHealthCheckList contains a list of MachineHealthCheck.
//
// Deprecated: This type will be removed in one of the next releases.
type MachineHealthCheckList struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#lists-and-simple-kinds
	metav1.ListMeta `json:"metadata,omitempty"`
	// items is the list of MachineHealthChecks.
	Items []MachineHealthCheck `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &MachineHealthCheck{}, &MachineHealthCheckList{})
}
