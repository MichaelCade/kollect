package kollect

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	k8sdata "github.com/michaelcade/kollect/api/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	data.VirtualMachines, err = fetchVirtualMachines(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch VirtualMachines: %v", err)
		data.VirtualMachines = []k8sdata.VirtualMachineInfo{}
	}
	data.DataVolumes, err = fetchDataVolumes(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch DataVolumes: %v", err)
		data.DataVolumes = []k8sdata.DataVolumeInfo{}
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
	data.CustomResourceDefs, err = fetchCustomResourceDefinitions(ctx, config)
	if err != nil {
		log.Printf("Warning: Failed to fetch CRDs: %v", err)
		data.CustomResourceDefs = []k8sdata.CRDInfo{}
	}

	data.VirtualMachines, err = fetchVirtualMachines(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch VirtualMachines: %v", err)
		data.VirtualMachines = []k8sdata.VirtualMachineInfo{}
	}

	data.DataVolumes, err = fetchDataVolumes(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch DataVolumes: %v", err)
		data.DataVolumes = []k8sdata.DataVolumeInfo{}
	}

	return data, nil
}

func CollectDataWithContext(ctx context.Context, kubeconfig string, contextName string) (k8sdata.K8sData, error) {
	var data k8sdata.K8sData

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	)
	config, err := clientConfig.ClientConfig()
	if err != nil {
		return k8sdata.K8sData{}, fmt.Errorf("error building kubeconfig with context %s: %v", contextName, err)
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

	data.CustomResourceDefs, err = fetchCustomResourceDefinitions(ctx, config)
	if err != nil {
		log.Printf("Warning: Failed to fetch CRDs: %v", err)
		data.CustomResourceDefs = []k8sdata.CRDInfo{}
	}

	data.VirtualMachines, err = fetchVirtualMachines(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch VirtualMachines: %v", err)
		data.VirtualMachines = []k8sdata.VirtualMachineInfo{}
	}

	data.DataVolumes, err = fetchDataVolumes(ctx, dynamicClient)
	if err != nil {
		log.Printf("Warning: Failed to fetch DataVolumes: %v", err)
		data.DataVolumes = []k8sdata.DataVolumeInfo{}
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
		age := formatDuration(time.Since(node.CreationTimestamp.Time))
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

func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	return fmt.Sprintf("%dd%dh%dm", days, hours, minutes)
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

func fetchCustomResourceDefinitions(ctx context.Context, config *rest.Config) ([]k8sdata.CRDInfo, error) {
	apiextensionsClientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create apiextensions clientset: %v", err)
	}

	crds, err := apiextensionsClientset.ApiextensionsV1().CustomResourceDefinitions().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CRDs: %v", err)
	}

	var crdInfos []k8sdata.CRDInfo
	for _, crd := range crds.Items {
		age := formatDuration(time.Since(crd.CreationTimestamp.Time))

		for _, version := range crd.Spec.Versions {
			if version.Served {
				scope := string(crd.Spec.Scope)
				crdInfos = append(crdInfos, k8sdata.CRDInfo{
					Name:    crd.Name,
					Group:   crd.Spec.Group,
					Version: version.Name,
					Kind:    crd.Spec.Names.Kind,
					Scope:   scope,
					Age:     age,
				})
			}
		}
	}

	return crdInfos, nil
}

