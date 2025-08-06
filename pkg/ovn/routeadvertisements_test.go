package ovn

import (
	"fmt"
	"testing"

	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/schemes/ovnk8s/routeadvertisementstypes"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	routeAdvertisementsTestSchemes    = []clients.SchemeAttacher{routeadvertisementstypes.AddToScheme}
	defaultRouteAdvertisementsName    = "default"
	defaultRouteAdvertisementsNsName  = "test-namespace"
	defaultRouteAdvertisementsAdverts = []string{"PodNetwork"}
)

func TestNewRouteAdvertisementsBuilder(t *testing.T) {
	generateRouteAdvertisements := NewRouteAdvertisementsBuilder

	testCases := []struct {
		name          string
		namespace     string
		client        bool
		expectedError string
	}{
		{
			name:          defaultRouteAdvertisementsName,
			namespace:     defaultRouteAdvertisementsNsName,
			client:        true,
			expectedError: "",
		},
		{
			name:          "",
			namespace:     defaultRouteAdvertisementsNsName,
			client:        true,
			expectedError: "RouteAdvertisements 'name' cannot be empty",
		},
		{
			name:          defaultRouteAdvertisementsName,
			namespace:     "",
			client:        true,
			expectedError: "RouteAdvertisements 'nsname' cannot be empty",
		},
		{
			name:          defaultRouteAdvertisementsName,
			namespace:     defaultRouteAdvertisementsNsName,
			client:        false,
			expectedError: "",
		},
	}

	for _, testCase := range testCases {
		var testSettings *clients.Settings

		if testCase.client {
			testSettings = clients.GetTestClients(clients.TestClientParams{})
		}

		testRouteAdvertisementsBuilder := generateRouteAdvertisements(testSettings, testCase.name, testCase.namespace)

		if testCase.client {
			assert.Equal(t, testCase.expectedError, testRouteAdvertisementsBuilder.errorMsg)

			if testCase.expectedError == "" {
				assert.Equal(t, testCase.name, testRouteAdvertisementsBuilder.Definition.Name)
			}
		} else {
			assert.Nil(t, testRouteAdvertisementsBuilder)
		}
	}
}

func TestRouteAdvertisementsCreate(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		expectedError           error
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: nil,
		},
		{
			testRouteAdvertisements: buildInValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: fmt.Errorf("RouteAdvertisements 'name' cannot be empty"),
		},
	}

	for _, testCase := range testCases {
		testRouteAdvertisementsBuilder, err := testCase.testRouteAdvertisements.Create()
		assert.Equal(t, testCase.expectedError, err)

		if testCase.expectedError == nil {
			assert.Equal(t, testRouteAdvertisementsBuilder.Definition.Name, testRouteAdvertisementsBuilder.Object.Name)
		}
	}
}

func TestRouteAdvertisementsGet(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		expectedError           string
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: "",
		},
		{
			testRouteAdvertisements: buildInValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: "RouteAdvertisements 'name' cannot be empty",
		},
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(
				clients.GetTestClients(clients.TestClientParams{})),
			expectedError: "routeadvertisements.k8s.ovn.org \"default\" not found",
		},
	}

	for _, testCase := range testCases {
		routeAdvertisements, err := testCase.testRouteAdvertisements.Get()

		if testCase.expectedError == "" {
			assert.Nil(t, err)
			assert.Equal(t, routeAdvertisements.Name, testCase.testRouteAdvertisements.Definition.Name)
		} else {
			assert.EqualError(t, err, testCase.expectedError)
		}
	}
}

