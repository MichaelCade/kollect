package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// TerraformData represents the parsed Terraform state data
type TerraformData struct {
	Resources []ResourceInfo `json:"Resources"`
	Outputs   []OutputInfo   `json:"Outputs"`
	Providers []ProviderInfo `json:"Providers"`
}

// ResourceInfo represents a Terraform resource
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

// OutputInfo represents a Terraform output
type OutputInfo struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
	Type  string `json:"Type"`
}

// ProviderInfo represents a Terraform provider
type ProviderInfo struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

// CollectTerraformData parses a Terraform state file
func CollectTerraformData(ctx context.Context, stateFile string) (TerraformData, error) {
	var data TerraformData

	// Read the state file
	fileContents, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return data, fmt.Errorf("error reading Terraform state file: %v", err)
	}

	// Parse the state file
	var rawState map[string]interface{}
	if err := json.Unmarshal(fileContents, &rawState); err != nil {
		return data, fmt.Errorf("error parsing Terraform state file: %v", err)
	}

	// Extract the version to ensure we're dealing with a state file
	version, ok := rawState["version"]
	if !ok {
		return data, fmt.Errorf("invalid terraform state file format: missing version")
	}

	// Check that we're dealing with a supported version
	// Most code handles version 4, but you might need specific handling for other versions
	versionFloat, ok := version.(float64)
	if !ok || versionFloat < 3 {
		return data, fmt.Errorf("unsupported terraform state file version: %v", version)
	}

	// Handle version 4 state files (most common)
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

	// Extract resources
	if resourcesRaw, ok := rawState["resources"]; ok {
		resourcesList, ok := resourcesRaw.([]interface{})
		if ok {
			for _, resRaw := range resourcesList {
				res, ok := resRaw.(map[string]interface{})
				if !ok {
					continue
				}

				// Extract resource info
				mode, _ := res["mode"].(string)
				rType, _ := res["type"].(string)
				name, _ := res["name"].(string)
				provider, _ := res["provider"].(string)
				module, _ := res["module"].(string)

				// Process instances
				if instances, ok := res["instances"].([]interface{}); ok {
					for _, instRaw := range instances {
						inst, ok := instRaw.(map[string]interface{})
						if !ok {
							continue
						}

						// Extract attributes
						attributes := make(map[string]string)
						if attrs, ok := inst["attributes"].(map[string]interface{}); ok {
							for k, v := range attrs {
								// Skip complex structures, only include simple values
								switch val := v.(type) {
								case string:
									attributes[k] = val
								case float64:
									attributes[k] = fmt.Sprintf("%v", val)
								case bool:
									attributes[k] = fmt.Sprintf("%v", val)
								default:
									// For complex types, just indicate type
									attributes[k] = fmt.Sprintf("[%T]", val)
								}
							}
						}

						// Extract dependencies
						var dependencies []string
						if deps, ok := inst["dependencies"].([]interface{}); ok {
							for _, dep := range deps {
								if depStr, ok := dep.(string); ok {
									dependencies = append(dependencies, depStr)
								}
							}
						}

						// Determine status
						status := "Created"
						// Check for pending states or other indicators
						if _, hasChanges := inst["changes"]; hasChanges {
							status = "Pending Changes"
						}

						// Create the resource info
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

	// Extract providers
	if providersRaw, ok := rawState["provider_hash"]; ok {
		providerHash, ok := providersRaw.(map[string]interface{})
		if ok {
			for provider, versionRaw := range providerHash {
				version := "unknown"
				if vStr, ok := versionRaw.(string); ok {
					version = vStr
				}

				// Remove provider prefix if present
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

	// Extract outputs
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
						// For complex values, just indicate it's complex
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
