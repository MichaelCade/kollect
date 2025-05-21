package mcp

import (
	"fmt"
	"log"
	"strings"
	"time"

	kollect "github.com/michaelcade/kollect/api/v1"
	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/terraform"
	"github.com/michaelcade/kollect/pkg/vault"
	"github.com/michaelcade/kollect/pkg/veeam"
)

// ProcessData converts collected data to MCP documents
func ProcessData(data interface{}, sourceType string) []MCPDocument {
	var documents []MCPDocument

	log.Printf("Processing %s data for MCP, data type: %T", sourceType, data)

	// Detailed logging of the data structure
	switch sourceType {
	case "kubernetes":
		if k8sData, ok := data.(kollect.K8sData); ok {
			documents = append(documents, processK8sData(k8sData)...)
			log.Printf("K8s data found: %d nodes, %d pods", len(k8sData.Nodes), len(k8sData.Pods))
		} else {
			log.Printf("WARNING: Data doesn't match expected K8sData type: %T", data)
		}
	case "aws":
		log.Printf("AWS data structure: %T", data)
	case "azure":
		log.Printf("Azure data structure: %T", data)
		if azureData, ok := data.(map[string]interface{}); ok {
			documents = append(documents, processAzureDataMap(azureData)...)
		}
	case "gcp":
		if gcpData, ok := data.(gcp.GCPData); ok {
			documents = append(documents, processGCPData(gcpData)...)
		}
	case "terraform":
		if tfData, ok := data.(terraform.TerraformData); ok {
			documents = append(documents, processTerraformData(tfData)...)
		}
	case "vault":
		if vaultData, ok := data.(vault.VaultData); ok {
			documents = append(documents, processVaultData(vaultData)...)
		}
	case "veeam":
		if veeamData, ok := data.(veeam.VeeamData); ok {
			documents = append(documents, processVeeamData(veeamData)...)
		}
	case "snapshots":
		if snapData, ok := data.(map[string]interface{}); ok {
			documents = append(documents, processSnapshotData(snapData)...)
		}
	default:
		// Try to determine data type from the data itself
		documents = detectAndProcessData(data)
	}

	log.Printf("MCP: Processing data of type '%s', found %d documents", sourceType, len(documents))

	if len(documents) > 0 {
		// Ensure MCP is initialized before indexing
		if docStore == nil {
			log.Println("MCP not initialized, initializing now...")
			InitMCP()
		}

		// Index the documents to make them searchable
		IndexDocuments(documents)
		log.Printf("MCP: Successfully indexed %d documents", len(documents))
	} else {
		log.Printf("MCP: No documents were generated for source type '%s'", sourceType)
	}

	return documents
}

