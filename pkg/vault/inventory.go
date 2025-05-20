package vault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

type VaultData struct {
	ServerInfo      ServerInfo             `json:"serverInfo"`
	ReplicationInfo ReplicationInfo        `json:"replicationInfo"`
	AuthMethods     []AuthMethodInfo       `json:"authMethods"`
	SecretEngines   []SecretEngineInfo     `json:"secretEngines"`
	Policies        []PolicyInfo           `json:"policies"`
	Namespaces      []NamespaceInfo        `json:"namespaces"`
	AuditDevices    []AuditDeviceInfo      `json:"auditDevices"`
	SecretStats     map[string]interface{} `json:"secretStats"`
	PerformanceInfo map[string]interface{} `json:"performanceInfo"`
	EntityCount     int                    `json:"entityCount"`
	GroupCount      int                    `json:"groupCount"`
	TokenCount      int                    `json:"tokenCount"`
	LicenseInfo     map[string]interface{} `json:"licenseInfo,omitempty"`
}

type ServerInfo struct {
	Version         string    `json:"version"`
	ClusterName     string    `json:"clusterName"`
	ClusterID       string    `json:"clusterId"`
	Initialized     bool      `json:"initialized"`
	Sealed          bool      `json:"sealed"`
	Standby         bool      `json:"standby"`
	HASEnabled      bool      `json:"haEnabled"`
	RaftLeader      bool      `json:"raftLeader,omitempty"`
	StorageType     string    `json:"storageType"`
	LastWALIndex    uint64    `json:"lastWALIndex,omitempty"`
	ServerTimestamp time.Time `json:"serverTimestamp"`
}

type ReplicationInfo struct {
	DREnabled            bool                   `json:"drEnabled"`
	DRMode               string                 `json:"drMode"`
	DRConnected          bool                   `json:"drConnected"`
	PerformanceEnabled   bool                   `json:"performanceEnabled"`
	PerformanceMode      string                 `json:"performanceMode"`
	PerformanceConnected bool                   `json:"performanceConnected"`
	PrimaryClusterAddr   string                 `json:"primaryClusterAddr,omitempty"`
	StateRaw             map[string]interface{} `json:"stateRaw,omitempty"`
}

type AuthMethodInfo struct {
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Accessor    string            `json:"accessor"`
	Local       bool              `json:"local"`
	Config      map[string]string `json:"config,omitempty"`
}

type SecretEngineInfo struct {
	Path        string            `json:"path"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Accessor    string            `json:"accessor"`
	Local       bool              `json:"local"`
	Version     int               `json:"version,omitempty"`
	Options     map[string]string `json:"options,omitempty"`
}

type PolicyInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Rules string `json:"rules,omitempty"`
}

type NamespaceInfo struct {
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	Parent      string `json:"parent,omitempty"`
}

type AuditDeviceInfo struct {
	Path        string                 `json:"path"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
}

