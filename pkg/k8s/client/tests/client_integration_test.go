//go:build integration
// +build integration

package k8s_client_test

import (
	"bytes"
	"context"
	"flag"
	k8s_client "go_testing/pkg/k8s/client"
	"log"
	"regexp"
	"testing"
)

const Namespace = "client"

var client k8s_client.IK8sClient
var logger *log.Logger
var buffer bytes.Buffer

var external = flag.Bool(k8s_client.External, false, "Name of location to greet")
var tests = map[string]func(*testing.T){ // tests to run for all types of clients
	"getPodCount": testGetPodCount,
	"logPodNames": testLogPodNames,
}

func TestInternalClient(t *testing.T) {
	if *external {
		t.Skip()
	}
	defer commonTeardown()()
	commonSetup()

	client = k8s_client.NewK8sClient(logger, k8s_client.InCluster, nil)

	for k, v := range tests {
		t.Run(k, v)
	}
}

func TestExternalClient(t *testing.T) {
	if !(*external) {
		t.Skip()
	}
	defer commonTeardown()()
	commonSetup()

	client = k8s_client.NewK8sClient(logger, k8s_client.OutCluster,
		&k8s_client.K8sClientConfig{
			OutConfig: &k8s_client.OutClusterClientConfig{
				KubeconfigFile: k8s_client.DefaultKubeConfigPath(),
			},
		})

	for k, v := range tests {
		t.Run(k, v)
	}
}

func commonSetup() {
	logger = log.New(&buffer, "", 0)
}

func commonTeardown() func() {
	return func() {
		client = nil
		logger = nil
		buffer.Reset()
	}
}

func testGetPodCount(t *testing.T) {
	got := 1
	want, err := client.GetPodCount(context.TODO(), Namespace)
	if got != want || err != nil {
		t.Errorf("Did not get the exptected amount of pods, or err: %s", err)
	}
}

func testLogPodNames(t *testing.T) {
	client.LogPodNames(context.TODO(), Namespace)
	logOutput := buffer.String()
	wantRegExp := "Found"
	failRegExp := "err|error|Error"

	if matched, _ := regexp.MatchString(failRegExp, logOutput); matched {
		t.Errorf("Logged error: %s", logOutput)
	} else if matched, _ := regexp.MatchString(wantRegExp, logOutput); !matched {
		t.Errorf("Did not match wanted output: %s", logOutput)
	}
}
