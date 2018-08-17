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
	"testing"

	"github.com/projectodd/kwsk/test/e2e/util"
)

func TestEvents(t *testing.T) {
	clients := Setup(t)

	//add test case specific name to its own logger
	logger := util.Logger.Named("TestEvent")

	names := &util.ResourceNames{}
	util.CleanupOnInterrupt(func() { TearDown(clients, *names) }, logger)
	defer TearDown(clients, *names)

	logger.Infof("Creating a new Flow")

	client, err := NewWhiskClient()
	if err != nil {
		t.Fatalf("Failed to initialize client %v", err)
	}

	//TODO: client.Insert(action *Action, overwrite bool)

	err := CreateFlow(clients, logger, names, "flow.yaml")
	if err != nil {
		t.Fatalf("Failed to create Flow: %v", err)
	}
	logger.Infof("Created Flow named %s", names.Flow)

	//TODO: Check that k8s events are routed to the action
}