func fetchVirtualMachines(ctx context.Context, dynamicClient dynamic.Interface) ([]k8sdata.VirtualMachineInfo, error) {
	gvr := schema.GroupVersionResource{
		Group:    "kubevirt.io",
		Version:  "v1",
		Resource: "virtualmachines",
	}

	vms, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "the server could not find the requested resource") {
			log.Printf("Warning: VirtualMachines resource not found in the cluster. Is KubeVirt installed?")
			return []k8sdata.VirtualMachineInfo{}, nil
		}
		return nil, err
	}

	var vmInfos []k8sdata.VirtualMachineInfo
	for _, vm := range vms.Items {
		vmInfo := k8sdata.VirtualMachineInfo{
			Name:        vm.GetName(),
			Namespace:   vm.GetNamespace(),
			DataVolumes: []string{},
		}

		status, found, _ := unstructured.NestedString(vm.Object, "status", "printableStatus")
		if !found {
			phase, phaseFound, _ := unstructured.NestedString(vm.Object, "status", "phase")
			if phaseFound {
				status = phase
			} else {
				status = "Unknown"
			}
		}
		vmInfo.Status = status

		conditions, found, _ := unstructured.NestedSlice(vm.Object, "status", "conditions")
		if found {
			for _, c := range conditions {
				condition, ok := c.(map[string]interface{})
				if !ok {
					continue
				}

				condType, typeFound, _ := unstructured.NestedString(condition, "type")
				status, statusFound, _ := unstructured.NestedString(condition, "status")

				if typeFound && statusFound && condType == "Ready" {
					vmInfo.Ready = status == "True"
					break
				}
			}
		}

		creationTimestamp := vm.GetCreationTimestamp()
		vmInfo.Age = formatDuration(time.Since(creationTimestamp.Time))

		runStrategy, found, _ := unstructured.NestedString(vm.Object, "spec", "runStrategy")
		if !found {
			running, runningFound, _ := unstructured.NestedBool(vm.Object, "spec", "running")
			if runningFound {
				if running {
					runStrategy = "Always"
				} else {
					runStrategy = "Manual"
				}
			} else {
				runStrategy = "Unknown"
			}
		}
		vmInfo.RunStrategy = runStrategy

		getCPU := func() string {
			paths := [][]string{
				{"spec", "template", "spec", "domain", "resources", "requests", "cpu"},
				{"spec", "template", "spec", "domain", "cpu", "cores"},
				{"spec", "domain", "resources", "requests", "cpu"},
				{"spec", "domain", "cpu", "cores"},
			}

			for _, path := range paths {
				cpu, found, _ := unstructured.NestedString(vm.Object, path...)
				if found && cpu != "" {
					return cpu
				}

				cpuFloat, found, _ := unstructured.NestedFloat64(vm.Object, path...)
				if found {
					return fmt.Sprintf("%g", cpuFloat)
				}

				cpuInt, found, _ := unstructured.NestedInt64(vm.Object, path...)
				if found {
					return fmt.Sprintf("%d", cpuInt)
				}
			}

			return "N/A"
		}

		vmInfo.CPU = getCPU()

		getMemory := func() string {
			paths := [][]string{
				{"spec", "template", "spec", "domain", "resources", "requests", "memory"},
				{"spec", "template", "spec", "domain", "memory", "guest"},
				{"spec", "domain", "resources", "requests", "memory"},
				{"spec", "domain", "memory", "guest"},
			}

			for _, path := range paths {
				memory, found, _ := unstructured.NestedString(vm.Object, path...)
				if found && memory != "" {
					return memory
				}
			}

			return "N/A"
		}

		vmInfo.Memory = getMemory()

		extractStorage := func() ([]string, []string) {
			var dataVols []string
			var storageVolumes []string

			volumePaths := [][]string{
				{"spec", "template", "spec", "volumes"},
				{"spec", "volumes"},
			}

			for _, volumePath := range volumePaths {
				volumes, found, _ := unstructured.NestedSlice(vm.Object, volumePath...)
				if !found || len(volumes) == 0 {
					continue
				}

				for _, v := range volumes {
					volume, ok := v.(map[string]interface{})
					if !ok {
						continue
					}

					volumeName, volNameFound, _ := unstructured.NestedString(volume, "name")
					if !volNameFound || volumeName == "" {
						continue
					}

					dataVolume, dvFound, _ := unstructured.NestedMap(volume, "dataVolume")
					if dvFound {
						name, nameFound, _ := unstructured.NestedString(dataVolume, "name")
						if nameFound && name != "" {
							dataVols = append(dataVols, name)
							storageVolumes = append(storageVolumes, fmt.Sprintf("DataVolume: %s", name))
						}
						continue
					}

					pvc, pvcFound, _ := unstructured.NestedMap(volume, "persistentVolumeClaim")
					if pvcFound {
						claimName, claimFound, _ := unstructured.NestedString(pvc, "claimName")
						if claimFound && claimName != "" {
							storageVolumes = append(storageVolumes, fmt.Sprintf("PVC: %s", claimName))

							if strings.Contains(claimName, vm.GetName()) || strings.HasPrefix(claimName, "dv-") {
								dataVols = append(dataVols, claimName)
							}
						}
						continue
					}

					if _, found, _ := unstructured.NestedMap(volume, "configMap"); found {
						configMapName, nameFound, _ := unstructured.NestedString(volume, "configMap", "name")
						if nameFound {
							storageVolumes = append(storageVolumes, fmt.Sprintf("ConfigMap: %s", configMapName))
						}
						continue
					}

					if _, found, _ := unstructured.NestedMap(volume, "secret"); found {
						secretName, nameFound, _ := unstructured.NestedString(volume, "secret", "secretName")
						if nameFound {
							storageVolumes = append(storageVolumes, fmt.Sprintf("Secret: %s", secretName))
						}
						continue
					}

					if _, found, _ := unstructured.NestedMap(volume, "emptyDir"); found {
						storageVolumes = append(storageVolumes, fmt.Sprintf("EmptyDir: %s", volumeName))
						continue
					}

					if _, found, _ := unstructured.NestedMap(volume, "hostPath"); found {
						path, pathFound, _ := unstructured.NestedString(volume, "hostPath", "path")
						if pathFound {
							storageVolumes = append(storageVolumes, fmt.Sprintf("HostPath: %s -> %s", volumeName, path))
						} else {
							storageVolumes = append(storageVolumes, fmt.Sprintf("HostPath: %s", volumeName))
						}
						continue
					}
					storageVolumes = append(storageVolumes, fmt.Sprintf("Volume: %s", volumeName))
				}
			}

			return dataVols, storageVolumes
		}

		dataVols, storageVols := extractStorage()
		vmInfo.DataVolumes = dataVols
		vmInfo.Storage = storageVols

		vmInfos = append(vmInfos, vmInfo)
	}

	return vmInfos, nil
}

