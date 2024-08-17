package kollect

type NodeInfo struct {
	Name    string
	Roles   string
	Age     string
	Version string
	OSImage string
}

type StatefulSetInfo struct {
	Name          string
	Namespace     string
	ReadyReplicas int32
	Image         string
}

type K8sData struct {
	Nodes                  []NodeInfo
	Namespaces             []string
	Pods                   []string
	Deployments            []string
	StatefulSets           []StatefulSetInfo
	Services               []string
	PersistentVolumes      []string
	PersistentVolumeClaims []string
	StorageClasses         []string
	VolumeSnapshotClasses  []string
	// Add other fields as needed
}
