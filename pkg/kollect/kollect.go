// pkg/kollect/kollect.go
package kollect

import (
	"context"
	"strings"
	"time"

	k8sdata "github.com/michaelcade/kollect/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CollectStorageData(kubeconfig string) (k8sdata.K8sData, error) {
	var data k8sdata.K8sData
	var err error

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

	data.VolumeSnapshotClasses, err = fetchVolumeSnapshotClasses(dynamicClient)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Nodes, err = fetchNodes(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Namespaces, err = fetchNamespaces(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.StatefulSets, err = fetchStatefulSets(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.PersistentVolumes, err = fetchPersistentVolumes(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.PersistentVolumeClaims, err = fetchPersistentVolumeClaims(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.StorageClasses, err = fetchStorageClasses(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	return data, nil
}

func CollectData(kubeconfig string) (k8sdata.K8sData, error) {
	var data k8sdata.K8sData
	var err error

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

	data.VolumeSnapshotClasses, err = fetchVolumeSnapshotClasses(dynamicClient)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Nodes, err = fetchNodes(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Namespaces, err = fetchNamespaces(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Pods, err = fetchPods(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Deployments, err = fetchDeployments(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.StatefulSets, err = fetchStatefulSets(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.Services, err = fetchServices(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.PersistentVolumes, err = fetchPersistentVolumes(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.PersistentVolumeClaims, err = fetchPersistentVolumeClaims(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	data.StorageClasses, err = fetchStorageClasses(clientset)
	if err != nil {
		return k8sdata.K8sData{}, err
	}

	return data, nil
}

func fetchNodes(clientset *kubernetes.Clientset) ([]k8sdata.NodeInfo, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
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

func fetchNamespaces(clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceNames []string
	for _, namespace := range namespaces.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}

	return namespaceNames, nil
}

func fetchPods(clientset *kubernetes.Clientset) ([]k8sdata.PodsInfo, error) {
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
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

func fetchDeployments(clientset *kubernetes.Clientset) ([]string, error) {
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentNames []string
	for _, deployment := range deployments.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}

	return deploymentNames, nil
}

func fetchStatefulSets(clientset *kubernetes.Clientset) ([]k8sdata.StatefulSetInfo, error) {
	statefulSets, err := clientset.AppsV1().StatefulSets("").List(context.Background(), v1.ListOptions{})
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

func fetchServices(clientset *kubernetes.Clientset) ([]string, error) {
	services, err := clientset.CoreV1().Services("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var serviceNames []string
	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}

	return serviceNames, nil
}

func fetchPersistentVolumes(clientset *kubernetes.Clientset) ([]string, error) {
	persistentVolumes, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var persistentVolumeNames []string
	for _, persistentVolume := range persistentVolumes.Items {
		persistentVolumeNames = append(persistentVolumeNames, persistentVolume.Name)
	}

	return persistentVolumeNames, nil
}

func fetchPersistentVolumeClaims(clientset *kubernetes.Clientset) ([]string, error) {
	persistentVolumeClaims, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var persistentVolumeClaimNames []string
	for _, persistentVolumeClaim := range persistentVolumeClaims.Items {
		persistentVolumeClaimNames = append(persistentVolumeClaimNames, persistentVolumeClaim.Name)
	}

	return persistentVolumeClaimNames, nil
}

func fetchStorageClasses(clientset *kubernetes.Clientset) ([]string, error) {
	storageClasses, err := clientset.StorageV1().StorageClasses().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var storageClassNames []string
	for _, storageClass := range storageClasses.Items {
		storageClassNames = append(storageClassNames, storageClass.Name)
	}

	return storageClassNames, nil
}

func fetchVolumeSnapshotClasses(client dynamic.Interface) ([]string, error) {
	gvr := schema.GroupVersionResource{Group: "snapshot.storage.k8s.io", Version: "v1", Resource: "volumesnapshotclasses"}
	list, err := client.Resource(gvr).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range list.Items {
		names = append(names, item.GetName())
	}

	return names, nil
}
