// pkg/kollect/kollect_test.go
package kollect

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestCollectData(t *testing.T) {
	// Create a fake clientset
	clientset := fake.NewSimpleClientset()

	// Call the CollectData function
	data, err := CollectData(clientset)
	if err != nil {
		t.Fatalf("CollectData failed: %v", err)
	}

	// Check that the data contains the correct number of objects
	// Since the clientset is fake and doesn't contain any objects, all counts should be 0
	if len(data.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(data.Nodes))
	}
	if len(data.Namespaces) != 0 {
		t.Errorf("Expected 0 namespaces, got %d", len(data.Namespaces))
	}
	if len(data.Pods) != 0 {
		t.Errorf("Expected 0 pods, got %d", len(data.Pods))
	}
	if len(data.Deployments) != 0 {
		t.Errorf("Expected 0 deployments, got %d", len(data.Deployments))
	}
	if len(data.StatefulSets) != 0 {
		t.Errorf("Expected 0 statefulSets, got %d", len(data.StatefulSets))
	}
	if len(data.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(data.Services))
	}
	if len(data.PersistentVolumes) != 0 {
		t.Errorf("Expected 0 persistentVolumes, got %d", len(data.PersistentVolumes))
	}
	if len(data.PersistentVolumeClaims) != 0 {
		t.Errorf("Expected 0 persistentVolumeClaims, got %d", len(data.PersistentVolumeClaims))
	}
	if len(data.StorageClasses) != 0 {
		t.Errorf("Expected 0 storageClasses, got %d", len(data.StorageClasses))
	}
	if len(data.VolumeSnapshotClasses) != 0 {
		t.Errorf("Expected 0 volumeSnapshotClasses, got %d", len(data.VolumeSnapshotClasses))
	}
}
