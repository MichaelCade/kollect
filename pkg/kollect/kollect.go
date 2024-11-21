package kollect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log"

	k8sdata "github.com/michaelcade/kollect/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CollectStorageData(ctx context.Context, kubeconfig string) (k8sdata.K8sData, error) {
	var data k8sdata.K8sData
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	data.PersistentVolumes, err = fetchPersistentVolumes(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching PersistentVolumes: %v", err)
	}
	data.PersistentVolumeClaims, err = fetchPersistentVolumeClaims(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching PersistentVolumeClaims: %v", err)
	}
	data.StorageClasses, err = fetchStorageClasses(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching StorageClasses: %v", err)
	}
	data.VolumeSnapshotClasses, err = fetchVolumeSnapshotClasses(ctx, dynamicClient)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching VolumeSnapshotClasses: %v", err)
	}
	data.VolumeSnapshots, err = fetchVolumeSnapshots(ctx, dynamicClient)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching VolumeSnapshots: %v", err)
	}
	return data, nil
}

func CollectData(ctx context.Context, kubeconfig string) (k8sdata.K8sData, error) {
	var data k8sdata.K8sData
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return k8sdata.K8sData{}, err
	}
	data.Nodes, err = fetchNodes(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching Nodes: %v", err)
	}
	data.Namespaces, err = fetchNamespaces(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching Namespaces: %v", err)
	}
	data.Pods, err = fetchPods(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching Pods: %v", err)
	}
	data.Deployments, err = fetchDeployments(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching Deployments: %v", err)
	}
	data.StatefulSets, err = fetchStatefulSets(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching StatefulSets: %v", err)
	}
	data.Services, err = fetchServices(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching Services: %v", err)
	}
	data.PersistentVolumes, err = fetchPersistentVolumes(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching PersistentVolumes: %v", err)
	}
	data.PersistentVolumeClaims, err = fetchPersistentVolumeClaims(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching PersistentVolumeClaims: %v", err)
	}
	data.StorageClasses, err = fetchStorageClasses(ctx, clientset)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching StorageClasses: %v", err)
	}
	data.VolumeSnapshotClasses, err = fetchVolumeSnapshotClasses(ctx, dynamicClient)
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error fetching VolumeSnapshotClasses: %v", err)
	}
	data.VolumeSnapshots, err = fetchVolumeSnapshots(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: VolumeSnapshots resource not found in the cluster: %v", err)
		data.VolumeSnapshots = []k8sdata.VolumeSnapshotInfo{}
	}
	return data, nil
}

func fetchNodes(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var nodeInfos []k8sdata.NodeInfo
	for _, node := range nodes.Items {
		roles := "none"
		for label := range node.Labels {
			if strings.HasPrefix(label, "node-role.kubernetes.io/") {
				role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
				if roles == "none" {
					roles = role
				} else {
					roles += "," + role
				}
			}
		}
		age := time.Since(node.CreationTimestamp.Time).String()
		version := node.Status.NodeInfo.KubeletVersion
		osImage := node.Status.NodeInfo.OSImage
		nodeInfos = append(nodeInfos, k8sdata.NodeInfo{
			Name:    node.Name,
			Roles:   roles,
			Age:     age,
			Version: version,
			OSImage: osImage,
		})
	}
	return nodeInfos, nil
}

func fetchNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceNames []string
	for _, namespace := range namespaces.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}

	return namespaceNames, nil
}

func fetchPods(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.PodsInfo, error) {
	pods, err := clientset.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podInfos []k8sdata.PodsInfo
	for _, pod := range pods.Items {
		podInfos = append(podInfos, k8sdata.PodsInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
		})
	}

	return podInfos, nil
}

