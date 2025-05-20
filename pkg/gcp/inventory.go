package gcp

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/run/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

type ComputeInstanceInfo struct {
	Name        string
	Zone        string
	MachineType string
	Status      string
	Project     string
}

type GCSBucketInfo struct {
	Name              string
	Location          string
	StorageClass      string
	RetentionPolicy   bool
	RetentionDuration int64
	Project           string
}

type CloudSQLInstanceInfo struct {
	Name            string
	DatabaseVersion string
	Region          string
	Tier            string
	Status          string
	Project         string
}

type CloudRunServiceInfo struct {
	Name      string
	Region    string
	URL       string
	Project   string
	Replicas  int64
	Container string
}

type CloudFunctionInfo struct {
	Name            string
	Region          string
	Runtime         string
	Status          string
	EntryPoint      string
	AvailableMemory string
	Project         string
}

type GCPData struct {
	ComputeInstances  []ComputeInstanceInfo
	GCSBuckets        []GCSBucketInfo
	CloudSQLInstances []CloudSQLInstanceInfo
	CloudRunServices  []CloudRunServiceInfo
	CloudFunctions    []CloudFunctionInfo
}

func CheckCredentials(ctx context.Context) (bool, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, err
	}
	defer client.Close()

	_, err = getCurrentProject()

	return err == nil, err
}

func CollectGCPData(ctx context.Context) (GCPData, error) {
	var data GCPData

	projectID, err := getCurrentProject()
	if err != nil {
		log.Printf("Warning: %v", formatAPIError(err))
	}

	instances, err := fetchComputeInstances(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch compute instances: %v", formatAPIError(err))
	} else {
		data.ComputeInstances = instances
	}

	buckets, err := fetchGCSBuckets(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch GCS buckets: %v", formatAPIError(err))
	} else {
		data.GCSBuckets = buckets
	}

	sqlInstances, err := fetchCloudSQLInstances(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud SQL instances: %v", formatAPIError(err))
	} else {
		data.CloudSQLInstances = sqlInstances
	}

	runServices, err := fetchCloudRunServices(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud Run services: %v", formatAPIError(err))
	} else {
		data.CloudRunServices = runServices
	}

	functions, err := fetchCloudFunctions(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud Functions: %v", formatAPIError(err))
	} else {
		data.CloudFunctions = functions
	}

	return data, nil
}

// formatAPIError shortens Google API errors to just the essential information
func formatAPIError(err error) string {
	errorStr := err.Error()

	// Check if this is an API not enabled error
	if strings.Contains(errorStr, "has not been used in project") && strings.Contains(errorStr, "or it is disabled") {
		parts := strings.Split(errorStr, ":")
		if len(parts) >= 3 {
			apiName := strings.TrimSpace(parts[2])
			if strings.Contains(apiName, "API") {
				// Find the API name (e.g., "Cloud SQL Admin API")
				apiNameEnd := strings.Index(apiName, "has not been used")
				if apiNameEnd > 0 {
					apiName = strings.TrimSpace(apiName[:apiNameEnd])
					return fmt.Sprintf("%s is not enabled for this project", apiName)
				}
			}
		}
		return "API not enabled for this project"
	}

	// If it's another kind of error or we couldn't parse properly
	if len(errorStr) > 150 {
		return errorStr[:150] + "... (truncated)"
	}

	return errorStr
}

func getCurrentProject() (string, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID != "" {
		return projectID, nil
	}

	cmd := exec.Command("gcloud", "config", "get-value", "project")
	output, err := cmd.Output()
	if err == nil {
		projectID = strings.TrimSpace(string(output))
		if projectID != "" {
			return projectID, nil
		}
	}

	return "demo-project", fmt.Errorf("could not determine GCP project ID, using demo-project")
}

