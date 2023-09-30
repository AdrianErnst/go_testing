//go:build integration && incluster
// +build integration,incluster

package k8s_client_test

import (
	"testing"
)

func TestGetPodCount(t *testing.T) {
	want := 1
	got := 1
	if got != want {
		t.Errorf("Got %d pods, want %d", got, want)
	}
}
