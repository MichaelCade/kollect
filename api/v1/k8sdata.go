package kollect

type NodeInfo struct {
	Name    string
	Roles   string
	Age     string
	Version string
	OSImage string
}

type PodsInfo struct {
	Name      string
	Namespace string
	Status    string
}

type DeploymentInfo struct {
	Name       string
	Namespace  string
	Containers []string
	Images     []string
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
	Pods                   []PodsInfo
	Deployments            []DeploymentInfo
	StatefulSets           []StatefulSetInfo
	Services               []string
	PersistentVolumes      []string
	PersistentVolumeClaims []string
	StorageClasses         []string
	VolumeSnapshotClasses  []string
	// Add other fields as needed
}
