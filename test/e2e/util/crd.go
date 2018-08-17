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

package util

// crd contains functions that construct boilerplate CRD definitions.

import (
	"math/rand"
	"os"
	"sync"
	"time"

	channelsv1alpha1 "github.com/knative/eventing/pkg/apis/channels/v1alpha1"
	feedsv1alpha1 "github.com/knative/eventing/pkg/apis/feeds/v1alpha1"
	flowsv1alpha1 "github.com/knative/eventing/pkg/apis/flows/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/kr/pretty"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ResourceNames holds names of related Config, Route, Revision, Flow objects.
type ResourceNames struct {
	Config      string
	Route       string
	Revision    string
	Flow        string
	ClusterBus  string
	EventType   string
	EventSource string
}

// Route returns a Route object in namespace using the route and configuration
// names in names.
func Route(namespace string, names ResourceNames) *v1alpha1.Route {
	return &v1alpha1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      names.Route,
		},
		Spec: v1alpha1.RouteSpec{
			Traffic: []v1alpha1.TrafficTarget{
				v1alpha1.TrafficTarget{
					Name:              names.Route,
					ConfigurationName: names.Config,
					Percent:           100,
				},
			},
		},
	}
}

// Configuration returns a Configuration object in namespace with the name names.Config
// that uses the image specifed by imagePath.
func Configuration(namespace string, names ResourceNames, imagePath string) *v1alpha1.Configuration {
	return &v1alpha1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      names.Config,
		},
		Spec: v1alpha1.ConfigurationSpec{
			RevisionTemplate: v1alpha1.RevisionTemplateSpec{
				Spec: v1alpha1.RevisionSpec{
					Container: corev1.Container{
						Image: imagePath,
					},
				},
			},
		},
	}
}

// FlowFromFile return a Flow object that is created by reading a yaml file
func FlowFromFile(yamlFile string) (flow *flowsv1alpha1.Flow) {
	flow = &flowsv1alpha1.Flow{}
	FromFile(yamlFile, flow)
	return
}

// ClusterBusFromFile return a Flow object that is created by reading a yaml file
func ClusterBusFromFile(yamlFile string) (bus *channelsv1alpha1.ClusterBus) {
	bus = &channelsv1alpha1.ClusterBus{}
	FromFile(yamlFile, bus)
	return
}

// EventTypeFromFile return a EventType object that is created by reading a yaml file
func EventTypeFromFile(yamlFile string) (eType *feedsv1alpha1.EventType) {
	eType = &feedsv1alpha1.EventType{}
	FromFile(yamlFile, eType)
	return
}

// EventSourceFromFile return a EventSource object that is created by reading a yaml file
func EventSourceFromFile(yamlFile string) (eSource *feedsv1alpha1.EventSource) {
	eSource = &feedsv1alpha1.EventSource{}
	FromFile(yamlFile, eSource)
	return
}

// FromFile pupulates a struct by reading a yaml file
func FromFile(yamlFile string, into interface{}) {
	file, err := os.Open(yamlFile)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	dec := yaml.NewYAMLToJSONDecoder(file)

	if err := dec.Decode(into); err != nil {
		panic(err.Error())
	}

	pretty.Println(into)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz"
	randSuffixLen = 8
)

// r is used by AppendRandomString to generate a random string. It is seeded with the time
// at import so the strings will be different between test runs.
var r *rand.Rand

// once is used to initialize r
var once sync.Once

func initSeed(logger *zap.SugaredLogger) func() {
	return func() {
		seed := time.Now().UTC().UnixNano()
		logger.Infof("Seeding rand.Rand with %v", seed)
		r = rand.New(rand.NewSource(seed))
	}
}

// AppendRandomString will generate a random string that begins with prefix. This is useful
// if you want to make sure that your tests can run at the same time against the same
// environment without conflicting. This method will seed rand with the current time when
// called for the first time.
func AppendRandomString(prefix string, logger *zap.SugaredLogger) string {
	once.Do(initSeed(logger))
	suffix := make([]byte, randSuffixLen)
	for i := range suffix {
		suffix[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return prefix + string(suffix)
}
