package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"cloud.google.com/go/storage"
)

type TerraformData struct {
	Resources []ResourceInfo `json:"Resources"`
	Outputs   []OutputInfo   `json:"Outputs"`
	Providers []ProviderInfo `json:"Providers"`
}
type ResourceInfo struct {
	Name         string            `json:"Name"`
	Type         string            `json:"Type"`
	Provider     string            `json:"Provider"`
	Module       string            `json:"Module"`
	Mode         string            `json:"Mode"`
	Attributes   map[string]string `json:"Attributes"`
	Dependencies []string          `json:"Dependencies"`
	Status       string            `json:"Status"`
}
type OutputInfo struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
	Type  string `json:"Type"`
}
type ProviderInfo struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

func CollectTerraformData(ctx context.Context, stateFile string) (TerraformData, error) {
	var data TerraformData

	fileContents, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return data, fmt.Errorf("error reading Terraform state file: %v", err)
	}

	var rawState map[string]interface{}
	if err := json.Unmarshal(fileContents, &rawState); err != nil {
		return data, fmt.Errorf("error parsing Terraform state file: %v", err)
	}

	version, ok := rawState["version"]
	if !ok {
		return data, fmt.Errorf("invalid terraform state file format: missing version")
	}

	versionFloat, ok := version.(float64)
	if !ok || versionFloat < 3 {
		return data, fmt.Errorf("unsupported terraform state file version: %v", version)
	}

	resources, providers, outputs, err := parseStateFile(rawState)
	if err != nil {
		return data, err
	}

	data.Resources = resources
	data.Providers = providers
	data.Outputs = outputs

	return data, nil
}

func CollectTerraformDataFromS3(ctx context.Context, bucket, key, region string) (TerraformData, error) {
	var data TerraformData

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return data, fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return data, fmt.Errorf("failed to retrieve state file from S3: %v", err)
	}
	defer result.Body.Close()

	stateBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return data, fmt.Errorf("failed to read state file content from S3: %v", err)
	}

	var rawState map[string]interface{}
	if err := json.Unmarshal(stateBytes, &rawState); err != nil {
		return data, fmt.Errorf("failed to parse state file from S3: %v", err)
	}

	version, ok := rawState["version"]
	if !ok {
		return data, fmt.Errorf("invalid terraform state file format: missing version")
	}

	versionFloat, ok := version.(float64)
	if !ok || versionFloat < 3 {
		return data, fmt.Errorf("unsupported terraform state file version: %v", version)
	}

	resources, providers, outputs, err := parseStateFile(rawState)
	if err != nil {
		return data, err
	}

	data.Resources = resources
	data.Providers = providers
	data.Outputs = outputs

	return data, nil
}

