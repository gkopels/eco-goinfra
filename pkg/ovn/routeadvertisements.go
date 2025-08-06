package ovn

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	routeadvertisementstypes "github.com/openshift-kni/eco-goinfra/pkg/schemes/ovnk8s/routeadvertisementstypes"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// RouteAdvertisementsBuilder provides struct for the RouteAdvertisements object containing connection to
// the cluster and the RouteAdvertisements definitions.
type RouteAdvertisementsBuilder struct {
	Definition *routeadvertisementstypes.RouteAdvertisements
	Object     *routeadvertisementstypes.RouteAdvertisements
	apiClient  runtimeClient.Client
	errorMsg   string
}

// NewRouteAdvertisementsBuilder creates a new instance of RouteAdvertisements.
func NewRouteAdvertisementsBuilder(
	apiClient *clients.Settings, name, nsname string) *RouteAdvertisementsBuilder {
	glog.V(100).Infof(
		"Initializing new RouteAdvertisements structure with the following params: %s, %s",
		name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("failed to initialize the apiclient is empty")

		return nil
	}

	err := apiClient.AttachScheme(routeadvertisementstypes.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add routeadvertisements scheme to client schemes")

		return nil
	}

	builder := &RouteAdvertisementsBuilder{
		apiClient: apiClient,
		Definition: &routeadvertisementstypes.RouteAdvertisements{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the RouteAdvertisements is empty")

		builder.errorMsg = "RouteAdvertisements 'name' cannot be empty"

		return builder
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the RouteAdvertisements is empty")

		builder.errorMsg = "RouteAdvertisements 'nsname' cannot be empty"

		return builder
	}

	return builder
}

// Exists checks whether the given RouteAdvertisements exists.
func (builder *RouteAdvertisementsBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof(
		"Checking if RouteAdvertisements %s exists in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// Create makes a RouteAdvertisements in the cluster and stores the created object in struct.
func (builder *RouteAdvertisementsBuilder) Create() (*RouteAdvertisementsBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating the RouteAdvertisements %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace,
	)

	var err error
	if !builder.Exists() {
		err = builder.apiClient.Create(context.TODO(), builder.Definition)

		if err != nil {
			glog.V(100).Infof("Failed to create RouteAdvertisements")

			return nil, err
		}

		builder.Object = builder.Definition
	}

	return builder, nil
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *RouteAdvertisementsBuilder) validate() (bool, error) {
	resourceCRD := "routeadvertisements"

	if builder == nil {
		glog.V(100).Infof("The %s builder is uninitialized", resourceCRD)

		return false, fmt.Errorf("error: received nil %s builder", resourceCRD)
	}

	if builder.Definition == nil {
		glog.V(100).Infof("The %s is undefined", resourceCRD)

		return false, fmt.Errorf("%s", msg.UndefinedCrdObjectErrString(resourceCRD))
	}

	if builder.apiClient == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceCRD)

		return false, fmt.Errorf("%s builder cannot have nil apiClient", resourceCRD)
	}

	if builder.errorMsg != "" {
		glog.V(100).Infof("the %s builder has error message: %s", resourceCRD, builder.errorMsg)

		return false, fmt.Errorf("%s", builder.errorMsg)
	}

	return true, nil
}

