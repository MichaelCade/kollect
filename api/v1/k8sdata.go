package kollect

type NodeInfo struct {
	Name    string
	Roles   string
	Age     string
	Version string
	OSImage string
}

type K8sData struct {
	Nodes                  []NodeInfo
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
