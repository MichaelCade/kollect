package veeam

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type VeeamData struct {
	ServerInfo           map[string]interface{}
	Credentials          []interface{}
	CloudCredentials     []interface{}
	KMSServers           []interface{}
	ManagedServers       []interface{}
	Repositories         []interface{}
	ScaleOutRepositories []interface{}
	Proxies              []interface{}
	BackupJobs           []map[string]interface{}
}

func CollectVeeamData(ctx context.Context, baseURL, username, password string) (VeeamData, error) {
	var data VeeamData

	token, err := authenticate(baseURL, username, password)
	if err != nil {
		return data, fmt.Errorf("authentication failed: %v", err)
	}

	data.ServerInfo, err = getServerInfo(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get server info: %v", err)
	}

	data.Credentials, err = getCredentials(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get credentials: %v", err)
	}

	data.CloudCredentials, err = getCloudCredentials(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get cloud credentials: %v", err)
	}

	data.KMSServers, err = getKMSServers(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get KMS servers: %v", err)
	}

	data.ManagedServers, err = getManagedServers(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get managed servers: %v", err)
	}

	data.Repositories, err = getRepositories(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get repositories: %v", err)
	}

	data.ScaleOutRepositories, err = getScaleOutRepositories(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get scale-out repositories: %v", err)
	}

	data.Proxies, err = getProxies(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to get proxies: %v", err)
	}

	data.BackupJobs, err = listBackupJobs(baseURL, token)
	if err != nil {
		return data, fmt.Errorf("failed to list backup jobs: %v", err)
	}

	return data, nil
}

func authenticate(baseURL, username, password string) (string, error) {
	authURL := fmt.Sprintf("%s/api/oauth2/token", baseURL)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-api-version", "1.1-rev2")
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to authenticate: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return token, nil
}

func getServerInfo(baseURL, token string) (map[string]interface{}, error) {
	return getAPIData(fmt.Sprintf("%s/api/v1/serverInfo", baseURL), token)
}

func getCredentials(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/credentials", baseURL), token)
}

func getCloudCredentials(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/cloudCredentials", baseURL), token)
}

func getKMSServers(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/kmsServers", baseURL), token)
}

func getManagedServers(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/backupInfrastructure/managedServers", baseURL), token)
}

func getRepositories(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/backupInfrastructure/repositories", baseURL), token)
}

func getScaleOutRepositories(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/backupInfrastructure/scaleOutRepositories", baseURL), token)
}

func getProxies(baseURL, token string) ([]interface{}, error) {
	return getAPIList(fmt.Sprintf("%s/api/v1/backupInfrastructure/proxies", baseURL), token)
}

func listBackupJobs(baseURL, token string) ([]map[string]interface{}, error) {
	jobsURL := fmt.Sprintf("%s/api/v1/jobs", baseURL)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", jobsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-version", "1.1-rev2")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get jobs: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	jobs, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("job data not found in response")
	}

	jobList := make([]map[string]interface{}, 0)
	for _, job := range jobs {
		if jobMap, ok := job.(map[string]interface{}); ok {
			jobList = append(jobList, jobMap)
		}
	}

	return jobList, nil
}

func getAPIData(url, token string) (map[string]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-version", "1.1-rev2")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func getAPIList(url, token string) ([]interface{}, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-api-version", "1.1-rev2")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get data: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("data not found in response")
	}

	return data, nil
}
