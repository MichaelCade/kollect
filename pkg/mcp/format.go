package mcp

import (
	"fmt"
	"sort"
	"strings"

	k8sdata "github.com/michaelcade/kollect/api/v1"
	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/terraform"
	"github.com/michaelcade/kollect/pkg/vault"
)

// Format functions for all resource types

// Kubernetes resource formatters
func formatK8sNodeInfo(node k8sdata.NodeInfo) string {
	info := fmt.Sprintf("Kubernetes Node: %s\n", node.Name)
	info += fmt.Sprintf("Roles: %s\n", node.Roles)
	info += fmt.Sprintf("Version: %s\n", node.Version)
	info += fmt.Sprintf("OS: %s\n", node.OSImage)
	info += fmt.Sprintf("Age: %s\n", node.Age)
	return info
}

func formatK8sPodInfo(pod k8sdata.PodsInfo) string {
	info := fmt.Sprintf("Kubernetes Pod: %s\n", pod.Name)
	info += fmt.Sprintf("Namespace: %s\n", pod.Namespace)
	info += fmt.Sprintf("Status: %s\n", pod.Status)
	return info
}

func formatK8sDeploymentInfo(deployment k8sdata.DeploymentInfo) string {
	info := fmt.Sprintf("Kubernetes Deployment: %s\n", deployment.Name)
	info += fmt.Sprintf("Namespace: %s\n", deployment.Namespace)

	if len(deployment.Containers) > 0 {
		info += "Containers:\n"
		for i, container := range deployment.Containers {
			info += fmt.Sprintf("  - %s\n", container)
			if i < len(deployment.Images) {
				info += fmt.Sprintf("    Image: %s\n", deployment.Images[i])
			}
		}
	}

	return info
}

func formatK8sStatefulSetInfo(sts k8sdata.StatefulSetInfo) string {
	info := fmt.Sprintf("Kubernetes StatefulSet: %s\n", sts.Name)
	info += fmt.Sprintf("Namespace: %s\n", sts.Namespace)
	info += fmt.Sprintf("Ready Replicas: %d\n", sts.ReadyReplicas)
	info += fmt.Sprintf("Image: %s\n", sts.Image)
	return info
}

func formatK8sServiceInfo(svc k8sdata.ServiceInfo) string {
	info := fmt.Sprintf("Kubernetes Service: %s\n", svc.Name)
	info += fmt.Sprintf("Namespace: %s\n", svc.Namespace)
	info += fmt.Sprintf("Type: %s\n", svc.Type)
	info += fmt.Sprintf("Cluster IP: %s\n", svc.ClusterIP)
	info += fmt.Sprintf("Ports: %s\n", svc.Ports)
	return info
}

func formatK8sPVInfo(pv k8sdata.PersistentVolumeInfo) string {
	info := fmt.Sprintf("Kubernetes PersistentVolume: %s\n", pv.Name)
	info += fmt.Sprintf("Capacity: %s\n", pv.Capacity)
	info += fmt.Sprintf("Access Modes: %s\n", pv.AccessModes)
	info += fmt.Sprintf("Status: %s\n", pv.Status)
	info += fmt.Sprintf("Claim: %s\n", pv.AssociatedClaim)
	info += fmt.Sprintf("Storage Class: %s\n", pv.StorageClass)
	return info
}

func formatK8sPVCInfo(pvc k8sdata.PersistentVolumeClaimInfo) string {
	info := fmt.Sprintf("Kubernetes PersistentVolumeClaim: %s\n", pvc.Name)
	info += fmt.Sprintf("Namespace: %s\n", pvc.Namespace)
	info += fmt.Sprintf("Status: %s\n", pvc.Status)
	info += fmt.Sprintf("Volume: %s\n", pvc.Volume)
	info += fmt.Sprintf("Capacity: %s\n", pvc.Capacity)
	info += fmt.Sprintf("Access Mode: %s\n", pvc.AccessMode)
	info += fmt.Sprintf("Storage Class: %s\n", pvc.StorageClass)
	return info
}

func formatK8sStorageClassInfo(sc k8sdata.StorageClassInfo) string {
	info := fmt.Sprintf("Kubernetes StorageClass: %s\n", sc.Name)
	info += fmt.Sprintf("Provisioner: %s\n", sc.Provisioner)
	info += fmt.Sprintf("Volume Expansion: %s\n", sc.VolumeExpansion)
	return info
}

func formatK8sVolumeSnapshotClassInfo(vsc k8sdata.VolumeSnapshotClassInfo) string {
	info := fmt.Sprintf("Kubernetes VolumeSnapshotClass: %s\n", vsc.Name)
	info += fmt.Sprintf("Driver: %s\n", vsc.Driver)
	return info
}

