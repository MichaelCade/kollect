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

type ServiceInfo struct {
	Name      string
	Namespace string
	Type      string
	ClusterIP string
	Ports     string
}

type PersistentVolumeInfo struct {
	Name            string
	Capacity        string
	AccessModes     string
	Status          string
	AssociatedClaim string
	StorageClass    string
	VolumeMode      string
}

type PersistentVolumeClaimInfo struct {
	Name         string
	Namespace    string
	Status       string
	Volume       string
	Capacity     string
	AccessMode   string
	StorageClass string
}

type StorageClassInfo struct {
	Name            string
	Provisioner     string
	VolumeExpansion string
}

type VolumeSnapshotClassInfo struct {
	Name   string
	Driver string
}

type VolumeSnapshotInfo struct {
	Name              string
	Namespace         string
	Volume            string
	CreationTimestamp string
	RestoreSize       string
	Status            bool
}

type K8sData struct {
	Nodes                  []NodeInfo
	Namespaces             []string
	Pods                   []PodsInfo
	Deployments            []DeploymentInfo
	StatefulSets           []StatefulSetInfo
	Services               []ServiceInfo
	PersistentVolumes      []PersistentVolumeInfo
	PersistentVolumeClaims []PersistentVolumeClaimInfo
	StorageClasses         []StorageClassInfo
	VolumeSnapshotClasses  []VolumeSnapshotClassInfo
	VolumeSnapshots        []VolumeSnapshotInfo
	// Add other fields as needed
}
