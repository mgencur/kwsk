// +build e2e

/*
Copyright 2018 The Knative Authors
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

package e2e

import (
	"net/http"
	"testing"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

func TestEvents(t *testing.T) {
	// clients := Setup(t)

	//add test case specific name to its own logger
	//logger := util.Logger.Named("TestBasicAction")

	// names := &util.ResourceNames{}

	// logger.Infof("Creating a new Flow")

	//err := CreateFlow(clients, logger, names, "flow.yaml")
	// if err != nil {
	// 	t.Fatalf("Failed to create Flow: %v", err)
	// }
	// logger.Infof("Created Flow named %s", names.Flow)
	//TODO: Wait for Flow to be ready (similar to waiting for Route/Config in crd_checks.go)

	// logger.Infof("Creating a new ClusterBus")

	// err := CreateClusterBus(clients, logger, names, "clusterbus-stub.yaml")
	// if err != nil {
	// 	t.Fatalf("Failed to create ClusterBus: %v", err)
	// }
	// logger.Infof("Created ClusterBus named %s", names.ClusterBus)

	// logger.Infof("Creating a new EventType")

	// err := CreateEventType(clients, logger, names, "k8s-eventtype.yaml")
	// if err != nil {
	// 	t.Fatalf("Failed to create EventType: %v", err)
	// }
	// logger.Infof("Created EventType named %s", names.EventType)

	// logger.Infof("Creating a new EventSource")

	// err := CreateEventSource(clients, logger, names, "k8s-eventsource.yaml")
	// if err != nil {
	// 	t.Fatalf("Failed to create EventSource: %v", err)
	// }
	// logger.Infof("Created EventSource named %s", names.EventSource)

	config, err := whisk.GetWhiskPropertiesConfig()
	if err != nil {
		t.Fatalf("Failed to read properties from file %v", err)
	}
	client, err := whisk.NewClient(http.DefaultClient, config)
	if err != nil {
		t.Fatalf("Failed to initialize client %v", err)
	}
	options := &whisk.ActionListOptions{
		Limit: 30,
		Skip:  30,
	}
	_, _, err = client.Actions.List("hello", options)
	if err != nil {
		t.Fatalf("Failed to call list %v", err)
	}
}