func fetchComputeInstances(ctx context.Context, projectID string) ([]ComputeInstanceInfo, error) {
	var instances []ComputeInstanceInfo

	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	zonesResp, err := computeService.Zones.List(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list zones: %v", err)
	}

	for _, zone := range zonesResp.Items {
		instancesResp, err := computeService.Instances.List(projectID, zone.Name).Do()
		if err != nil {
			continue // Skip zones with errors rather than failing completely
		}

		for _, instance := range instancesResp.Items {
			machineType := instance.MachineType
			if parts := strings.Split(machineType, "/"); len(parts) > 0 {
				machineType = parts[len(parts)-1]
			}

			instances = append(instances, ComputeInstanceInfo{
				Name:        instance.Name,
				Zone:        zone.Name,
				MachineType: machineType,
				Status:      instance.Status,
				Project:     projectID,
			})
		}
	}

	return instances, nil
}

func fetchGCSBuckets(ctx context.Context, projectID string) ([]GCSBucketInfo, error) {
	var buckets []GCSBucketInfo

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	it := client.Buckets(ctx, projectID)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error iterating buckets: %v", err)
		}

		bucketInfo := GCSBucketInfo{
			Name:         bucketAttrs.Name,
			Location:     bucketAttrs.Location,
			StorageClass: bucketAttrs.StorageClass,
			Project:      projectID,
		}

		if bucketAttrs.RetentionPolicy != nil {
			bucketInfo.RetentionPolicy = true
			bucketInfo.RetentionDuration = int64(bucketAttrs.RetentionPolicy.RetentionPeriod.Seconds())
		}

		buckets = append(buckets, bucketInfo)
	}

	return buckets, nil
}

func fetchCloudSQLInstances(ctx context.Context, projectID string) ([]CloudSQLInstanceInfo, error) {
	var instances []CloudSQLInstanceInfo

	sqlService, err := sqladmin.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud SQL service: %v", err)
	}

	instancesResp, err := sqlService.Instances.List(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list Cloud SQL instances: %v", err)
	}

	for _, instance := range instancesResp.Items {
		instances = append(instances, CloudSQLInstanceInfo{
			Name:            instance.Name,
			DatabaseVersion: instance.DatabaseVersion,
			Region:          instance.Region,
			Tier:            instance.Settings.Tier,
			Status:          instance.State,
			Project:         projectID,
		})
	}

	return instances, nil
}

func fetchCloudRunServices(ctx context.Context, projectID string) ([]CloudRunServiceInfo, error) {
	var services []CloudRunServiceInfo

	runService, err := run.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Run service: %v", err)
	}

	locationsResp, err := runService.Projects.Locations.List("projects/" + projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list Cloud Run locations: %v", err)
	}

	for _, location := range locationsResp.Locations {
		servicesResp, err := runService.Projects.Locations.Services.List(
			fmt.Sprintf("projects/%s/locations/%s", projectID, location.LocationId)).Do()
		if err != nil {
			continue
		}

		for _, service := range servicesResp.Items {
			var replicas int64
			var container string

			if service.Spec != nil && len(service.Spec.Template.Spec.Containers) > 0 {
				container = service.Spec.Template.Spec.Containers[0].Image
			}

			if service.Status != nil && service.Status.Traffic != nil && len(service.Status.Traffic) > 0 {
				replicas = service.Status.Traffic[0].Percent
			}

			services = append(services, CloudRunServiceInfo{
				Name:      service.Metadata.Name,
				Region:    location.LocationId,
				URL:       service.Status.Url,
				Project:   projectID,
				Replicas:  replicas,
				Container: container,
			})
		}
	}

	return services, nil
}

func fetchCloudFunctions(ctx context.Context, projectID string) ([]CloudFunctionInfo, error) {
	var functions []CloudFunctionInfo

	functionsService, err := cloudfunctions.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Functions service: %v", err)
	}

	locationsCall := functionsService.Projects.Locations.List("projects/" + projectID)
	locationsResp, err := locationsCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list Cloud Functions locations: %v", err)
	}

	for _, location := range locationsResp.Locations {
		functionsCall := functionsService.Projects.Locations.Functions.List(
			fmt.Sprintf("projects/%s/locations/%s", projectID, location.LocationId))
		functionsResp, err := functionsCall.Do()
		if err != nil {
			continue
		}

		for _, function := range functionsResp.Functions {
			functions = append(functions, CloudFunctionInfo{
				Name:            function.Name,
				Region:          location.LocationId,
				Runtime:         function.Runtime,
				Status:          function.Status,
				EntryPoint:      function.EntryPoint,
				AvailableMemory: fmt.Sprintf("%dMB", function.AvailableMemoryMb),
				Project:         projectID,
			})
		}
	}

	return functions, nil
}

