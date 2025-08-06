package routeadvertisementstypes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteAdvertisementsSpec defines the desired state of RouteAdvertisements
type RouteAdvertisementsSpec struct {
	// Advertisements specifies what types of routes to advertise.
	// +kubebuilder:validation:Required
	Advertisements []string `json:"advertisements"`

	// FRRConfigurationSelector is a label selector for FRR configurations.
	// +optional
	FRRConfigurationSelector map[string]string `json:"frrConfigurationSelector,omitempty"`

	// NetworkSelectors specifies which networks to advertise routes for.
	// +optional
	NetworkSelectors []NetworkSelectorSpec `json:"networkSelectors,omitempty"`

	// NodeSelector is a label selector for nodes to advertise routes from.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// NetworkSelectorSpec defines a network selector
type NetworkSelectorSpec struct {
	// NetworkSelectionType specifies the type of network selection.
	// +kubebuilder:validation:Required
	NetworkSelectionType string `json:"networkSelectionType"`

	// ClusterUserDefinedNetworkSelector specifies selector for cluster user defined networks.
	// +optional
	ClusterUserDefinedNetworkSelector *NetworkSelector `json:"clusterUserDefinedNetworkSelector,omitempty"`
}

// NetworkSelector defines a network selector with match labels
type NetworkSelector struct {
	// NetworkSelector specifies the network selector with match labels.
	// +optional
	NetworkSelector *metav1.LabelSelector `json:"networkSelector,omitempty"`
}

// RouteAdvertisementsStatus defines the observed state of RouteAdvertisements
type RouteAdvertisementsStatus struct {
	// Conditions represents the current state of the RouteAdvertisements
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced

// RouteAdvertisements is the Schema for the routeadvertisements API
type RouteAdvertisements struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteAdvertisementsSpec   `json:"spec,omitempty"`
	Status RouteAdvertisementsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RouteAdvertisementsList contains a list of RouteAdvertisements
type RouteAdvertisementsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RouteAdvertisements `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RouteAdvertisements{}, &RouteAdvertisementsList{})
}