func formatK8sVolumeSnapshotInfo(vs k8sdata.VolumeSnapshotInfo) string {
	info := fmt.Sprintf("Kubernetes VolumeSnapshot: %s\n", vs.Name)
	info += fmt.Sprintf("Namespace: %s\n", vs.Namespace)
	info += fmt.Sprintf("Volume: %s\n", vs.Volume)
	info += fmt.Sprintf("Creation Timestamp: %s\n", vs.CreationTimestamp)
	info += fmt.Sprintf("Restore Size: %s\n", vs.RestoreSize)
	info += fmt.Sprintf("State: %s\n", vs.State)
	return info
}

func formatK8sDataVolumeInfo(dv k8sdata.DataVolumeInfo) string {
	info := fmt.Sprintf("Kubernetes DataVolume: %s\n", dv.Name)
	info += fmt.Sprintf("Namespace: %s\n", dv.Namespace)
	info += fmt.Sprintf("Phase: %s\n", dv.Phase)
	info += fmt.Sprintf("Size: %s\n", dv.Size)
	info += fmt.Sprintf("Source Type: %s\n", dv.SourceType)
	info += fmt.Sprintf("Source Info: %s\n", dv.SourceInfo)
	info += fmt.Sprintf("Age: %s\n", dv.Age)
	return info
}

func formatK8sVMInfo(vm k8sdata.VirtualMachineInfo) string {
	info := fmt.Sprintf("Kubernetes VirtualMachine: %s\n", vm.Name)
	info += fmt.Sprintf("Namespace: %s\n", vm.Namespace)
	info += fmt.Sprintf("Status: %s\n", vm.Status)
	info += fmt.Sprintf("Ready: %t\n", vm.Ready)
	info += fmt.Sprintf("CPU: %s\n", vm.CPU)
	info += fmt.Sprintf("Memory: %s\n", vm.Memory)
	info += fmt.Sprintf("Run Strategy: %s\n", vm.RunStrategy)

	if len(vm.Storage) > 0 {
		info += "Storage:\n"
		for _, storage := range vm.Storage {
			info += fmt.Sprintf("  - %s\n", storage)
		}
	}

	if len(vm.DataVolumes) > 0 {
		info += "Data Volumes:\n"
		for _, dv := range vm.DataVolumes {
			info += fmt.Sprintf("  - %s\n", dv)
		}
	}

	info += fmt.Sprintf("Age: %s\n", vm.Age)
	return info
}

func formatK8sCRDInfo(crd k8sdata.CRDInfo) string {
	info := fmt.Sprintf("Kubernetes CustomResourceDefinition: %s\n", crd.Name)
	info += fmt.Sprintf("Group: %s\n", crd.Group)
	info += fmt.Sprintf("Version: %s\n", crd.Version)
	info += fmt.Sprintf("Scope: %s\n", crd.Scope)
	return info
}

// AWS resource formatters
func formatEC2Info(instance aws.EC2InstanceInfo) string {
	info := fmt.Sprintf("EC2 Instance: %s\n", instance.InstanceID)
	info += fmt.Sprintf("Name: %s\n", instance.Name)
	info += fmt.Sprintf("Type: %s\n", instance.Type)
	info += fmt.Sprintf("State: %s\n", instance.State)
	info += fmt.Sprintf("Region: %s\n", instance.Region)
	return info
}

func formatS3BucketInfo(bucket aws.S3BucketInfo) string {
	info := fmt.Sprintf("S3 Bucket: %s\n", bucket.Name)
	info += fmt.Sprintf("Region: %s\n", bucket.Region)
	if bucket.Immutable {
		info += "Immutable: Yes\n"
	} else {
		info += "Immutable: No\n"
	}
	return info
}

func formatRDSInfo(rds aws.RDSInstanceInfo) string {
	info := fmt.Sprintf("RDS Instance: %s\n", rds.InstanceID)
	info += fmt.Sprintf("Engine: %s\n", rds.Engine)
	info += fmt.Sprintf("Status: %s\n", rds.Status)
	info += fmt.Sprintf("Region: %s\n", rds.Region)
	return info
}

func formatDynamoDBInfo(table aws.DynamoDBTableInfo) string {
	info := fmt.Sprintf("DynamoDB Table: %s\n", table.TableName)
	info += fmt.Sprintf("Status: %s\n", table.Status)
	info += fmt.Sprintf("Region: %s\n", table.Region)
	return info
}

func formatVPCInfo(vpc aws.VPCInfo) string {
	info := fmt.Sprintf("VPC: %s\n", vpc.VPCID)
	info += fmt.Sprintf("State: %s\n", vpc.State)
	info += fmt.Sprintf("Region: %s\n", vpc.Region)
	return info
}

