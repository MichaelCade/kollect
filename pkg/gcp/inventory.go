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
		log.Printf("Warning: %v", err)

	}

	instances, err := fetchComputeInstances(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch compute instances: %v", err)
	} else {
		data.ComputeInstances = instances
	}

	buckets, err := fetchGCSBuckets(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch GCS buckets: %v", err)
	} else {
		data.GCSBuckets = buckets
	}

	sqlInstances, err := fetchCloudSQLInstances(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud SQL instances: %v", err)
	} else {
		data.CloudSQLInstances = sqlInstances
	}

	runServices, err := fetchCloudRunServices(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud Run services: %v", err)
	} else {
		data.CloudRunServices = runServices
	}

	functions, err := fetchCloudFunctions(ctx, projectID)
	if err != nil {
		log.Printf("Warning: Failed to fetch Cloud Functions: %v", err)
	} else {
		data.CloudFunctions = functions
	}

	return data, nil
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
			return nil, fmt.Errorf("failed to list instances in zone %s: %v", zone.Name, err)
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

func CollectSnapshotData(ctx context.Context) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}

	project, err := getCurrentProject()
	if err != nil {
		return nil, fmt.Errorf("failed to get current GCP project: %v", err)
	}

	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	regionsList, err := computeService.Regions.List(project).Do()
	if err != nil {
		log.Printf("Warning: Failed to fetch GCP regions: %v", err)
		return collectSnapshotsFromDefaultRegions(ctx, project)
	}

	var regions []string
	for _, region := range regionsList.Items {
		regions = append(regions, region.Name)
	}

	if len(regions) == 0 {
		log.Printf("No regions found in GCP project, using default regions")
		return collectSnapshotsFromDefaultRegions(ctx, project)
	}

	diskSnapshots, err := collectDiskSnapshots(ctx, project, regions)
	if err != nil {
		log.Printf("Warning: Failed to collect disk snapshots: %v", err)
	} else if len(diskSnapshots) > 0 {
		snapshots["DiskSnapshots"] = diskSnapshots
	}

	return snapshots, nil
}

func collectSnapshotsFromDefaultRegions(ctx context.Context, project string) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}
	defaultRegions := []string{
		"us-central1", "us-east1", "us-west1", "us-west2", "us-east4",
		"europe-west1", "europe-west2", "europe-west3", "europe-west4",
		"asia-east1", "asia-southeast1", "asia-northeast1",
		"australia-southeast1",
	}

	diskSnapshots, err := collectDiskSnapshots(ctx, project, defaultRegions)
	if err != nil {
		log.Printf("Warning: Failed to collect disk snapshots: %v", err)
	} else if len(diskSnapshots) > 0 {
		snapshots["DiskSnapshots"] = diskSnapshots
	}

	return snapshots, nil
}

func collectDiskSnapshots(ctx context.Context, project string, regions []string) ([]map[string]string, error) {
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	snapshotList, err := computeService.Snapshots.List(project).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshots: %v", err)
	}

	diskRegionMap := make(map[string]string)

	zoneToRegion := make(map[string]string)
	for _, region := range regions {
		zoneList, err := computeService.Zones.List(project).Filter(fmt.Sprintf("region eq .*/regions/%s", region)).Do()
		if err != nil {
			log.Printf("Warning: Failed to list zones for region %s: %v", region, err)
			continue
		}

		for _, zone := range zoneList.Items {
			zoneToRegion[zone.Name] = region
		}
	}

	for zone, region := range zoneToRegion {
		diskList, err := computeService.Disks.List(project, zone).Do()
		if err != nil {
			log.Printf("Warning: Failed to list disks in zone %s: %v", zone, err)
			continue
		}

		for _, disk := range diskList.Items {
			diskName := disk.Name
			if diskName != "" {
				diskRegionMap[diskName] = region
			}
		}
	}

	var snapshots []map[string]string
	defaultRegion := regions[0]

	for _, snapshot := range snapshotList.Items {
		region := defaultRegion
		if snapshot.SourceDisk != "" {
			parts := strings.Split(snapshot.SourceDisk, "/")
			if len(parts) > 0 {
				diskName := parts[len(parts)-1]
				if r, ok := diskRegionMap[diskName]; ok {
					region = r
				} else {
					for _, r := range regions {
						if strings.Contains(snapshot.SourceDisk, r) {
							region = r
							break
						}
					}
				}
			}
		}

		snapshotInfo := map[string]string{
			"Name":    snapshot.Name,
			"Region":  region,
			"Status":  snapshot.Status,
			"Project": project,
		}

		if snapshot.DiskSizeGb > 0 {
			snapshotInfo["DiskSizeGB"] = fmt.Sprintf("%d", snapshot.DiskSizeGb)
		}

		if snapshot.CreationTimestamp != "" {
			snapshotInfo["CreationTime"] = snapshot.CreationTimestamp
		}

		if snapshot.SourceDisk != "" {
			snapshotInfo["SourceDisk"] = snapshot.SourceDisk
		}

		if snapshot.StorageBytes > 0 {
			snapshotInfo["StorageBytes"] = fmt.Sprintf("%d", snapshot.StorageBytes)
		}

		snapshots = append(snapshots, snapshotInfo)
	}

	return snapshots, nil
}

func getGCPProjectID() (string, error) {
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get project ID from gcloud: %v", err)
	}

	projectID := strings.TrimSpace(string(output))
	if projectID == "" {
		return "", fmt.Errorf("no project ID configured in gcloud")
	}

	return projectID, nil
}

func GetProjectID() string {
	projectID, err := getCurrentProject()
	if err != nil {
		log.Printf("Warning: %v", err)
		return "demo-project"
	}
	return projectID
}
