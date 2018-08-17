package e2e

import (
	"net/http"
	"testing"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/projectodd/kwsk/test/e2e/util"
	"go.uber.org/zap"

	// Mysteriously required to support GCP auth (required by k8s libs).
	// Apparently just importing it is enough. @_@ side effects @_@.
	// https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	// NamespaceName is the namespace used for the e2e tests.
	NamespaceName = "noodleburg"

	configName = "prod"
	routeName  = "noodleburg"
	//flowName = "k8s-event-flow"
)

// Setup creates the client objects needed in the e2e tests.
func Setup(t *testing.T) *util.Clients {
	clients, err := util.NewClients(
		util.Flags.Kubeconfig,
		util.Flags.Cluster,
		NamespaceName)
	if err != nil {
		t.Fatalf("Couldn't initialize clients: %v", err)
	}
	return clients
}

// TearDown will delete created names using clients.
func TearDown(clients *util.Clients, names util.ResourceNames) {
	if clients != nil {
		clients.Delete([]string{names.Route}, []string{names.Config}, []string{names.Flow}, []string{names.ClusterBus}, []string{names.EventType}, []string{names.EventSource})
	}
}

// CreateRouteAndConfig will create Route and Config objects using clients.
// The Config object will serve requests to a container started from the image at imagePath.
func CreateRouteAndConfig(clients *util.Clients, logger *zap.SugaredLogger, names *util.ResourceNames, imagePath string) error {
	names.Config = util.AppendRandomString(configName, logger)
	names.Route = util.AppendRandomString(routeName, logger)

	_, err := clients.Configs.Create(
		util.Configuration(NamespaceName, *names, imagePath))
	if err != nil {
		return err
	}
	_, err = clients.Routes.Create(
		util.Route(NamespaceName, *names))
	return err
}

// CreateFlow will create a Flow object using clients. This flow is read from a file named yamlFile.
func CreateFlow(clients *util.Clients, logger *zap.SugaredLogger, names *util.ResourceNames, yamlFile string) error {
	flow, err := clients.Flows.Create(util.FlowFromFile(yamlFile))
	if err != nil {
		return err
	}
	names.Flow = flow.ObjectMeta.Name
	return err
}

// CreateClusterBus will create a ClusterBus object using clients. This clusterBus is read from a file named yamlFile.
func CreateClusterBus(clients *util.Clients, logger *zap.SugaredLogger, names *util.ResourceNames, yamlFile string) error {
	bus, err := clients.ClusterBuses.Create(util.ClusterBusFromFile(yamlFile))
	if err != nil {
		return err
	}
	names.ClusterBus = bus.ObjectMeta.Name
	return err
}

// CreateEventType will create an EventType object using clients. This eventType is read from a file named yamlFile.
func CreateEventType(clients *util.Clients, logger *zap.SugaredLogger, names *util.ResourceNames, yamlFile string) error {
	eType, err := clients.EventTypes.Create(util.EventTypeFromFile(yamlFile))
	if err != nil {
		return err
	}
	names.EventType = eType.ObjectMeta.Name
	return err
}

// CreateEventSource will create an EventSource object using clients. This eventSource is read from a file named yamlFile.
func CreateEventSource(clients *util.Clients, logger *zap.SugaredLogger, names *util.ResourceNames, yamlFile string) error {
	eSource, err := clients.EventSources.Create(util.EventSourceFromFile(yamlFile))
	if err != nil {
		return err
	}
	names.EventSource = eSource.ObjectMeta.Name
	return err
}

// NewClient creates new OpenWhisk client based on configuration from whick.properties that is found under OPENWHISK_HOME
// env variable
func NewWhiskClient() (*whisk.Client, error) {
	config, err := whisk.GetWhiskPropertiesConfig()
	if err != nil {
		return nil, err
	}
	client, err := whisk.NewClient(http.DefaultClient, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