// Get returns RouteAdvertisements object if found.
func (builder *RouteAdvertisementsBuilder) Get() (*routeadvertisementstypes.RouteAdvertisements, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof(
		"Collecting RouteAdvertisements object %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	routeAdvertisements := &routeadvertisementstypes.RouteAdvertisements{}
	err := builder.apiClient.Get(context.TODO(),
		runtimeClient.ObjectKey{Name: builder.Definition.Name, Namespace: builder.Definition.Namespace}, routeAdvertisements)

	if err != nil {
		glog.V(100).Infof(
			"RouteAdvertisements object %s does not exist in namespace %s",
			builder.Definition.Name, builder.Definition.Namespace)

		return nil, err
	}

	return routeAdvertisements, nil
}

// PullRouteAdvertisements pulls existing RouteAdvertisements from cluster.
func PullRouteAdvertisements(
	apiClient *clients.Settings, name, nsname string) (*RouteAdvertisementsBuilder, error) {
	glog.V(100).Infof(
		"Pulling existing RouteAdvertisements name %s under namespace %s from cluster", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is empty")

		return nil, fmt.Errorf("RouteAdvertisements 'apiClient' cannot be empty")
	}

	err := apiClient.AttachScheme(routeadvertisementstypes.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add RouteAdvertisements scheme to client schemes")

		return nil, err
	}

	builder := &RouteAdvertisementsBuilder{
		apiClient: apiClient,
		Definition: &routeadvertisementstypes.RouteAdvertisements{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the RouteAdvertisements is empty")

		return nil, fmt.Errorf("RouteAdvertisements 'name' cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the RouteAdvertisements is empty")

		return nil, fmt.Errorf("RouteAdvertisements 'namespace' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("RouteAdvertisements object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return builder, nil
}

// Delete removes RouteAdvertisements object from a cluster.
func (builder *RouteAdvertisementsBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the RouteAdvertisements object %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace,
	)

	if !builder.Exists() {
		glog.V(100).Infof(
			"RouteAdvertisements %s in namespace %s does not exist",
			builder.Definition.Name, builder.Definition.Namespace)

		builder.Object = nil

		return nil
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return fmt.Errorf("can not delete RouteAdvertisements: %w", err)
	}

	builder.Object = nil

	return nil
}

// WithAdvertisements sets the advertisements in the RouteAdvertisements spec.
func (builder *RouteAdvertisementsBuilder) WithAdvertisements(
	advertisements []string) *RouteAdvertisementsBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Creating RouteAdvertisements %s in namespace %s with advertisements: %v",
		builder.Definition.Name, builder.Definition.Namespace, advertisements)

	if len(advertisements) == 0 {
		glog.V(100).Infof("Can not redefine RouteAdvertisements with empty advertisements")

		builder.errorMsg = "RouteAdvertisements 'advertisements' cannot be empty"

		return builder
	}

	builder.Definition.Spec.Advertisements = advertisements

	return builder
}

// WithNodeSelector sets the nodeSelector in the RouteAdvertisements spec.
func (builder *RouteAdvertisementsBuilder) WithNodeSelector(
	nodeSelector map[string]string) *RouteAdvertisementsBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Setting RouteAdvertisements %s in namespace %s with nodeSelector: %v",
		builder.Definition.Name, builder.Definition.Namespace, nodeSelector)

	if len(nodeSelector) == 0 {
		glog.V(100).Infof("Can not set RouteAdvertisements with empty nodeSelector map")

		builder.errorMsg = "RouteAdvertisements 'nodeSelector' cannot be empty map"

		return builder
	}

	builder.Definition.Spec.NodeSelector = nodeSelector

	return builder
}

// WithFRRConfigurationSelector sets the frrConfigurationSelector in the RouteAdvertisements spec.
func (builder *RouteAdvertisementsBuilder) WithFRRConfigurationSelector(
	frrConfigurationSelector map[string]string) *RouteAdvertisementsBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Setting RouteAdvertisements %s in namespace %s with frrConfigurationSelector: %v",
		builder.Definition.Name, builder.Definition.Namespace, frrConfigurationSelector)

	builder.Definition.Spec.FRRConfigurationSelector = frrConfigurationSelector

	return builder
}

// WithNetworkSelector adds a network selector to the RouteAdvertisements spec.
func (builder *RouteAdvertisementsBuilder) WithNetworkSelector(
	networkSelectionType string, clusterUserDefinedNetworkSelector *routeadvertisementstypes.NetworkSelector) *RouteAdvertisementsBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Adding network selector to RouteAdvertisements %s in namespace %s with type: %s",
		builder.Definition.Name, builder.Definition.Namespace, networkSelectionType)

	if networkSelectionType == "" {
		glog.V(100).Infof("Can not add network selector with empty networkSelectionType")

		builder.errorMsg = "networkSelectionType cannot be empty"

		return builder
	}

	networkSelector := routeadvertisementstypes.NetworkSelectorSpec{
		NetworkSelectionType: networkSelectionType,
	}

	if clusterUserDefinedNetworkSelector != nil {
		networkSelector.ClusterUserDefinedNetworkSelector = clusterUserDefinedNetworkSelector
	}

	builder.Definition.Spec.NetworkSelectors = append(builder.Definition.Spec.NetworkSelectors, networkSelector)

	return builder
}

// GetRouteAdvertisementsGVR returns RouteAdvertisements GroupVersionResource.
func GetRouteAdvertisementsGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group: "k8s.ovn.org", Version: "v1", Resource: "routeadvertisements",
	}
}
