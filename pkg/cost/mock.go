package cost

func GenerateMockSnapshotData(platform string) map[string]interface{} {
	switch platform {
	case "aws":
		return map[string]interface{}{
			"EBSSnapshots": []map[string]string{
				{
					"SnapshotId": "snap-0abc123456789def0",
					"VolumeId":   "vol-0abc123456789def0",
					"State":      "completed",
					"VolumeSize": "100",
					"StartTime":  "2023-05-15T00:00:00Z",
				},
				{
					"SnapshotId": "snap-0def987654321abc0",
					"VolumeId":   "vol-0def987654321abc0",
					"State":      "completed",
					"VolumeSize": "250",
					"StartTime":  "2023-05-10T00:00:00Z",
				},
			},
			"RDSSnapshots": []map[string]string{
				{
					"SnapshotId":         "rds:database-1-snapshot-2023-05-01",
					"Engine":             "postgres",
					"AllocatedStorage":   "500",
					"Status":             "available",
					"SnapshotCreateTime": "2023-05-01T00:00:00Z",
				},
			},
		}
	case "azure":
		return map[string]interface{}{
			"DiskSnapshots": []map[string]string{
				{
					"Name":              "snapshot-vm1-osdisk-20230501",
					"DiskSizeGB":        "128",
					"Location":          "eastus",
					"ProvisioningState": "Succeeded",
					"TimeCreated":       "2023-05-01T00:00:00Z",
				},
				{
					"Name":              "snapshot-vm2-datadisk-20230515",
					"DiskSizeGB":        "256",
					"Location":          "westeurope",
					"ProvisioningState": "Succeeded",
					"TimeCreated":       "2023-05-15T00:00:00Z",
				},
			},
		}
	case "gcp":
		return map[string]interface{}{
			"DiskSnapshots": []map[string]string{
				{
					"Name":              "snapshot-instance1-boot-disk",
					"DiskSizeGB":        "50",
					"Location":          "us-central1",
					"Status":            "READY",
					"CreationTimestamp": "2023-05-01T00:00:00Z",
				},
				{
					"Name":              "snapshot-instance2-data-disk",
					"DiskSizeGB":        "200",
					"Location":          "europe-west1",
					"Status":            "READY",
					"CreationTimestamp": "2023-05-15T00:00:00Z",
				},
			},
		}
	default:
		return map[string]interface{}{}
	}
}
