package e2e

import (
	"testing"

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
		clients.Delete([]string{names.Route}, []string{names.Config}, []string{names.Flow})
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