// GCP resource formatters
func formatGCPInstanceInfo(instance gcp.ComputeInstanceInfo) string {
	info := fmt.Sprintf("GCP Compute Instance: %s\n", instance.Name)
	info += fmt.Sprintf("Zone: %s\n", instance.Zone)
	info += fmt.Sprintf("Machine Type: %s\n", instance.MachineType)
	info += fmt.Sprintf("Status: %s\n", instance.Status)
	info += fmt.Sprintf("Project: %s\n", instance.Project)
	return info
}

func formatGCSBucketInfo(bucket gcp.GCSBucketInfo) string {
	info := fmt.Sprintf("GCS Bucket: %s\n", bucket.Name)
	info += fmt.Sprintf("Location: %s\n", bucket.Location)
	info += fmt.Sprintf("Storage Class: %s\n", bucket.StorageClass)
	info += fmt.Sprintf("Project: %s\n", bucket.Project)
	return info
}

func formatCloudSQLInfo(sql gcp.CloudSQLInstanceInfo) string {
	info := fmt.Sprintf("Cloud SQL Instance: %s\n", sql.Name)
	info += fmt.Sprintf("Database Version: %s\n", sql.DatabaseVersion)
	info += fmt.Sprintf("Region: %s\n", sql.Region)
	info += fmt.Sprintf("Tier: %s\n", sql.Tier)
	info += fmt.Sprintf("Status: %s\n", sql.Status)
	info += fmt.Sprintf("Project: %s\n", sql.Project)
	return info
}

func formatCloudRunInfo(service gcp.CloudRunServiceInfo) string {
	info := fmt.Sprintf("Cloud Run Service: %s\n", service.Name)
	info += fmt.Sprintf("Region: %s\n", service.Region)
	info += fmt.Sprintf("Project: %s\n", service.Project)
	return info
}

func formatCloudFunctionInfo(function gcp.CloudFunctionInfo) string {
	info := fmt.Sprintf("Cloud Function: %s\n", function.Name)
	info += fmt.Sprintf("Region: %s\n", function.Region)
	info += fmt.Sprintf("Runtime: %s\n", function.Runtime)
	info += fmt.Sprintf("Status: %s\n", function.Status)
	info += fmt.Sprintf("Entry Point: %s\n", function.EntryPoint)
	info += fmt.Sprintf("Available Memory: %s\n", function.AvailableMemory)
	info += fmt.Sprintf("Project: %s\n", function.Project)
	return info
}

// Azure resource formatters
func formatAzureVMInfo(vm map[string]interface{}) string {
	info := fmt.Sprintf("Azure Virtual Machine: %v\n", vm["name"])
	info += fmt.Sprintf("Resource Group: %v\n", vm["resourceGroup"])
	info += fmt.Sprintf("Location: %v\n", vm["location"])
	info += fmt.Sprintf("Size: %v\n", vm["size"])
	info += fmt.Sprintf("OS: %v\n", vm["os"])
	return info
}

func formatAzureDiskInfo(disk map[string]interface{}) string {
	info := fmt.Sprintf("Azure Managed Disk: %v\n", disk["name"])
	info += fmt.Sprintf("Resource Group: %v\n", disk["resourceGroup"])
	info += fmt.Sprintf("Location: %v\n", disk["location"])
	info += fmt.Sprintf("Size: %v\n", disk["diskSizeGB"])
	info += fmt.Sprintf("State: %v\n", disk["diskState"])
	return info
}

func formatAzureStorageAccountInfo(account map[string]interface{}) string {
	info := fmt.Sprintf("Azure Storage Account: %v\n", account["name"])
	info += fmt.Sprintf("Resource Group: %v\n", account["resourceGroup"])
	info += fmt.Sprintf("Location: %v\n", account["location"])
	info += fmt.Sprintf("Type: %v\n", account["accountType"])
	return info
}

// Terraform resource formatters
func formatTerraformResourceInfo(resource terraform.ResourceInfo) string {
	info := fmt.Sprintf("Terraform Resource: %s.%s\n", resource.Type, resource.Name)
	info += fmt.Sprintf("Provider: %s\n", resource.Provider)
	info += fmt.Sprintf("Mode: %s\n", resource.Mode)
	if resource.Module != "" {
		info += fmt.Sprintf("Module: %s\n", resource.Module)
	}

	if len(resource.Attributes) > 0 {
		info += "Attributes:\n"
		var attrKeys []string
		for k := range resource.Attributes {
			attrKeys = append(attrKeys, k)
		}
		sort.Strings(attrKeys)

		for _, k := range attrKeys {
			if len(resource.Attributes[k]) < 50 {
				info += fmt.Sprintf("  %s = %s\n", k, resource.Attributes[k])
			} else {
				info += fmt.Sprintf("  %s = %s...\n", k, resource.Attributes[k][:50])
			}
		}
	}

	return info
}

