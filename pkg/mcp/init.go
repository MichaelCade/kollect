package mcp

import (
	"fmt"
	"time"

	k8sdata "github.com/michaelcade/kollect/api/v1"
	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/terraform"
	"github.com/michaelcade/kollect/pkg/vault"
	"github.com/michaelcade/kollect/pkg/veeam"
)

// RegisterHandlers sets up all resource type handlers
// This is called during InitMCP()
func RegisterHandlers() {
	// =====================================================================
	// KUBERNETES HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "node",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, node := range k8sData.Nodes {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "namespace",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, ns := range k8sData.Namespaces {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-namespace-%s", ns),
					Content: fmt.Sprintf("Kubernetes Namespace: %s\n", ns),
					Metadata: map[string]interface{}{
						"type": "namespace",
						"name": ns,
					},
					Source:     "kubernetes",
					SourceType: "namespace",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "pod",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, pod := range k8sData.Pods {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "deployment",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, deployment := range k8sData.Deployments {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-deployment-%s-%s", deployment.Namespace, deployment.Name),
					Content: formatK8sDeploymentInfo(deployment),
					Metadata: map[string]interface{}{
						"type":      "deployment",
						"namespace": deployment.Namespace,
					},
					Source:     "kubernetes",
					SourceType: "deployment",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "stateful_set",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, sts := range k8sData.StatefulSets {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-statefulset-%s-%s", sts.Namespace, sts.Name),
					Content: formatK8sStatefulSetInfo(sts),
					Metadata: map[string]interface{}{
						"type":      "statefulset",
						"namespace": sts.Namespace,
						"image":     sts.Image,
					},
					Source:     "kubernetes",
					SourceType: "stateful_set",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "service",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, svc := range k8sData.Services {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-service-%s-%s", svc.Namespace, svc.Name),
					Content: formatK8sServiceInfo(svc),
					Metadata: map[string]interface{}{
						"type":      "service",
						"namespace": svc.Namespace,
						"svc_type":  svc.Type,
					},
					Source:     "kubernetes",
					SourceType: "service",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "persistent_volume",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, pv := range k8sData.PersistentVolumes {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "persistent_volume_claim",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, pvc := range k8sData.PersistentVolumeClaims {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-pvc-%s-%s", pvc.Namespace, pvc.Name),
					Content: formatK8sPVCInfo(pvc),
					Metadata: map[string]interface{}{
						"type":          "persistentvolumeclaim",
						"namespace":     pvc.Namespace,
						"status":        pvc.Status,
						"volume":        pvc.Volume,
						"storage_class": pvc.StorageClass,
					},
					Source:     "kubernetes",
					SourceType: "persistent_volume_claim",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "storage_class",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, sc := range k8sData.StorageClasses {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-storageclass-%s", sc.Name),
					Content: formatK8sStorageClassInfo(sc),
					Metadata: map[string]interface{}{
						"type":             "storageclass",
						"provisioner":      sc.Provisioner,
						"volume_expansion": sc.VolumeExpansion,
					},
					Source:     "kubernetes",
					SourceType: "storage_class",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "volume_snapshot_class",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, vsc := range k8sData.VolumeSnapshotClasses {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-volumesnapshotclass-%s", vsc.Name),
					Content: formatK8sVolumeSnapshotClassInfo(vsc),
					Metadata: map[string]interface{}{
						"type":   "volumesnapshotclass",
						"driver": vsc.Driver,
					},
					Source:     "kubernetes",
					SourceType: "volume_snapshot_class",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "kubernetes",
		ResourceType: "volume_snapshot",
		ExtractFunc: func(data interface{}) []MCPDocument {
			k8sData, ok := data.(k8sdata.K8sData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, vs := range k8sData.VolumeSnapshots {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("k8s-volumesnapshot-%s-%s", vs.Namespace, vs.Name),
					Content: formatK8sVolumeSnapshotInfo(vs),
					Metadata: map[string]interface{}{
						"type":         "volumesnapshot",
						"namespace":    vs.Namespace,
						"volume":       vs.Volume,
						"restore_size": vs.RestoreSize,
						"state":        vs.State,
					},
					Source:     "kubernetes",
					SourceType: "volume_snapshot",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	// =====================================================================
	// AWS HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "aws",
		ResourceType: "ec2_instance",
		ExtractFunc: func(data interface{}) []MCPDocument {
			awsData, ok := data.(aws.AWSData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, instance := range awsData.EC2Instances {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "aws",
		ResourceType: "s3_bucket",
		ExtractFunc: func(data interface{}) []MCPDocument {
			awsData, ok := data.(aws.AWSData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, bucket := range awsData.S3Buckets {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "aws",
		ResourceType: "rds_instance",
		ExtractFunc: func(data interface{}) []MCPDocument {
			awsData, ok := data.(aws.AWSData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, rds := range awsData.RDSInstances {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "aws",
		ResourceType: "dynamodb_table",
		ExtractFunc: func(data interface{}) []MCPDocument {
			awsData, ok := data.(aws.AWSData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, table := range awsData.DynamoDBTables {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("aws-dynamodb-%s", table.TableName),
					Content: formatDynamoDBInfo(table),
					Metadata: map[string]interface{}{
						"type":   "dynamodb",
						"region": table.Region,
					},
					Source:     "aws",
					SourceType: "dynamodb_table",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "aws",
		ResourceType: "vpc",
		ExtractFunc: func(data interface{}) []MCPDocument {
			awsData, ok := data.(aws.AWSData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, vpc := range awsData.VPCs {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	// =====================================================================
	// AZURE HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "azure",
		ResourceType: "resource_group",
		ExtractFunc: func(data interface{}) []MCPDocument {
			// Azure data is typically a map[string]interface{}
			azureData, ok := data.(map[string]interface{})
			if !ok {
				return nil
			}

			var docs []MCPDocument
			if rgs, ok := azureData["AzureResourceGroups"].([]interface{}); ok {
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
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "azure",
		ResourceType: "virtual_machine",
		ExtractFunc: func(data interface{}) []MCPDocument {
			// Azure data is typically a map[string]interface{}
			azureData, ok := data.(map[string]interface{})
			if !ok {
				return nil
			}

			var docs []MCPDocument
			if vms, ok := azureData["AzureVMs"].([]interface{}); ok {
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
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "azure",
		ResourceType: "storage_account",
		ExtractFunc: func(data interface{}) []MCPDocument {
			// Azure data is typically a map[string]interface{}
			azureData, ok := data.(map[string]interface{})
			if !ok {
				return nil
			}

			var docs []MCPDocument
			if accts, ok := azureData["AzureStorageAccounts"].([]interface{}); ok {
				for _, acct := range accts {
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
			return docs
		},
	})

	// =====================================================================
	// GCP HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "gcp",
		ResourceType: "compute_instance",
		ExtractFunc: func(data interface{}) []MCPDocument {
			gcpData, ok := data.(gcp.GCPData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, instance := range gcpData.ComputeInstances {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "gcp",
		ResourceType: "gcs_bucket",
		ExtractFunc: func(data interface{}) []MCPDocument {
			gcpData, ok := data.(gcp.GCPData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, bucket := range gcpData.GCSBuckets {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "gcp",
		ResourceType: "cloudsql_instance",
		ExtractFunc: func(data interface{}) []MCPDocument {
			gcpData, ok := data.(gcp.GCPData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, sql := range gcpData.CloudSQLInstances {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "gcp",
		ResourceType: "cloudrun_service",
		ExtractFunc: func(data interface{}) []MCPDocument {
			gcpData, ok := data.(gcp.GCPData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, service := range gcpData.CloudRunServices {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "gcp",
		ResourceType: "cloud_function",
		ExtractFunc: func(data interface{}) []MCPDocument {
			gcpData, ok := data.(gcp.GCPData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, function := range gcpData.CloudFunctions {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	// =====================================================================
	// TERRAFORM HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "terraform",
		ResourceType: "resource",
		ExtractFunc: func(data interface{}) []MCPDocument {
			tfData, ok := data.(terraform.TerraformData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, resource := range tfData.Resources {
				docs = append(docs, MCPDocument{
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
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "terraform",
		ResourceType: "provider",
		ExtractFunc: func(data interface{}) []MCPDocument {
			tfData, ok := data.(terraform.TerraformData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, provider := range tfData.Providers {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("terraform-provider-%s", provider.Name),
					Content: formatTerraformProviderInfo(provider),
					Metadata: map[string]interface{}{
						"name":    provider.Name,
						"version": provider.Version,
					},
					Source:     "terraform",
					SourceType: "provider",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "terraform",
		ResourceType: "output",
		ExtractFunc: func(data interface{}) []MCPDocument {
			tfData, ok := data.(terraform.TerraformData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, output := range tfData.Outputs {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("terraform-output-%s", output.Name),
					Content: formatTerraformOutputInfo(output),
					Metadata: map[string]interface{}{
						"name": output.Name,
						"type": output.Type,
					},
					Source:     "terraform",
					SourceType: "output",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	// =====================================================================
	// VAULT HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "vault",
		ResourceType: "server_info",
		ExtractFunc: func(data interface{}) []MCPDocument {
			vaultData, ok := data.(vault.VaultData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			serverInfo := vaultData.ServerInfo
			docs = append(docs, MCPDocument{
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
			})
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "vault",
		ResourceType: "auth_method",
		ExtractFunc: func(data interface{}) []MCPDocument {
			vaultData, ok := data.(vault.VaultData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, auth := range vaultData.AuthMethods {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("vault-auth-%s", auth.Path),
					Content: formatVaultAuthInfo(auth),
					Metadata: map[string]interface{}{
						"type":        auth.Type,
						"description": auth.Description,
					},
					Source:     "vault",
					SourceType: "auth_method",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "vault",
		ResourceType: "secret_engine",
		ExtractFunc: func(data interface{}) []MCPDocument {
			vaultData, ok := data.(vault.VaultData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for _, engine := range vaultData.SecretEngines {
				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("vault-secret-engine-%s", engine.Path),
					Content: formatVaultSecretEngineInfo(engine),
					Metadata: map[string]interface{}{
						"type":        engine.Type,
						"description": engine.Description,
					},
					Source:     "vault",
					SourceType: "secret_engine",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	// =====================================================================
	// VEEAM HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "veeam",
		ResourceType: "server_info",
		ExtractFunc: func(data interface{}) []MCPDocument {
			veeamData, ok := data.(veeam.VeeamData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			if veeamData.ServerInfo != nil {
				serverInfoStr := fmt.Sprintf("%v", veeamData.ServerInfo)
				docs = append(docs, MCPDocument{
					ID:      "veeam-server-info",
					Content: fmt.Sprintf("Veeam Backup & Replication Server:\n%s", serverInfoStr),
					Metadata: map[string]interface{}{
						"type": "server_info",
					},
					Source:     "veeam",
					SourceType: "server_info",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "veeam",
		ResourceType: "backup_job",
		ExtractFunc: func(data interface{}) []MCPDocument {
			veeamData, ok := data.(veeam.VeeamData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for i, job := range veeamData.BackupJobs {
				jobID := fmt.Sprintf("%d", i)
				jobStr := fmt.Sprintf("%v", job)

				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("veeam-job-%s", jobID),
					Content: fmt.Sprintf("Veeam Backup Job:\n%s", jobStr),
					Metadata: map[string]interface{}{
						"type": "backup_job",
					},
					Source:     "veeam",
					SourceType: "backup_job",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "veeam",
		ResourceType: "repository",
		ExtractFunc: func(data interface{}) []MCPDocument {
			veeamData, ok := data.(veeam.VeeamData)
			if !ok {
				return nil
			}

			var docs []MCPDocument
			for i, repo := range veeamData.Repositories {
				repoID := fmt.Sprintf("%d", i)
				repoStr := fmt.Sprintf("%v", repo)

				docs = append(docs, MCPDocument{
					ID:      fmt.Sprintf("veeam-repository-%s", repoID),
					Content: fmt.Sprintf("Veeam Repository:\n%s", repoStr),
					Metadata: map[string]interface{}{
						"type": "repository",
					},
					Source:     "veeam",
					SourceType: "repository",
					CreatedAt:  time.Now(),
				})
			}
			return docs
		},
	})

	// =====================================================================
	// SNAPSHOTS HANDLERS
	// =====================================================================
	RegisterResourceHandler(ResourceTypeHandler{
		Platform:     "snapshots",
		ResourceType: "all",
		ExtractFunc: func(data interface{}) []MCPDocument {
			// Snapshots data is a map[string]interface{}
			snapData, ok := data.(map[string]interface{})
			if !ok {
				return nil
			}

			var docs []MCPDocument
			// Process each platform's snapshots
			for platform, platformData := range snapData {
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
						var snapshotMap map[string]string

						// Try to convert to map[string]string
						if sMap, ok := snapshot.(map[string]string); ok {
							snapshotMap = sMap
						} else if sMapGeneric, ok := snapshot.(map[string]interface{}); ok {
							// Convert map[string]interface{} to map[string]string
							snapshotMap = make(map[string]string)
							for k, v := range sMapGeneric {
								snapshotMap[k] = fmt.Sprintf("%v", v)
							}
						} else {
							continue
						}

						// Try to get ID from the snapshot
						snapshotID := fmt.Sprintf("%d", i)
						for _, idField := range []string{"ID", "SnapshotId", "SnapshotID", "Name"} {
							if id, ok := snapshotMap[idField]; ok && id != "" {
								snapshotID = id
								break
							}
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
		},
	})
}