func CollectTerraformDataFromAzure(ctx context.Context, storageAccount, container, blob string) (TerraformData, error) {
	var data TerraformData

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return data, fmt.Errorf("failed to create Azure credential: %v", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", storageAccount)
	client, err := azblob.NewClient(serviceURL, credential, nil)
	if err != nil {
		return data, fmt.Errorf("failed to create Azure blob client: %v", err)
	}

	downloadResponse, err := client.DownloadStream(ctx, container, blob, nil)
	if err != nil {
		return data, fmt.Errorf("failed to download blob: %v", err)
	}

	stateBytes, err := io.ReadAll(downloadResponse.Body)
	if err != nil {
		return data, fmt.Errorf("failed to read blob content: %v", err)
	}

	var rawState map[string]interface{}
	if err := json.Unmarshal(stateBytes, &rawState); err != nil {
		return data, fmt.Errorf("failed to parse state file from Azure blob: %v", err)
	}

	version, ok := rawState["version"]
	if !ok {
		return data, fmt.Errorf("invalid terraform state file format: missing version")
	}

	versionFloat, ok := version.(float64)
	if !ok || versionFloat < 3 {
		return data, fmt.Errorf("unsupported terraform state file version: %v", version)
	}

	resources, providers, outputs, err := parseStateFile(rawState)
	if err != nil {
		return data, err
	}

	data.Resources = resources
	data.Providers = providers
	data.Outputs = outputs

	return data, nil
}

func CollectTerraformDataFromGCS(ctx context.Context, bucket, object string) (TerraformData, error) {
	var data TerraformData

	client, err := storage.NewClient(ctx)
	if err != nil {
		return data, fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	reader, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return data, fmt.Errorf("failed to read GCS object: %v", err)
	}
	defer reader.Close()

	stateBytes, err := io.ReadAll(reader)
	if err != nil {
		return data, fmt.Errorf("failed to read object content: %v", err)
	}

	var rawState map[string]interface{}
	if err := json.Unmarshal(stateBytes, &rawState); err != nil {
		return data, fmt.Errorf("failed to parse state file from GCS: %v", err)
	}

	version, ok := rawState["version"]
	if !ok {
		return data, fmt.Errorf("invalid terraform state file format: missing version")
	}

	versionFloat, ok := version.(float64)
	if !ok || versionFloat < 3 {
		return data, fmt.Errorf("unsupported terraform state file version: %v", version)
	}

	resources, providers, outputs, err := parseStateFile(rawState)
	if err != nil {
		return data, err
	}

	data.Resources = resources
	data.Providers = providers
	data.Outputs = outputs

	return data, nil
}

func parseStateFile(rawState map[string]interface{}) ([]ResourceInfo, []ProviderInfo, []OutputInfo, error) {
	var resources []ResourceInfo
	var providers []ProviderInfo
	var outputs []OutputInfo

	if resourcesRaw, ok := rawState["resources"]; ok {
		resourcesList, ok := resourcesRaw.([]interface{})
		if ok {
			for _, resRaw := range resourcesList {
				res, ok := resRaw.(map[string]interface{})
				if !ok {
					continue
				}

				mode, _ := res["mode"].(string)
				rType, _ := res["type"].(string)
				name, _ := res["name"].(string)
				provider, _ := res["provider"].(string)
				module, _ := res["module"].(string)

				if instances, ok := res["instances"].([]interface{}); ok {
					for _, instRaw := range instances {
						inst, ok := instRaw.(map[string]interface{})
						if !ok {
							continue
						}

						attributes := make(map[string]string)
						if attrs, ok := inst["attributes"].(map[string]interface{}); ok {
							for k, v := range attrs {
								switch val := v.(type) {
								case string:
									attributes[k] = val
								case float64:
									attributes[k] = fmt.Sprintf("%v", val)
								case bool:
									attributes[k] = fmt.Sprintf("%v", val)
								default:
									attributes[k] = fmt.Sprintf("[%T]", val)
								}
							}
						}

						var dependencies []string
						if deps, ok := inst["dependencies"].([]interface{}); ok {
							for _, dep := range deps {
								if depStr, ok := dep.(string); ok {
									dependencies = append(dependencies, depStr)
								}
							}
						}

						status := "Created"
						if _, hasChanges := inst["changes"]; hasChanges {
							status = "Pending Changes"
						}

						resourceInfo := ResourceInfo{
							Name:         name,
							Type:         rType,
							Provider:     provider,
							Module:       module,
							Mode:         mode,
							Attributes:   attributes,
							Dependencies: dependencies,
							Status:       status,
						}
						resources = append(resources, resourceInfo)
					}
				}
			}
		}
	}

	if providersRaw, ok := rawState["provider_hash"]; ok {
		providerHash, ok := providersRaw.(map[string]interface{})
		if ok {
			for provider, versionRaw := range providerHash {
				version := "unknown"
				if vStr, ok := versionRaw.(string); ok {
					version = vStr
				}

				name := provider
				if strings.HasPrefix(name, "provider.") {
					name = strings.TrimPrefix(name, "provider.")
				}
				if strings.HasPrefix(name, "provider[") {
					name = strings.TrimPrefix(name, "provider[")
					if idx := strings.Index(name, "]"); idx > 0 {
						name = name[:idx]
					}
				}

				providers = append(providers, ProviderInfo{
					Name:    name,
					Version: version,
				})
			}
		}
	}

	if outputsRaw, ok := rawState["outputs"]; ok {
		outputsMap, ok := outputsRaw.(map[string]interface{})
		if ok {
			for name, outRaw := range outputsMap {
				out, ok := outRaw.(map[string]interface{})
				if !ok {
					continue
				}

				value := "complex value"
				if val, ok := out["value"]; ok {
					switch v := val.(type) {
					case string:
						value = v
					case float64:
						value = fmt.Sprintf("%v", v)
					case bool:
						value = fmt.Sprintf("%v", v)
					default:
						value = fmt.Sprintf("[%T]", v)
					}
				}

				typeStr := "string"
				if typ, ok := out["type"]; ok {
					if ts, ok := typ.(string); ok {
						typeStr = ts
					}
				}

				outputs = append(outputs, OutputInfo{
					Name:  name,
					Value: value,
					Type:  typeStr,
				})
			}
		}
	}

	return resources, providers, outputs, nil
}