func formatTerraformProviderInfo(provider terraform.ProviderInfo) string {
	info := fmt.Sprintf("Terraform Provider: %s\n", provider.Name)
	info += fmt.Sprintf("Version: %s\n", provider.Version)
	return info
}

func formatTerraformOutputInfo(output terraform.OutputInfo) string {
	info := fmt.Sprintf("Terraform Output: %s\n", output.Name)
	info += fmt.Sprintf("Type: %s\n", output.Type)
	info += fmt.Sprintf("Value: %v\n", output.Value)
	return info
}

// Vault formatters
func formatVaultServerInfo(info vault.ServerInfo) string {
	result := fmt.Sprintf("Vault Server:\n")
	result += fmt.Sprintf("Version: %s\n", info.Version)
	result += fmt.Sprintf("Cluster Name: %s\n", info.ClusterName)
	result += fmt.Sprintf("Cluster ID: %s\n", info.ClusterID)

	if info.Sealed {
		result += "Status: Sealed\n"
	} else {
		result += "Status: Unsealed\n"
	}

	return result
}

func formatVaultAuthInfo(auth vault.AuthMethodInfo) string {
	info := fmt.Sprintf("Vault Auth Method: %s\n", auth.Path)
	info += fmt.Sprintf("Type: %s\n", auth.Type)
	if auth.Description != "" {
		info += fmt.Sprintf("Description: %s\n", auth.Description)
	}
	return info
}

func formatVaultSecretEngineInfo(engine vault.SecretEngineInfo) string {
	info := fmt.Sprintf("Vault Secret Engine: %s\n", engine.Path)
	info += fmt.Sprintf("Type: %s\n", engine.Type)
	if engine.Description != "" {
		info += fmt.Sprintf("Description: %s\n", engine.Description)
	}
	return info
}

// Veeam formatters
func formatVeeamJobInfo(job map[string]interface{}) string {
	info := fmt.Sprintf("Veeam Backup Job: %v\n", job["name"])

	if schedule, ok := job["schedule"].(map[string]interface{}); ok {
		info += fmt.Sprintf("Schedule: %v\n", schedule["type"])
	}

	if repository, ok := job["repository"].(string); ok {
		info += fmt.Sprintf("Repository: %s\n", repository)
	}

	return info
}

// Generic map formatter
func formatMapAsText(title string, data map[string]interface{}) string {
	info := fmt.Sprintf("%s:\n", title)

	var keys []string
	for k := range data {
		if k != "properties" && k != "tags" && k != "metadata" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		if str, ok := v.(string); ok {
			info += fmt.Sprintf("  %s: %s\n", strings.Title(k), str)
		} else {
			info += fmt.Sprintf("  %s: %v\n", strings.Title(k), v)
		}
	}

	return info
}

// Snapshot formatter
func formatSnapshotInfo(platform, snapshotType string, snapshot map[string]string) string {
	info := fmt.Sprintf("%s %s Snapshot:\n", platform, snapshotType)

	// Sort and display keys
	var keys []string
	for k := range snapshot {
		keys = append(keys, k)
	}

	// Important keys to show first
	importantKeys := []string{"ID", "SnapshotId", "SnapshotID", "Name", "CreationDate", "Size", "State", "Status", "VolumeSize", "DiskSizeGB", "StartTime", "TimeCreated"}
	for _, key := range importantKeys {
		if value, exists := snapshot[key]; exists {
			info += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}

	// Add remaining keys
	sort.Strings(keys)
	for _, k := range keys {
		isImportantKey := false
		for _, key := range importantKeys {
			if k == key {
				isImportantKey = true
				break
			}
		}

		if !isImportantKey {
			info += fmt.Sprintf("  %s: %s\n", k, snapshot[k])
		}
	}

	return info
}

// Generic content formatter
func formatGenericContent(title string, data interface{}) string {
	result := fmt.Sprintf("%s:\n", title)

	switch v := data.(type) {
	case map[string]interface{}:
		var keys []string
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			val := v[k]
			result += fmt.Sprintf("  %s: %v\n", k, val)
		}
	case []interface{}:
		for i, item := range v {
			if i >= 10 { // Limit number of items
				result += "  ... (more items)\n"
				break
			}
			result += fmt.Sprintf("  [%d]: %v\n", i, item)
		}
	default:
		result += fmt.Sprintf("  %v\n", v)
	}

	return result
}