func fetchDataVolumes(ctx context.Context, dynamicClient dynamic.Interface) ([]k8sdata.DataVolumeInfo, error) {
	gvr := schema.GroupVersionResource{
		Group:    "cdi.kubevirt.io",
		Version:  "v1beta1",
		Resource: "datavolumes",
	}

	dvs, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "the server could not find the requested resource") {
			log.Printf("Warning: DataVolumes resource not found in the cluster. Is CDI installed?")
			return []k8sdata.DataVolumeInfo{}, nil
		}
		return nil, err
	}

	var dvInfos []k8sdata.DataVolumeInfo
	for _, dv := range dvs.Items {
		dvInfo := k8sdata.DataVolumeInfo{
			Name:      dv.GetName(),
			Namespace: dv.GetNamespace(),
		}

		phase, found, _ := unstructured.NestedString(dv.Object, "status", "phase")
		if found {
			dvInfo.Phase = phase
		}

		capacity, found, _ := unstructured.NestedString(dv.Object, "spec", "pvc", "resources", "requests", "storage")
		if found {
			dvInfo.Size = capacity
		}

		source, found, _ := unstructured.NestedMap(dv.Object, "spec", "source")
		if found {
			for sourceType, sourceData := range source {
				dvInfo.SourceType = sourceType

				switch sourceType {
				case "http":
					sourceMap, ok := sourceData.(map[string]interface{})
					if ok {
						url, urlFound, _ := unstructured.NestedString(sourceMap, "url")
						if urlFound {
							dvInfo.SourceInfo = url
						}
					}
				case "pvc":
					sourceMap, ok := sourceData.(map[string]interface{})
					if ok {
						name, nameFound, _ := unstructured.NestedString(sourceMap, "name")
						namespace, nsFound, _ := unstructured.NestedString(sourceMap, "namespace")
						if nameFound && nsFound {
							dvInfo.SourceInfo = fmt.Sprintf("%s/%s", namespace, name)
						} else if nameFound {
							dvInfo.SourceInfo = name
						}
					}
				case "registry":
					sourceMap, ok := sourceData.(map[string]interface{})
					if ok {
						url, urlFound, _ := unstructured.NestedString(sourceMap, "url")
						if urlFound {
							dvInfo.SourceInfo = url
						}
					}
				}
				break
			}
		}

		creationTimestamp := dv.GetCreationTimestamp()
		dvInfo.Age = formatDuration(time.Since(creationTimestamp.Time))

		dvInfos = append(dvInfos, dvInfo)
	}

	return dvInfos, nil
}