func fetchDeployments(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.DeploymentInfo, error) {
	deployments, err := clientset.AppsV1().Deployments("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentInfos []k8sdata.DeploymentInfo
	for _, deployment := range deployments.Items {
		var containers []string
		var images []string
		for _, container := range deployment.Spec.Template.Spec.Containers {
			containers = append(containers, container.Name)
			images = append(images, container.Image)
		}
		deploymentInfos = append(deploymentInfos, k8sdata.DeploymentInfo{
			Name:       deployment.Name,
			Namespace:  deployment.Namespace,
			Containers: containers,
			Images:     images,
		})
	}

	return deploymentInfos, nil
}

func fetchStatefulSets(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.StatefulSetInfo, error) {
	statefulSets, err := clientset.AppsV1().StatefulSets("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var statefulSetInfos []k8sdata.StatefulSetInfo
	for _, statefulSet := range statefulSets.Items {
		image := ""
		if len(statefulSet.Spec.Template.Spec.Containers) > 0 {
			image = statefulSet.Spec.Template.Spec.Containers[0].Image
		}
		statefulSetInfos = append(statefulSetInfos, k8sdata.StatefulSetInfo{
			Name:          statefulSet.Name,
			Namespace:     statefulSet.Namespace,
			ReadyReplicas: statefulSet.Status.ReadyReplicas,
			Image:         image,
		})
	}

	return statefulSetInfos, nil
}

func fetchServices(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.ServiceInfo, error) {
	services, err := clientset.CoreV1().Services("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceInfos []k8sdata.ServiceInfo
	for _, service := range services.Items {
		ports := []string{}
		for _, port := range service.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
		}
		serviceInfos = append(serviceInfos, k8sdata.ServiceInfo{
			Name:      service.Name,
			Namespace: service.Namespace,
			Type:      string(service.Spec.Type),
			ClusterIP: service.Spec.ClusterIP,
			Ports:     strings.Join(ports, ","),
		})
	}

	return serviceInfos, nil
}

func fetchPersistentVolumes(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.PersistentVolumeInfo, error) {
	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pvInfos []k8sdata.PersistentVolumeInfo
	for _, pv := range pvs.Items {
		accessModes := []string{}
		for _, mode := range pv.Spec.AccessModes {
			accessModes = append(accessModes, string(mode))
		}
		accessModesStr := strings.Join(accessModes, ",")
		pvInfos = append(pvInfos, k8sdata.PersistentVolumeInfo{
			Name:            pv.Name,
			Capacity:        pv.Spec.Capacity.Storage().String(),
			AccessModes:     accessModesStr,
			Status:          string(pv.Status.Phase),
			AssociatedClaim: pv.Spec.ClaimRef.Name,
			StorageClass:    pv.Spec.StorageClassName,
			VolumeMode:      string(*pv.Spec.VolumeMode),
		})
	}

	return pvInfos, nil
}

func fetchPersistentVolumeClaims(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.PersistentVolumeClaimInfo, error) {
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pvcInfos []k8sdata.PersistentVolumeClaimInfo
	for _, pvc := range pvcs.Items {
		storageClassName := ""
		if pvc.Spec.StorageClassName != nil {
			storageClassName = *pvc.Spec.StorageClassName
		}
		pvcInfos = append(pvcInfos, k8sdata.PersistentVolumeClaimInfo{
			Name:         pvc.Name,
			Namespace:    pvc.Namespace,
			Status:       string(pvc.Status.Phase),
			Volume:       pvc.Spec.VolumeName,
			Capacity:     pvc.Spec.Resources.Requests.Storage().String(),
			AccessMode:   string(pvc.Spec.AccessModes[0]),
			StorageClass: storageClassName,
		})
	}

	return pvcInfos, nil
}

func fetchStorageClasses(ctx context.Context, clientset *kubernetes.Clientset) ([]k8sdata.StorageClassInfo, error) {
	storageClasses, err := clientset.StorageV1().StorageClasses().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var storageClassInfos []k8sdata.StorageClassInfo
	for _, sc := range storageClasses.Items {
		allowVolumeExpansion := "false"
		if sc.AllowVolumeExpansion != nil && *sc.AllowVolumeExpansion {
			allowVolumeExpansion = "true"
		}
		storageClassInfos = append(storageClassInfos, k8sdata.StorageClassInfo{
			Name:            sc.Name,
			Provisioner:     sc.Provisioner,
			VolumeExpansion: allowVolumeExpansion,
		})
	}

	return storageClassInfos, nil
}

func fetchVolumeSnapshotClasses(ctx context.Context, dynamicClient dynamic.Interface) ([]k8sdata.VolumeSnapshotClassInfo, error) {
	gvr := schema.GroupVersionResource{
		Group:    "snapshot.storage.k8s.io",
		Version:  "v1",
		Resource: "volumesnapshotclasses",
	}
	volumeSnapshotClasses, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "the server could not find the requested resource") {
			log.Printf("Warning: VolumeSnapshotClasses resource not found in the cluster")
			return nil, nil
		}
		return nil, err
	}

	var volumeSnapshotClassInfos []k8sdata.VolumeSnapshotClassInfo
	for _, vsc := range volumeSnapshotClasses.Items {
		driver, found, err := unstructured.NestedString(vsc.Object, "driver")
		if err != nil || !found {
			return nil, fmt.Errorf("failed to get driver for volume snapshot class %s: %v", vsc.GetName(), err)
		}
		volumeSnapshotClassInfos = append(volumeSnapshotClassInfos, k8sdata.VolumeSnapshotClassInfo{
			Name:   vsc.GetName(),
			Driver: driver,
		})
	}

	return volumeSnapshotClassInfos, nil
}

func fetchVolumeSnapshots(ctx context.Context, dynamicClient dynamic.Interface) ([]k8sdata.VolumeSnapshotInfo, error) {
	gvr := schema.GroupVersionResource{
		Group:    "snapshot.storage.k8s.io",
		Version:  "v1",
		Resource: "volumesnapshots",
	}
	volumeSnapshots, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "the server could not find the requested resource") {
			log.Printf("Warning: VolumeSnapshots resource not found in the cluster")
			return []k8sdata.VolumeSnapshotInfo{}, nil
		}
		return nil, err
	}

	var volumeSnapshotInfos []k8sdata.VolumeSnapshotInfo
	for _, vs := range volumeSnapshots.Items {
		volumeSnapshot := k8sdata.VolumeSnapshotInfo{
			Name:      vs.GetName(),
			Namespace: vs.GetNamespace(),
		}

		if volumeName, found, err := unstructured.NestedString(vs.Object, "spec", "source", "persistentVolumeClaimName"); err == nil && found {
			volumeSnapshot.Volume = volumeName
		}

		if creationTimestamp, found, err := unstructured.NestedString(vs.Object, "metadata", "creationTimestamp"); err == nil && found {
			volumeSnapshot.CreationTimestamp = creationTimestamp
		}

		if restoreSize, found, err := unstructured.NestedString(vs.Object, "status", "restoreSize"); err == nil && found {
			volumeSnapshot.RestoreSize = restoreSize
		}

		if status, found, err := unstructured.NestedBool(vs.Object, "status", "readyToUse"); err == nil && found {
			volumeSnapshot.Status = status
		}

		volumeSnapshotInfos = append(volumeSnapshotInfos, volumeSnapshot)
	}

	return volumeSnapshotInfos, nil
}