// Helper functions for each data type
func processAWSData(data aws.AWSData) []MCPDocument {
	var docs []MCPDocument

	// Process EC2 instances
	for _, instance := range data.EC2Instances {
		doc := MCPDocument{
			ID:      fmt.Sprintf("aws-ec2-%s", instance.InstanceID),
			Content: formatEC2Info(instance),
			Metadata: map[string]interface{}{
				"type":          "ec2",
				"region":        instance.Region,
				"instance_type": instance.Type,
				"state":         instance.State,
			},
			Source:     "aws",
			SourceType: "ec2_instance",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process S3 buckets
	for _, bucket := range data.S3Buckets {
		doc := MCPDocument{
			ID:      fmt.Sprintf("aws-s3-%s", bucket.Name),
			Content: formatS3BucketInfo(bucket),
			Metadata: map[string]interface{}{
				"type":      "s3",
				"region":    bucket.Region,
				"immutable": bucket.Immutable,
			},
			Source:     "aws",
			SourceType: "s3_bucket",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process RDS instances
	for _, rds := range data.RDSInstances {
		doc := MCPDocument{
			ID:      fmt.Sprintf("aws-rds-%s", rds.InstanceID),
			Content: formatRDSInfo(rds),
			Metadata: map[string]interface{}{
				"type":   "rds",
				"region": rds.Region,
				"engine": rds.Engine,
				"status": rds.Status,
			},
			Source:     "aws",
			SourceType: "rds_instance",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process DynamoDB tables
	for _, table := range data.DynamoDBTables {
		doc := MCPDocument{
			ID:      fmt.Sprintf("aws-dynamodb-%s", table.TableName),
			Content: formatDynamoDBInfo(table),
			Metadata: map[string]interface{}{
				"type":   "dynamodb",
				"region": table.Region,
			},
			Source:     "aws",
			SourceType: "dynamodb_table",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process VPCs
	for _, vpc := range data.VPCs {
		doc := MCPDocument{
			ID:      fmt.Sprintf("aws-vpc-%s", vpc.VPCID),
			Content: formatVPCInfo(vpc),
			Metadata: map[string]interface{}{
				"type":   "vpc",
				"region": vpc.Region,
				"state":  vpc.State,
			},
			Source:     "aws",
			SourceType: "vpc",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

// Process Azure data that's in map format
func processAzureDataMap(data map[string]interface{}) []MCPDocument {
	var docs []MCPDocument

	// Process resource groups
	if rgs, ok := data["AzureResourceGroups"].([]interface{}); ok {
		for _, rg := range rgs {
			if rgMap, ok := rg.(map[string]interface{}); ok {
				name := rgMap["name"]
				location := rgMap["location"]
				doc := MCPDocument{
					ID:      fmt.Sprintf("azure-rg-%v", name),
					Content: formatMapAsText("Azure Resource Group", rgMap),
					Metadata: map[string]interface{}{
						"type":     "resourcegroup",
						"location": location,
					},
					Source:     "azure",
					SourceType: "resource_group",
					CreatedAt:  time.Now(),
				}
				docs = append(docs, doc)
			}
		}
	}

	// Process VMs
	if vms, ok := data["AzureVMs"].([]interface{}); ok {
		for _, vm := range vms {
			if vmMap, ok := vm.(map[string]interface{}); ok {
				name := vmMap["name"]
				doc := MCPDocument{
					ID:      fmt.Sprintf("azure-vm-%v", name),
					Content: formatMapAsText("Azure Virtual Machine", vmMap),
					Metadata: map[string]interface{}{
						"type": "virtualmachine",
					},
					Source:     "azure",
					SourceType: "virtual_machine",
					CreatedAt:  time.Now(),
				}
				docs = append(docs, doc)
			}
		}
	}

	// Process storage accounts
	if accounts, ok := data["AzureStorageAccounts"].([]interface{}); ok {
		for _, acct := range accounts {
			if acctMap, ok := acct.(map[string]interface{}); ok {
				name := acctMap["name"]
				doc := MCPDocument{
					ID:      fmt.Sprintf("azure-storage-%v", name),
					Content: formatMapAsText("Azure Storage Account", acctMap),
					Metadata: map[string]interface{}{
						"type": "storageaccount",
					},
					Source:     "azure",
					SourceType: "storage_account",
					CreatedAt:  time.Now(),
				}
				docs = append(docs, doc)
			}
		}
	}

	// Add other Azure resource types as needed

	return docs
}

func processGCPData(data gcp.GCPData) []MCPDocument {
	var docs []MCPDocument

	// Process Compute Instances
	for _, instance := range data.ComputeInstances {
		doc := MCPDocument{
			ID:      fmt.Sprintf("gcp-compute-%s", instance.Name),
			Content: formatGCPInstanceInfo(instance),
			Metadata: map[string]interface{}{
				"type":         "compute",
				"zone":         instance.Zone,
				"machine_type": instance.MachineType,
				"status":       instance.Status,
				"project":      instance.Project,
			},
			Source:     "gcp",
			SourceType: "compute_instance",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process GCS Buckets
	for _, bucket := range data.GCSBuckets {
		doc := MCPDocument{
			ID:      fmt.Sprintf("gcp-gcs-%s", bucket.Name),
			Content: formatGCSBucketInfo(bucket),
			Metadata: map[string]interface{}{
				"type":          "storage",
				"location":      bucket.Location,
				"storage_class": bucket.StorageClass,
				"project":       bucket.Project,
			},
			Source:     "gcp",
			SourceType: "gcs_bucket",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process Cloud SQL Instances
	for _, sql := range data.CloudSQLInstances {
		doc := MCPDocument{
			ID:      fmt.Sprintf("gcp-sql-%s", sql.Name),
			Content: formatCloudSQLInfo(sql),
			Metadata: map[string]interface{}{
				"type":             "cloudsql",
				"region":           sql.Region,
				"database_version": sql.DatabaseVersion,
				"tier":             sql.Tier,
				"status":           sql.Status,
				"project":          sql.Project,
			},
			Source:     "gcp",
			SourceType: "cloudsql_instance",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process Cloud Run Services
	for _, service := range data.CloudRunServices {
		doc := MCPDocument{
			ID:      fmt.Sprintf("gcp-run-%s", service.Name),
			Content: formatCloudRunInfo(service),
			Metadata: map[string]interface{}{
				"type":    "cloudrun",
				"region":  service.Region,
				"project": service.Project,
			},
			Source:     "gcp",
			SourceType: "cloudrun_service",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process Cloud Functions
	for _, function := range data.CloudFunctions {
		doc := MCPDocument{
			ID:      fmt.Sprintf("gcp-function-%s", function.Name),
			Content: formatCloudFunctionInfo(function),
			Metadata: map[string]interface{}{
				"type":    "cloudfunction",
				"region":  function.Region,
				"runtime": function.Runtime,
				"status":  function.Status,
				"project": function.Project,
			},
			Source:     "gcp",
			SourceType: "cloud_function",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

func processK8sData(data kollect.K8sData) []MCPDocument {
	var docs []MCPDocument

	// Process nodes
	for _, node := range data.Nodes {
		doc := MCPDocument{
			ID:      fmt.Sprintf("k8s-node-%s", node.Name),
			Content: formatK8sNodeInfo(node),
			Metadata: map[string]interface{}{
				"type":    "node",
				"roles":   node.Roles,
				"version": node.Version,
			},
			Source:     "kubernetes",
			SourceType: "node",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process pods
	for _, pod := range data.Pods {
		doc := MCPDocument{
			ID:      fmt.Sprintf("k8s-pod-%s-%s", pod.Namespace, pod.Name),
			Content: formatK8sPodInfo(pod),
			Metadata: map[string]interface{}{
				"type":      "pod",
				"namespace": pod.Namespace,
				"status":    pod.Status,
			},
			Source:     "kubernetes",
			SourceType: "pod",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process persistent volumes
	for _, pv := range data.PersistentVolumes {
		doc := MCPDocument{
			ID:      fmt.Sprintf("k8s-pv-%s", pv.Name),
			Content: formatK8sPVInfo(pv),
			Metadata: map[string]interface{}{
				"type":          "persistentvolume",
				"capacity":      pv.Capacity,
				"status":        pv.Status,
				"storage_class": pv.StorageClass,
			},
			Source:     "kubernetes",
			SourceType: "persistent_volume",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

func processTerraformData(data terraform.TerraformData) []MCPDocument {
	var docs []MCPDocument

	// Process resources
	for _, resource := range data.Resources {
		doc := MCPDocument{
			ID:      fmt.Sprintf("terraform-resource-%s-%s", resource.Type, resource.Name),
			Content: formatTerraformResourceInfo(resource),
			Metadata: map[string]interface{}{
				"type":     resource.Type,
				"provider": resource.Provider,
				"mode":     resource.Mode,
				"module":   resource.Module,
			},
			Source:     "terraform",
			SourceType: "resource",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process providers
	for _, provider := range data.Providers {
		doc := MCPDocument{
			ID:      fmt.Sprintf("terraform-provider-%s", provider.Name),
			Content: formatTerraformProviderInfo(provider),
			Metadata: map[string]interface{}{
				"name":    provider.Name,
				"version": provider.Version,
			},
			Source:     "terraform",
			SourceType: "provider",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process outputs
	for _, output := range data.Outputs {
		doc := MCPDocument{
			ID:      fmt.Sprintf("terraform-output-%s", output.Name),
			Content: formatTerraformOutputInfo(output),
			Metadata: map[string]interface{}{
				"name": output.Name,
				"type": output.Type,
			},
			Source:     "terraform",
			SourceType: "output",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

func processVaultData(data vault.VaultData) []MCPDocument {
	var docs []MCPDocument

	// Process server info
	serverInfo := data.ServerInfo
	doc := MCPDocument{
		ID:      fmt.Sprintf("vault-server-%s", serverInfo.Version),
		Content: formatVaultServerInfo(serverInfo),
		Metadata: map[string]interface{}{
			"version":      serverInfo.Version,
			"cluster_name": serverInfo.ClusterName,
			"cluster_id":   serverInfo.ClusterID,
			"sealed":       serverInfo.Sealed,
		},
		Source:     "vault",
		SourceType: "server_info",
		CreatedAt:  time.Now(),
	}
	docs = append(docs, doc)

	// Process auth methods
	for _, auth := range data.AuthMethods {
		doc := MCPDocument{
			ID:      fmt.Sprintf("vault-auth-%s", auth.Path),
			Content: formatVaultAuthInfo(auth),
			Metadata: map[string]interface{}{
				"type":        auth.Type,
				"description": auth.Description,
			},
			Source:     "vault",
			SourceType: "auth_method",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process secret engines
	for _, engine := range data.SecretEngines {
		doc := MCPDocument{
			ID:      fmt.Sprintf("vault-engine-%s", engine.Path),
			Content: formatVaultSecretEngineInfo(engine),
			Metadata: map[string]interface{}{
				"type":        engine.Type,
				"description": engine.Description,
			},
			Source:     "vault",
			SourceType: "secret_engine",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

// Replace the processVeeamData function with this simplified implementation
func processVeeamData(data veeam.VeeamData) []MCPDocument {
	var docs []MCPDocument

	// Create a direct reference to veeam.go to understand the structure
	// Most likely ServerInfo is already map[string]interface{}, not an interface{} type

	// Process server info - skip type assertions
	if data.ServerInfo != nil {
		// Convert to string representation
		serverInfoStr := fmt.Sprintf("%v", data.ServerInfo)
		doc := MCPDocument{
			ID:      "veeam-server-info",
			Content: fmt.Sprintf("Veeam Backup & Replication Server:\n%s", serverInfoStr),
			Metadata: map[string]interface{}{
				"type": "server_info",
			},
			Source:     "veeam",
			SourceType: "server_info",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process backup jobs - without type assertions
	for i, job := range data.BackupJobs {
		jobID := fmt.Sprintf("%d", i)

		// Create a simple string representation
		jobStr := fmt.Sprintf("%v", job)

		doc := MCPDocument{
			ID:      fmt.Sprintf("veeam-job-%s", jobID),
			Content: fmt.Sprintf("Veeam Backup Job:\n%s", jobStr),
			Metadata: map[string]interface{}{
				"type": "backup_job",
			},
			Source:     "veeam",
			SourceType: "backup_job",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	// Process repositories - similar approach
	for i, repo := range data.Repositories {
		repoID := fmt.Sprintf("%d", i)

		// Create a simple string representation
		repoStr := fmt.Sprintf("%v", repo)

		doc := MCPDocument{
			ID:      fmt.Sprintf("veeam-repository-%s", repoID),
			Content: fmt.Sprintf("Veeam Repository:\n%s", repoStr),
			Metadata: map[string]interface{}{
				"type": "repository",
			},
			Source:     "veeam",
			SourceType: "repository",
			CreatedAt:  time.Now(),
		}
		docs = append(docs, doc)
	}

	return docs
}

func processSnapshotData(data map[string]interface{}) []MCPDocument {
	var docs []MCPDocument

	// Process each platform's snapshots
	for platform, platformData := range data {
		platformMap, ok := platformData.(map[string]interface{})
		if !ok {
			continue
		}

		// Process each snapshot type
		for snapshotType, snapshots := range platformMap {
			snapshotsList, ok := snapshots.([]interface{})
			if !ok {
				continue
			}

			for i, snapshot := range snapshotsList {
				snapshotMap, ok := snapshot.(map[string]string)
				if !ok {
					// Try to convert if it's not already a map[string]string
					snapshotMapGeneric, isMap := snapshot.(map[string]interface{})
					if isMap {
						snapshotMap = make(map[string]string)
						for k, v := range snapshotMapGeneric {
							snapshotMap[k] = fmt.Sprintf("%v", v)
						}
					} else {
						continue
					}
				}

				// Try to get ID from the snapshot
				snapshotID := fmt.Sprintf("%d", i)
				if id, ok := snapshotMap["ID"]; ok && id != "" {
					snapshotID = id
				} else if id, ok := snapshotMap["SnapshotId"]; ok && id != "" {
					snapshotID = id
				} else if id, ok := snapshotMap["SnapshotID"]; ok && id != "" {
					snapshotID = id
				} else if id, ok := snapshotMap["Name"]; ok && id != "" {
					snapshotID = id
				}

				doc := MCPDocument{
					ID:      fmt.Sprintf("%s-%s-%s", platform, snapshotType, snapshotID),
					Content: formatSnapshotInfo(platform, snapshotType, snapshotMap),
					Metadata: map[string]interface{}{
						"platform": platform,
						"type":     snapshotType,
					},
					Source:     platform,
					SourceType: fmt.Sprintf("%s_snapshot", snapshotType),
					CreatedAt:  time.Now(),
				}
				docs = append(docs, doc)
			}
		}
	}

	return docs
}

func detectAndProcessData(data interface{}) []MCPDocument {
	var documents []MCPDocument

	// Try to determine what kind of data we have

	// Map-based data (could be snapshots or generic)
	if mapData, ok := data.(map[string]interface{}); ok {
		// Check if it's snapshot data
		if _, hasAWS := mapData["aws"]; hasAWS {
			if _, hasAzure := mapData["azure"]; hasAzure {
				if _, hasGCP := mapData["gcp"]; hasGCP {
					return processSnapshotData(mapData)
				}
			}
		}

		// Azure data detection
		if _, hasVMs := mapData["AzureVMs"]; hasVMs {
			return processAzureDataMap(mapData)
		}

		// Generic map processing
		for key, value := range mapData {
			doc := MCPDocument{
				ID:      fmt.Sprintf("generic-%s", key),
				Content: formatGenericMapInfo(key, value),
				Metadata: map[string]interface{}{
					"key": key,
				},
				Source:     "generic",
				SourceType: "map_data",
				CreatedAt:  time.Now(),
			}
			documents = append(documents, doc)
		}
	}

	return documents
}

// Helper function to format map as text
func formatMapAsText(title string, data map[string]interface{}) string {
	info := fmt.Sprintf("%s:\n", title)
	for k, v := range data {
		if k == "properties" || k == "tags" {
			continue // Skip complex nested objects
		}
		info += fmt.Sprintf("%s: %v\n", strings.Title(fmt.Sprintf("%v", k)), v)
	}
	return info
}

// Formatting functions for each resource type
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

func formatRDSInfo(instance aws.RDSInstanceInfo) string {
	info := fmt.Sprintf("RDS Instance: %s\n", instance.InstanceID)
	info += fmt.Sprintf("Engine: %s\n", instance.Engine)
	info += fmt.Sprintf("Status: %s\n", instance.Status)
	info += fmt.Sprintf("Region: %s\n", instance.Region)
	return info
}

func formatDynamoDBInfo(table aws.DynamoDBTableInfo) string {
	info := fmt.Sprintf("DynamoDB Table: %s\n", table.TableName)
	info += fmt.Sprintf("Region: %s\n", table.Region)
	info += fmt.Sprintf("Status: %s\n", table.Status)
	return info
}

func formatVPCInfo(vpc aws.VPCInfo) string {
	info := fmt.Sprintf("VPC: %s\n", vpc.VPCID)
	info += fmt.Sprintf("State: %s\n", vpc.State)
	info += fmt.Sprintf("Region: %s\n", vpc.Region)
	return info
}

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

func formatCloudSQLInfo(instance gcp.CloudSQLInstanceInfo) string {
	info := fmt.Sprintf("Cloud SQL Instance: %s\n", instance.Name)
	info += fmt.Sprintf("Database Version: %s\n", instance.DatabaseVersion)
	info += fmt.Sprintf("Region: %s\n", instance.Region)
	info += fmt.Sprintf("Tier: %s\n", instance.Tier)
	info += fmt.Sprintf("Status: %s\n", instance.Status)
	info += fmt.Sprintf("Project: %s\n", instance.Project)
	return info
}

func formatCloudRunInfo(service gcp.CloudRunServiceInfo) string {
	info := fmt.Sprintf("Cloud Run Service: %s\n", service.Name)
	info += fmt.Sprintf("Region: %s\n", service.Region)
	info += fmt.Sprintf("URL: %s\n", service.URL)
	info += fmt.Sprintf("Replicas: %d\n", service.Replicas)
	if service.Container != "" {
		info += fmt.Sprintf("Container: %s\n", service.Container)
	}
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

func formatK8sNodeInfo(node kollect.NodeInfo) string {
	info := fmt.Sprintf("Kubernetes Node: %s\n", node.Name)
	info += fmt.Sprintf("Roles: %s\n", node.Roles)
	info += fmt.Sprintf("Version: %s\n", node.Version)
	info += fmt.Sprintf("OS: %s\n", node.OSImage)
	info += fmt.Sprintf("Age: %s\n", node.Age)
	return info
}

func formatK8sPodInfo(pod kollect.PodsInfo) string {
	info := fmt.Sprintf("Kubernetes Pod: %s\n", pod.Name)
	info += fmt.Sprintf("Namespace: %s\n", pod.Namespace)
	info += fmt.Sprintf("Status: %s\n", pod.Status)
	return info
}

func formatK8sPVInfo(pv kollect.PersistentVolumeInfo) string {
	info := fmt.Sprintf("Kubernetes PersistentVolume: %s\n", pv.Name)
	info += fmt.Sprintf("Capacity: %s\n", pv.Capacity)
	info += fmt.Sprintf("Access Modes: %s\n", pv.AccessModes)
	info += fmt.Sprintf("Status: %s\n", pv.Status)
	info += fmt.Sprintf("Associated Claim: %s\n", pv.AssociatedClaim)
	info += fmt.Sprintf("Storage Class: %s\n", pv.StorageClass)
	return info
}

func formatTerraformResourceInfo(resource terraform.ResourceInfo) string {
	info := fmt.Sprintf("Terraform Resource: %s.%s\n", resource.Type, resource.Name)
	info += fmt.Sprintf("Provider: %s\n", resource.Provider)
	info += fmt.Sprintf("Mode: %s\n", resource.Mode)

	if resource.Module != "" {
		info += fmt.Sprintf("Module: %s\n", resource.Module)
	}

	if len(resource.Attributes) > 0 {
		info += "Attributes:\n"
		for k, v := range resource.Attributes {
			info += fmt.Sprintf("  %s = %s\n", k, v)
		}
	}

	if len(resource.Dependencies) > 0 {
		info += "Dependencies:\n"
		for _, dep := range resource.Dependencies {
			info += fmt.Sprintf("  - %s\n", dep)
		}
	}

	info += fmt.Sprintf("Status: %s\n", resource.Status)

	return info
}

func formatTerraformProviderInfo(provider terraform.ProviderInfo) string {
	info := fmt.Sprintf("Terraform Provider: %s\n", provider.Name)
	info += fmt.Sprintf("Version: %s\n", provider.Version)
	return info
}

func formatTerraformOutputInfo(output terraform.OutputInfo) string {
	info := fmt.Sprintf("Terraform Output: %s\n", output.Name)
	info += fmt.Sprintf("Value: %s\n", output.Value)
	info += fmt.Sprintf("Type: %s\n", output.Type)
	return info
}

// Format Vault ServerInfo
func formatVaultServerInfo(info vault.ServerInfo) string {
	result := "Vault Server:\n"

	if info.ClusterName != "" {
		result += fmt.Sprintf("Vault Server: %s\n", info.ClusterName)
	}

	if info.Version != "" {
		result += fmt.Sprintf("Version: %s\n", info.Version)
	}

	if info.ClusterID != "" {
		result += fmt.Sprintf("Cluster ID: %s\n", info.ClusterID)
	}

	if info.Sealed {
		result += "Status: Sealed\n"
	} else {
		result += "Status: Unsealed\n"
	}

	if info.Initialized {
		result += "Initialized: Yes\n"
	} else {
		result += "Initialized: No\n"
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

func formatVeeamServerInfo(info map[string]interface{}) string {
	result := "Veeam Backup & Replication Server:\n"

	if version, ok := info["version"].(string); ok {
		result += fmt.Sprintf("Version: %s\n", version)
	}

	if edition, ok := info["edition"].(string); ok {
		result += fmt.Sprintf("Edition: %s\n", edition)
	}

	if name, ok := info["name"].(string); ok {
		result += fmt.Sprintf("Server Name: %s\n", name)
	}

	return result
}

func formatVeeamJobInfo(job map[string]interface{}) string {
	info := "Veeam Backup Job:\n"

	// Use map access rather than type assertion on the map itself
	if name, ok := job["name"].(string); ok {
		info += fmt.Sprintf("Name: %s\n", name)
	}

	if id, ok := job["id"].(string); ok {
		info += fmt.Sprintf("ID: %s\n", id)
	}

	if jobType, ok := job["type"].(string); ok {
		info += fmt.Sprintf("Type: %s\n", jobType)
	}

	if scheduleObj, ok := job["schedule"]; ok {
		if schedule, ok := scheduleObj.(map[string]interface{}); ok {
			info += "Schedule:\n"
			if enabled, ok := schedule["enabled"].(bool); ok {
				info += fmt.Sprintf("  Enabled: %t\n", enabled)
			}
			if daily, ok := schedule["daily"].(string); ok {
				info += fmt.Sprintf("  Daily: %s\n", daily)
			}
		}
	}

	if status, ok := job["status"].(string); ok {
		info += fmt.Sprintf("Status: %s\n", status)
	}

	return info
}

func formatVeeamRepoInfo(repo map[string]interface{}) string {
	info := "Veeam Repository:\n"

	if name, ok := repo["name"].(string); ok {
		info += fmt.Sprintf("Name: %s\n", name)
	}

	if id, ok := repo["id"].(string); ok {
		info += fmt.Sprintf("ID: %s\n", id)
	}

	if repoType, ok := repo["type"].(string); ok {
		info += fmt.Sprintf("Type: %s\n", repoType)
	}

	if path, ok := repo["path"].(string); ok {
		info += fmt.Sprintf("Path: %s\n", path)
	}

	if capacity, ok := repo["capacity"].(string); ok {
		info += fmt.Sprintf("Capacity: %s\n", capacity)
	} else if capacity, ok := repo["capacity"].(float64); ok {
		info += fmt.Sprintf("Capacity: %.2f GB\n", capacity)
	}

	if free, ok := repo["freeSpace"].(string); ok {
		info += fmt.Sprintf("Free Space: %s\n", free)
	} else if free, ok := repo["freeSpace"].(float64); ok {
		info += fmt.Sprintf("Free Space: %.2f GB\n", free)
	}

	return info
}

func formatSnapshotInfo(platform string, snapshotType string, snapshot map[string]string) string {
	info := fmt.Sprintf("%s %s Snapshot:\n", platform, snapshotType)

	for key, value := range snapshot {
		// Skip empty values and some internal keys
		if value == "" || key == "id" {
			continue
		}

		// Format the key for better readability
		formattedKey := strings.ReplaceAll(key, "_", " ")
		formattedKey = strings.Title(formattedKey)

		info += fmt.Sprintf("%s: %s\n", formattedKey, value)
	}

	return info
}

func formatGenericMapInfo(key string, value interface{}) string {
	info := fmt.Sprintf("Resource: %s\n", key)

	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			info += fmt.Sprintf("%s: %v\n", k, val)
		}
	case []interface{}:
		info += fmt.Sprintf("Count: %d\n", len(v))
		for i, item := range v {
			if i >= 5 {
				info += "... (more items)\n"
				break
			}
			info += fmt.Sprintf("Item %d: %v\n", i+1, item)
		}
	default:
		info += fmt.Sprintf("Value: %v\n", v)
	}

	return info
}
