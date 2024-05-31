// api/v1/k8sdata.go
package v1

type K8sData struct {
	Nodes                  []string `json:"nodes"`
	Namespaces             []string `json:"namespaces"`
	Pods                   []string `json:"pods"`
	Deployments            []string `json:"deployments"`
	StatefulSets           []string `json:"statefulSets"`
	Services               []string `json:"services"`
	PersistentVolumes      []string `json:"persistentVolumes"`
	PersistentVolumeClaims []string `json:"persistentVolumeClaims"`
	StorageClasses         []string `json:"storageClasses"`
	VolumeSnapshotClasses  []string `json:"volumeSnapshotClasses"`
	// Add other Kubernetes objects here
}