func TestPullRouteAdvertisements(t *testing.T) {
	generateRouteAdvertisements := func(name, namespace string) *routeadvertisementstypes.RouteAdvertisements {
		return &routeadvertisementstypes.RouteAdvertisements{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		name                string
		namespace           string
		addToRuntimeObjects bool
		expectedError       error
		client              bool
	}{
		{
			name:                defaultRouteAdvertisementsName,
			namespace:           defaultRouteAdvertisementsNsName,
			addToRuntimeObjects: true,
			expectedError:       nil,
			client:              true,
		},
		{
			name:                "",
			namespace:           defaultRouteAdvertisementsNsName,
			addToRuntimeObjects: true,
			expectedError:       fmt.Errorf("RouteAdvertisements 'name' cannot be empty"),
			client:              true,
		},
		{
			name:                defaultRouteAdvertisementsName,
			namespace:           "",
			addToRuntimeObjects: true,
			expectedError:       fmt.Errorf("RouteAdvertisements 'namespace' cannot be empty"),
			client:              true,
		},
		{
			name:                defaultRouteAdvertisementsName,
			namespace:           defaultRouteAdvertisementsNsName,
			addToRuntimeObjects: false,
			expectedError: fmt.Errorf("RouteAdvertisements object default does not" +
				" exist in namespace test-namespace"),
			client: true,
		},
		{
			name:                defaultRouteAdvertisementsName,
			namespace:           defaultRouteAdvertisementsNsName,
			addToRuntimeObjects: true,
			expectedError:       fmt.Errorf("RouteAdvertisements 'apiClient' cannot be empty"),
			client:              false,
		},
	}

	for _, testCase := range testCases {
		// Pre-populate the runtime objects
		var runtimeObjects []runtime.Object

		var testSettings *clients.Settings

		testRouteAdvertisements := generateRouteAdvertisements(testCase.name, testCase.namespace)

		if testCase.addToRuntimeObjects {
			runtimeObjects = append(runtimeObjects, testRouteAdvertisements)
		}

		if testCase.client {
			testSettings = clients.GetTestClients(clients.TestClientParams{
				K8sMockObjects:  runtimeObjects,
				SchemeAttachers: routeAdvertisementsTestSchemes,
			})
		}

		builderResult, err := PullRouteAdvertisements(testSettings, testCase.name, testCase.namespace)
		assert.Equal(t, testCase.expectedError, err)

		if testCase.expectedError == nil {
			assert.Equal(t, testCase.name, builderResult.Object.Name)
			assert.Equal(t, testCase.namespace, builderResult.Object.Namespace)
		}
	}
}

func TestRouteAdvertisementsExist(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		expectedStatus          bool
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedStatus: true,
		},
		{
			testRouteAdvertisements: buildInValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedStatus: false,
		},
	}

	for _, testCase := range testCases {
		exist := testCase.testRouteAdvertisements.Exists()
		assert.Equal(t, testCase.expectedStatus, exist)
	}
}

func TestRouteAdvertisementsDelete(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		expectedError           error
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: nil,
		},
		{
			testRouteAdvertisements: buildInValidRouteAdvertisementsBuilder(
				buildRouteAdvertisementsTestClientWithDummyObject()),
			expectedError: fmt.Errorf("RouteAdvertisements 'name' cannot be empty"),
		},
	}

	for _, testCase := range testCases {
		err := testCase.testRouteAdvertisements.Delete()
		assert.Equal(t, testCase.expectedError, err)

		if testCase.expectedError == nil {
			assert.Nil(t, testCase.testRouteAdvertisements.Object)
		}
	}
}

func TestRouteAdvertisementsWithAdvertisements(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		advertisements          []string
		expectedError           string
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			advertisements:          []string{"PodNetwork", "ServiceNetwork"},
		},
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			advertisements:          []string{},
			expectedError:           "RouteAdvertisements 'advertisements' cannot be empty",
		},
	}

	for _, testCase := range testCases {
		routeAdvertisementsBuilder := testCase.testRouteAdvertisements.WithAdvertisements(testCase.advertisements)
		assert.Equal(t, testCase.expectedError, routeAdvertisementsBuilder.errorMsg)

		if testCase.expectedError == "" {
			assert.Equal(t, testCase.advertisements, routeAdvertisementsBuilder.Definition.Spec.Advertisements)
		}
	}
}