func CheckCredentials(ctx context.Context, addr, token string, insecure bool) (bool, error) {
	config := api.DefaultConfig()
	if addr != "" {
		config.Address = addr
	}

	if insecure {
		config.ConfigureTLS(&api.TLSConfig{
			Insecure: true,
		})
	}

	client, err := api.NewClient(config)
	if err != nil {
		return false, fmt.Errorf("failed to create vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	}

	_, err = client.Sys().Health()
	if err != nil {
		return false, fmt.Errorf("failed to access vault: %w", err)
	}

	return true, nil
}

func CollectVaultData(ctx context.Context, addr, token string, insecure bool) (VaultData, error) {
	var data VaultData

	config := api.DefaultConfig()
	if addr != "" {
		config.Address = addr
	}

	if insecure {
		config.ConfigureTLS(&api.TLSConfig{
			Insecure: true,
		})
	}

	client, err := api.NewClient(config)
	if err != nil {
		return data, fmt.Errorf("failed to create vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	}

	health, err := client.Sys().Health()
	if err != nil {
		return data, fmt.Errorf("failed to get vault health: %w", err)
	}

	data.ServerInfo, err = collectServerInfo(client, health)
	if err != nil {
		log.Printf("Warning: Failed to collect full server info: %v", err)
	}

	if health.Sealed {
		return data, nil
	}

	data.ReplicationInfo, err = collectReplicationInfo(client)
	if err != nil {
		log.Printf("Warning: Failed to collect replication info: %v", err)
	}

	data.AuthMethods, err = collectAuthMethods(client)
	if err != nil {
		log.Printf("Warning: Failed to collect auth methods: %v", err)
	}

	data.SecretEngines, err = collectSecretEngines(client)
	if err != nil {
		log.Printf("Warning: Failed to collect secret engines: %v", err)
	}

	data.Policies, err = collectPolicies(client)
	if err != nil {
		log.Printf("Warning: Failed to collect policies: %v", err)
	}

	data.Namespaces, err = collectNamespaces(client)
	if err != nil {
		if !strings.Contains(err.Error(), "not supported") {
			log.Printf("Warning: Failed to collect namespaces: %v", err)
		}
	}

	data.AuditDevices, err = collectAuditDevices(client)
	if err != nil {
		log.Printf("Warning: Failed to collect audit devices: %v", err)
	}

	identityCount, err := collectIdentityCount(client)
	if err != nil {
		log.Printf("Warning: Failed to collect identity counts: %v", err)
	} else {
		data.EntityCount = identityCount.EntityCount
		data.GroupCount = identityCount.GroupCount
	}

	licenseInfo, err := collectLicenseInfo(client)
	if err == nil {
		data.LicenseInfo = licenseInfo
	}

	data.PerformanceInfo, err = collectPerformanceInfo(client)
	if err != nil {
		log.Printf("Warning: Failed to collect performance info: %v", err)
	}

	return data, nil
}

func collectServerInfo(client *api.Client, health *api.HealthResponse) (ServerInfo, error) {
	var info ServerInfo

	info.Version = health.Version
	info.Initialized = health.Initialized
	info.Sealed = health.Sealed
	info.Standby = health.Standby
	info.ServerTimestamp = time.Now()

	status, err := client.Sys().SealStatus()
	if err != nil {
		return info, fmt.Errorf("failed to get seal status: %w", err)
	}

	info.ClusterName = status.ClusterName
	info.ClusterID = status.ClusterID
	info.StorageType = status.Type

	leaderStatus, err := client.Sys().Leader()
	if err == nil {
		info.HASEnabled = leaderStatus.HAEnabled
		info.RaftLeader = leaderStatus.IsSelf
	}

	return info, nil
}
func collectReplicationInfo(client *api.Client) (ReplicationInfo, error) {
	var info ReplicationInfo

	statusRaw, err := client.Logical().Read("sys/replication/status")
	if err != nil {
		return info, fmt.Errorf("failed to get replication status: %w", err)
	}

	if statusRaw == nil || statusRaw.Data == nil {
		return info, nil
	}

	info.StateRaw = statusRaw.Data

	if dr, ok := statusRaw.Data["dr"]; ok && dr != nil {
		if drMap, ok := dr.(map[string]interface{}); ok {
			if mode, ok := drMap["mode"]; ok {
				info.DRMode = fmt.Sprintf("%v", mode)
				info.DREnabled = info.DRMode != "disabled"
			}

			if state, ok := drMap["connection_state"]; ok && state != nil {
				info.DRConnected = fmt.Sprintf("%v", state) == "connected"
			}
		}
	}

	if perf, ok := statusRaw.Data["performance"]; ok && perf != nil {
		if perfMap, ok := perf.(map[string]interface{}); ok {
			if mode, ok := perfMap["mode"]; ok {
				info.PerformanceMode = fmt.Sprintf("%v", mode)
				info.PerformanceEnabled = info.PerformanceMode != "disabled"
			}

			if state, ok := perfMap["connection_state"]; ok && state != nil {
				info.PerformanceConnected = fmt.Sprintf("%v", state) == "connected"
			}

			if addr, ok := perfMap["primary_cluster_addr"]; ok && addr != nil {
				info.PrimaryClusterAddr = fmt.Sprintf("%v", addr)
			}
		}
	}

	return info, nil
}

func collectAuthMethods(client *api.Client) ([]AuthMethodInfo, error) {
	var methods []AuthMethodInfo

	auths, err := client.Sys().ListAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to list auth methods: %w", err)
	}

	for path, auth := range auths {
		method := AuthMethodInfo{
			Path:        strings.TrimSuffix(path, "/"),
			Type:        auth.Type,
			Description: auth.Description,
			Accessor:    auth.Accessor,
			Local:       auth.Local,
			Config:      make(map[string]string),
		}

		configPath := fmt.Sprintf("sys/auth/%s/tune", method.Path)
		config, err := client.Logical().Read(configPath)
		if err == nil && config != nil && config.Data != nil {
			for k, v := range config.Data {
				if k != "options" && v != nil {
					method.Config[k] = fmt.Sprintf("%v", v)
				}
			}
		}

		methods = append(methods, method)
	}

	return methods, nil
}

func collectSecretEngines(client *api.Client) ([]SecretEngineInfo, error) {
	var engines []SecretEngineInfo

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return nil, fmt.Errorf("failed to list secret engines: %w", err)
	}

	for path, mount := range mounts {
		engine := SecretEngineInfo{
			Path:        strings.TrimSuffix(path, "/"),
			Type:        mount.Type,
			Description: mount.Description,
			Accessor:    mount.Accessor,
			Local:       mount.Local,
			Options:     make(map[string]string),
		}

		if mount.Options != nil {
			for k, v := range mount.Options {
				engine.Options[k] = v
			}
		}

		if mount.Type == "kv" {
			if v, ok := mount.Options["version"]; ok {
				if v == "2" {
					engine.Version = 2
				} else {
					engine.Version = 1
				}
			}
		}

		engines = append(engines, engine)
	}

	return engines, nil
}

func collectPolicies(client *api.Client) ([]PolicyInfo, error) {
	var policies []PolicyInfo

	policyNames, err := client.Sys().ListPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to list policies: %w", err)
	}

	for _, name := range policyNames {
		if name == "default" || name == "root" {
			policies = append(policies, PolicyInfo{
				Name: name,
				Type: "acl",
			})
			continue
		}

		policy, err := client.Sys().GetPolicy(name)
		if err != nil {
			log.Printf("Warning: Failed to get policy %s: %v", name, err)
			policies = append(policies, PolicyInfo{
				Name: name,
				Type: "acl",
			})
			continue
		}

		policies = append(policies, PolicyInfo{
			Name:  name,
			Type:  "acl",
			Rules: policy,
		})
	}

	rgpNames, err := client.Logical().List("sys/policies/rgp")
	if err == nil && rgpNames != nil && rgpNames.Data != nil {
		if keys, ok := rgpNames.Data["keys"].([]interface{}); ok {
			for _, k := range keys {
				name := fmt.Sprintf("%v", k)
				rgpPolicy, err := client.Logical().Read(fmt.Sprintf("sys/policies/rgp/%s", name))
				if err != nil || rgpPolicy == nil {
					policies = append(policies, PolicyInfo{
						Name: name,
						Type: "rgp",
					})
					continue
				}

				var rules string
				if policy, ok := rgpPolicy.Data["policy"]; ok && policy != nil {
					rules = fmt.Sprintf("%v", policy)
				}

				policies = append(policies, PolicyInfo{
					Name:  name,
					Type:  "rgp",
					Rules: rules,
				})
			}
		}
	}

	return policies, nil
}

func collectNamespaces(client *api.Client) ([]NamespaceInfo, error) {
	var namespaces []NamespaceInfo

	namespacesRaw, err := client.Logical().List("sys/namespaces")
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	if namespacesRaw == nil || namespacesRaw.Data == nil {
		return namespaces, nil
	}

	if keys, ok := namespacesRaw.Data["keys"].([]interface{}); ok {
		for _, k := range keys {
			path := fmt.Sprintf("%v", k)
			ns := NamespaceInfo{
				Path: strings.TrimSuffix(path, "/"),
			}

			nsInfo, err := client.Logical().Read(fmt.Sprintf("sys/namespaces/%s", path))
			if err == nil && nsInfo != nil && nsInfo.Data != nil {
				if desc, ok := nsInfo.Data["description"]; ok && desc != nil {
					ns.Description = fmt.Sprintf("%v", desc)
				}
				if parent, ok := nsInfo.Data["parent_namespace_path"]; ok && parent != nil {
					ns.Parent = fmt.Sprintf("%v", parent)
				}
			}

			namespaces = append(namespaces, ns)
		}
	}

	return namespaces, nil
}

func collectAuditDevices(client *api.Client) ([]AuditDeviceInfo, error) {
	var devices []AuditDeviceInfo

	auditDevices, err := client.Sys().ListAudit()
	if err != nil {
		return nil, fmt.Errorf("failed to list audit devices: %w", err)
	}

	for path, device := range auditDevices {
		options := make(map[string]interface{})
		for k, v := range device.Options {
			options[k] = v
		}
		auditDevice := AuditDeviceInfo{
			Path:        strings.TrimSuffix(path, "/"),
			Type:        device.Type,
			Description: device.Description,
			Options:     options,
		}

		devices = append(devices, auditDevice)
	}

	return devices, nil
}

type identityCounts struct {
	EntityCount int
	GroupCount  int
}

func collectIdentityCount(client *api.Client) (identityCounts, error) {
	var counts identityCounts

	entityList, err := client.Logical().List("identity/entity/id")
	if err == nil && entityList != nil && entityList.Data != nil {
		if keys, ok := entityList.Data["keys"].([]interface{}); ok {
			counts.EntityCount = len(keys)
		}
	}

	groupList, err := client.Logical().List("identity/group/id")
	if err == nil && groupList != nil && groupList.Data != nil {
		if keys, ok := groupList.Data["keys"].([]interface{}); ok {
			counts.GroupCount = len(keys)
		}
	}

	return counts, nil
}

func collectLicenseInfo(client *api.Client) (map[string]interface{}, error) {
	licenseStatus, err := client.Logical().Read("sys/license/status")
	if err != nil {
		return nil, fmt.Errorf("failed to read license status: %w", err)
	}

	if licenseStatus == nil || licenseStatus.Data == nil {
		return nil, fmt.Errorf("no license data found")
	}

	return licenseStatus.Data, nil
}

func collectPerformanceInfo(client *api.Client) (map[string]interface{}, error) {
	perfInfo := make(map[string]interface{})

	metrics, err := client.Logical().Read("sys/metrics")
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics: %w", err)
	}

	if metrics == nil || metrics.Data == nil {
		return perfInfo, nil
	}

	if gauges, ok := metrics.Data["Gauges"].([]interface{}); ok {
		for _, g := range gauges {
			if gauge, ok := g.(map[string]interface{}); ok {
				name, ok1 := gauge["Name"].(string)
				value, ok2 := gauge["Value"].(float64)
				if ok1 && ok2 {
					if strings.Contains(name, "vault.token.count") ||
						strings.Contains(name, "vault.expire.num_leases") ||
						strings.Contains(name, "vault.raft.") ||
						strings.Contains(name, "vault.runtime.") {
						perfInfo[name] = value
					}
				}
			}
		}
	}

	if counters, ok := metrics.Data["Counters"].([]interface{}); ok {
		for _, c := range counters {
			if counter, ok := c.(map[string]interface{}); ok {
				name, ok1 := counter["Name"].(string)
				count, ok2 := counter["Count"].(float64)
				rate, ok3 := counter["Rate"].(float64)
				if ok1 && (ok2 || ok3) {
					if strings.Contains(name, "vault.route.") ||
						strings.Contains(name, "vault.audit.") ||
						strings.Contains(name, "vault.core.") {
						if ok2 {
							perfInfo[name+".count"] = count
						}
						if ok3 {
							perfInfo[name+".rate"] = rate
						}
					}
				}
			}
		}
	}

	return perfInfo, nil
}