// CollectSnapshotData collects all disk snapshots in a GCP project
func CollectSnapshotData(ctx context.Context) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}

	project, err := getCurrentProject()
	if err != nil {
		return nil, fmt.Errorf("failed to get current GCP project: %v", formatAPIError(err))
	}

	diskSnapshots, err := collectDiskSnapshots(ctx, project)
	if err != nil {
		log.Printf("Warning: Failed to collect disk snapshots: %v", formatAPIError(err))
	} else if len(diskSnapshots) > 0 {
		snapshots["DiskSnapshots"] = diskSnapshots
		log.Printf("Found %d GCP disk snapshots", len(diskSnapshots))
	} else {
		log.Printf("No GCP disk snapshots found")
	}

	return snapshots, nil
}

// collectDiskSnapshots gets all disk snapshots in a project
func collectDiskSnapshots(ctx context.Context, project string) ([]map[string]string, error) {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	// Get all snapshots (this returns all snapshots regardless of region)
	snapshotList, err := computeService.Snapshots.List(project).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %v", err)
	}

	log.Printf("Found %d GCP snapshots", len(snapshotList.Items))

	// Map to store disk details
	diskDetails := make(map[string]string)

	// Create snapshot info objects
	var snapshots []map[string]string
	for _, snapshot := range snapshotList.Items {
		location := "global" // Default location

		// Try to determine location from source disk if possible
		if snapshot.SourceDisk != "" {
			// If we don't already have this disk's details, try to get them
			diskName := getDiskNameFromURL(snapshot.SourceDisk)
			if _, exists := diskDetails[diskName]; !exists {
				// Extract zone from disk URL if possible
				zone := getZoneFromDiskURL(snapshot.SourceDisk)
				if zone != "" {
					// If we extracted a zone, try to get the region from it
					if strings.Count(zone, "-") >= 2 {
						// Most zones follow the pattern: region-letter, e.g., us-central1-a
						regionParts := strings.Split(zone, "-")
						if len(regionParts) >= 3 {
							// Remove the last segment (the zone letter)
							location = strings.Join(regionParts[:len(regionParts)-1], "-")
							diskDetails[diskName] = location
						}
					}
				}
			} else {
				location = diskDetails[diskName]
			}
		}

		snapshotInfo := map[string]string{
			"Name":     snapshot.Name,
			"Location": location,
			"Status":   snapshot.Status,
			"Project":  project,
		}

		if snapshot.DiskSizeGb > 0 {
			snapshotInfo["DiskSizeGB"] = fmt.Sprintf("%d", snapshot.DiskSizeGb)
		}

		if snapshot.CreationTimestamp != "" {
			snapshotInfo["CreationTime"] = snapshot.CreationTimestamp
		}

		if snapshot.SourceDisk != "" {
			snapshotInfo["SourceDisk"] = getDiskNameFromURL(snapshot.SourceDisk)
		}

		if snapshot.StorageBytes > 0 {
			snapshotInfo["StorageBytes"] = fmt.Sprintf("%d", snapshot.StorageBytes)
		}

		snapshots = append(snapshots, snapshotInfo)
	}

	return snapshots, nil
}

// Helper function to extract disk name from a URL
func getDiskNameFromURL(diskURL string) string {
	parts := strings.Split(diskURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// Helper function to extract zone from a disk URL
func getZoneFromDiskURL(diskURL string) string {
	// URL format is typically like:
	// https://www.googleapis.com/compute/v1/projects/[PROJECT]/zones/[ZONE]/disks/[DISK]
	parts := strings.Split(diskURL, "/")
	for i, part := range parts {
		if part == "zones" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func GetProjectID() string {
	projectID, err := getCurrentProject()
	if err != nil {
		log.Printf("Warning: %v", formatAPIError(err))
		return "demo-project"
	}
	return projectID
}