func TestRouteAdvertisementsWithNodeSelector(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements *RouteAdvertisementsBuilder
		nodeSelector            map[string]string
		expectedError           string
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			nodeSelector:            map[string]string{"node-role": "worker"},
		},
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			nodeSelector:            map[string]string{},
			expectedError:           "RouteAdvertisements 'nodeSelector' cannot be empty map",
		},
	}

	for _, testCase := range testCases {
		routeAdvertisementsBuilder := testCase.testRouteAdvertisements.WithNodeSelector(testCase.nodeSelector)
		assert.Equal(t, testCase.expectedError, routeAdvertisementsBuilder.errorMsg)

		if testCase.expectedError == "" {
			assert.Equal(t, testCase.nodeSelector, routeAdvertisementsBuilder.Definition.Spec.NodeSelector)
		}
	}
}

func TestRouteAdvertisementsWithNetworkSelector(t *testing.T) {
	testCases := []struct {
		testRouteAdvertisements           *RouteAdvertisementsBuilder
		networkSelectionType              string
		clusterUserDefinedNetworkSelector *routeadvertisementstypes.NetworkSelector
		expectedError                     string
	}{
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			networkSelectionType:    "DefaultNetwork",
		},
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			networkSelectionType:    "ClusterUserDefinedNetworks",
			clusterUserDefinedNetworkSelector: &routeadvertisementstypes.NetworkSelector{
				NetworkSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"cudn-bgp": "red"},
				},
			},
		},
		{
			testRouteAdvertisements: buildValidRouteAdvertisementsBuilder(buildRouteAdvertisementsTestClientWithDummyObject()),
			networkSelectionType:    "",
			expectedError:           "networkSelectionType cannot be empty",
		},
	}

	for _, testCase := range testCases {
		routeAdvertisementsBuilder := testCase.testRouteAdvertisements.WithNetworkSelector(
			testCase.networkSelectionType, testCase.clusterUserDefinedNetworkSelector)
		assert.Equal(t, testCase.expectedError, routeAdvertisementsBuilder.errorMsg)

		if testCase.expectedError == "" {
			assert.Equal(t, testCase.networkSelectionType,
				routeAdvertisementsBuilder.Definition.Spec.NetworkSelectors[0].NetworkSelectionType)
			if testCase.clusterUserDefinedNetworkSelector != nil {
				assert.Equal(t, testCase.clusterUserDefinedNetworkSelector,
					routeAdvertisementsBuilder.Definition.Spec.NetworkSelectors[0].ClusterUserDefinedNetworkSelector)
			}
		}
	}
}

func buildValidRouteAdvertisementsBuilder(apiClient *clients.Settings) *RouteAdvertisementsBuilder {
	return NewRouteAdvertisementsBuilder(apiClient, defaultRouteAdvertisementsName, defaultRouteAdvertisementsNsName)
}

func buildInValidRouteAdvertisementsBuilder(apiClient *clients.Settings) *RouteAdvertisementsBuilder {
	return NewRouteAdvertisementsBuilder(apiClient, "", defaultRouteAdvertisementsNsName)
}

func buildRouteAdvertisementsTestClientWithDummyObject() *clients.Settings {
	return clients.GetTestClients(clients.TestClientParams{
		K8sMockObjects: []runtime.Object{
			buildDummyRouteAdvertisements(defaultRouteAdvertisementsName),
		},
		SchemeAttachers: routeAdvertisementsTestSchemes,
	})
}

// buildDummyRouteAdvertisements returns a RouteAdvertisements with the provided name.
func buildDummyRouteAdvertisements(name string) *routeadvertisementstypes.RouteAdvertisements {
	return &routeadvertisementstypes.RouteAdvertisements{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: defaultRouteAdvertisementsNsName,
		},
		Spec: routeadvertisementstypes.RouteAdvertisementsSpec{
			Advertisements: defaultRouteAdvertisementsAdverts,
		},
	}
}
