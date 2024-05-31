package kollect

type K8sData struct {
	Nodes                  []string
	Namespaces             []string
	Pods                   []string
	Deployments            []string
	StatefulSets           []string
	Services               []string
	PersistentVolumes      []string
	PersistentVolumeClaims []string
	StorageClasses         []string
	VolumeSnapshotClasses  []string
	// Add other fields as needed
}
