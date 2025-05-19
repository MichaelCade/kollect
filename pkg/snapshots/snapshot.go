package snapshots

import (
	"context"
	"fmt"
	"log"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/azure"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/kollect"
)

func CollectAllSnapshots(ctx context.Context, kubeconfigPath string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	log.Println("Collecting Kubernetes snapshots...")
	k8sSnapshots, err := kollect.CollectSnapshotData(ctx, kubeconfigPath)
	if err != nil {
		log.Printf("Warning: Error collecting Kubernetes snapshots: %v", err)
	} else if k8sSnapshots != nil {
		log.Println("Successfully collected Kubernetes snapshots")
		results["kubernetes"] = k8sSnapshots
	}

	log.Println("Collecting AWS snapshots...")
	awsSnapshots, err := aws.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Error collecting AWS snapshots: %v", err)
	} else if awsSnapshots != nil {
		log.Println("Successfully collected AWS snapshots")
		results["aws"] = awsSnapshots
	}

	log.Println("Collecting Azure snapshots...")
	azureSnapshots, err := azure.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Error collecting Azure snapshots: %v", err)
	} else if azureSnapshots != nil {
		log.Println("Successfully collected Azure snapshots")
		results["azure"] = azureSnapshots
	}

	log.Println("Collecting GCP snapshots...")
	gcpSnapshots, err := gcp.CollectSnapshotData(ctx)
	if err != nil {
		log.Printf("Warning: Error collecting GCP snapshots: %v", err)
	} else if gcpSnapshots != nil {
		log.Println("Successfully collected GCP snapshots")
		results["gcp"] = gcpSnapshots
	}

	return results, nil
}

func CollectPlatformSnapshots(ctx context.Context, platform string, kubeconfigPath string) (map[string]interface{}, error) {
	switch platform {
	case "kubernetes":
		snapshots, err := kollect.CollectSnapshotData(ctx, kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("error collecting Kubernetes snapshots: %v", err)
		}
		return snapshots, nil
	case "aws":
		snapshots, err := aws.CollectSnapshotData(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting AWS snapshots: %v", err)
		}
		return snapshots, nil
	case "azure":
		snapshots, err := azure.CollectSnapshotData(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting Azure snapshots: %v", err)
		}
		return snapshots, nil
	case "gcp":
		snapshots, err := gcp.CollectSnapshotData(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting GCP snapshots: %v", err)
		}
		return snapshots, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
