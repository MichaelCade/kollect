package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type EC2InstanceInfo struct {
	Name       string
	InstanceID string
	Type       string
	State      string
	Region     string
}

type S3BucketInfo struct {
	Name      string
	Immutable bool
	Region    string
}

type RDSInstanceInfo struct {
	InstanceID string
	Engine     string
	Status     string
	Region     string
}

type DynamoDBTableInfo struct {
	TableName string
	Status    string
	Region    string
}

type VPCInfo struct {
	VPCID  string
	State  string
	Region string
}

type AWSData struct {
	EC2Instances   []EC2InstanceInfo
	S3Buckets      []S3BucketInfo
	RDSInstances   []RDSInstanceInfo
	DynamoDBTables []DynamoDBTableInfo
	VPCs           []VPCInfo
}

func CheckCredentials(ctx context.Context) (bool, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return false, err
	}

	ec2Client := ec2.NewFromConfig(cfg)
	_, err = ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})

	return err == nil, err
}

func FetchEC2Instances(ctx context.Context) ([]EC2InstanceInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	regionsOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe regions, %v", err)
	}

	var allInstances []EC2InstanceInfo
	for _, region := range regionsOutput.Regions {
		regionCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region.RegionName))
		if err != nil {
			return nil, fmt.Errorf("unable to load SDK config for region %s, %v", *region.RegionName, err)
		}

		regionEc2Client := ec2.NewFromConfig(regionCfg)
		result, err := regionEc2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
		if err != nil {
			return nil, fmt.Errorf("unable to describe instances in region %s, %v", *region.RegionName, err)
		}

		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				allInstances = append(allInstances, EC2InstanceInfo{
					Name:       aws.ToString(instance.KeyName),
					InstanceID: aws.ToString(instance.InstanceId),
					Type:       string(instance.InstanceType),
					State:      string(instance.State.Name),
					Region:     *region.RegionName,
				})
			}
		}
	}

	return allInstances, nil
}

func FetchS3Buckets(ctx context.Context) ([]S3BucketInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	result, err := svc.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list buckets, %v", err)
	}

	var buckets []S3BucketInfo
	for _, bucket := range result.Buckets {
		region, err := svc.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get bucket location for %s, %v", aws.ToString(bucket.Name), err)
		}
		objectLockConfig, err := svc.GetObjectLockConfiguration(ctx, &s3.GetObjectLockConfigurationInput{
			Bucket: bucket.Name,
		})
		immutable := false
		if err == nil && objectLockConfig.ObjectLockConfiguration != nil {
			immutable = objectLockConfig.ObjectLockConfiguration.ObjectLockEnabled == "Enabled"
		}
		buckets = append(buckets, S3BucketInfo{
			Name:      aws.ToString(bucket.Name),
			Region:    string(region.LocationConstraint),
			Immutable: immutable,
		})
	}

	return buckets, nil
}

func FetchRDSInstances(ctx context.Context) ([]RDSInstanceInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	regionsOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe regions, %v", err)
	}

	var allInstances []RDSInstanceInfo
	for _, region := range regionsOutput.Regions {
		regionCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region.RegionName))
		if err != nil {
			return nil, fmt.Errorf("unable to load SDK config for region %s, %v", *region.RegionName, err)
		}

		regionRdsClient := rds.NewFromConfig(regionCfg)
		result, err := regionRdsClient.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
		if err != nil {
			return nil, fmt.Errorf("unable to describe DB instances in region %s, %v", *region.RegionName, err)
		}

		for _, instance := range result.DBInstances {
			allInstances = append(allInstances, RDSInstanceInfo{
				InstanceID: aws.ToString(instance.DBInstanceIdentifier),
				Engine:     aws.ToString(instance.Engine),
				Status:     aws.ToString(instance.DBInstanceStatus),
				Region:     *region.RegionName,
			})
		}
	}

	return allInstances, nil
}

func FetchDynamoDBTables(ctx context.Context) ([]DynamoDBTableInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	regionsOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe regions, %v", err)
	}

	var allTables []DynamoDBTableInfo
	for _, region := range regionsOutput.Regions {
		regionCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region.RegionName))
		if err != nil {
			return nil, fmt.Errorf("unable to load SDK config for region %s, %v", *region.RegionName, err)
		}

		svc := dynamodb.NewFromConfig(regionCfg)
		result, err := svc.ListTables(ctx, &dynamodb.ListTablesInput{})
		if err != nil {
			return nil, fmt.Errorf("unable to list tables in region %s, %v", *region.RegionName, err)
		}

		for _, tableName := range result.TableNames {
			describeResult, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &tableName,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to describe table %s in region %s, %v", tableName, *region.RegionName, err)
			}
			allTables = append(allTables, DynamoDBTableInfo{
				TableName: tableName,
				Status:    string(describeResult.Table.TableStatus),
				Region:    *region.RegionName,
			})
		}
	}

	return allTables, nil
}

func FetchVPCs(ctx context.Context) ([]VPCInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	ec2Client := ec2.NewFromConfig(cfg)
	regionsOutput, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe regions, %v", err)
	}

	var allVPCs []VPCInfo
	for _, region := range regionsOutput.Regions {
		regionCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region.RegionName))
		if err != nil {
			return nil, fmt.Errorf("unable to load SDK config for region %s, %v", *region.RegionName, err)
		}

		regionEc2Client := ec2.NewFromConfig(regionCfg)
		result, err := regionEc2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
		if err != nil {
			return nil, fmt.Errorf("unable to describe VPCs in region %s, %v", *region.RegionName, err)
		}

		for _, vpc := range result.Vpcs {
			allVPCs = append(allVPCs, VPCInfo{
				VPCID:  aws.ToString(vpc.VpcId),
				State:  string(vpc.State),
				Region: *region.RegionName,
			})
		}
	}

	return allVPCs, nil
}

func CollectAWSData(ctx context.Context) (AWSData, error) {
	var data AWSData
	var err error

	data.EC2Instances, err = FetchEC2Instances(ctx)
	if err != nil {
		return AWSData{}, err
	}

	data.S3Buckets, err = FetchS3Buckets(ctx)
	if err != nil {
		return AWSData{}, err
	}

	data.RDSInstances, err = FetchRDSInstances(ctx)
	if err != nil {
		return AWSData{}, err
	}

	data.DynamoDBTables, err = FetchDynamoDBTables(ctx)
	if err != nil {
		return AWSData{}, err
	}

	data.VPCs, err = FetchVPCs(ctx)
	if err != nil {
		return AWSData{}, err
	}

	return data, nil
}

func CollectSnapshotData(ctx context.Context) (map[string]interface{}, error) {
	snapshots := map[string]interface{}{}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	ebsSnapshots, err := collectEBSSnapshots(ctx, cfg)
	if err != nil {
		log.Printf("Warning: Failed to collect EBS snapshots: %v", err)
	} else {
		snapshots["EBSSnapshots"] = ebsSnapshots
	}

	rdsSnapshots, err := collectRDSSnapshots(ctx, cfg)
	if err != nil {
		log.Printf("Warning: Failed to collect RDS snapshots: %v", err)
	} else {
		snapshots["RDSSnapshots"] = rdsSnapshots
	}

	return snapshots, nil
}

func collectEBSSnapshots(ctx context.Context, cfg aws.Config) ([]map[string]string, error) {
	client := ec2.NewFromConfig(cfg)

	resp, err := client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{"self"},
	})
	if err != nil {
		return nil, err
	}

	var snapshots []map[string]string
	for _, snapshot := range resp.Snapshots {
		snapshotInfo := map[string]string{
			"SnapshotId": *snapshot.SnapshotId,
			"VolumeId":   *snapshot.VolumeId,
			"State":      string(snapshot.State),
		}

		if snapshot.VolumeSize != nil {
			snapshotInfo["VolumeSize"] = fmt.Sprintf("%d GiB", *snapshot.VolumeSize)
		}

		if snapshot.StartTime != nil {
			snapshotInfo["StartTime"] = snapshot.StartTime.Format(time.RFC3339)
		}

		if snapshot.Description != nil {
			snapshotInfo["Description"] = *snapshot.Description
		}

		if snapshot.Encrypted != nil {
			if *snapshot.Encrypted {
				snapshotInfo["Encrypted"] = "true"
			} else {
				snapshotInfo["Encrypted"] = "false"
			}
		}

		snapshots = append(snapshots, snapshotInfo)
	}

	return snapshots, nil
}

func collectRDSSnapshots(ctx context.Context, cfg aws.Config) ([]map[string]string, error) {
	client := rds.NewFromConfig(cfg)

	resp, err := client.DescribeDBSnapshots(ctx, &rds.DescribeDBSnapshotsInput{})
	if err != nil {
		return nil, err
	}

	var snapshots []map[string]string
	for _, snapshot := range resp.DBSnapshots {
		snapshotInfo := map[string]string{
			"SnapshotId":   aws.ToString(snapshot.DBSnapshotIdentifier),
			"DBInstanceId": aws.ToString(snapshot.DBInstanceIdentifier),
			"SnapshotType": aws.ToString(snapshot.SnapshotType),
			"Status":       aws.ToString(snapshot.Status),
		}

		if snapshot.Engine != nil {
			snapshotInfo["Engine"] = *snapshot.Engine
		}

		if snapshot.AllocatedStorage != nil && *snapshot.AllocatedStorage != 0 {
			snapshotInfo["AllocatedStorage"] = fmt.Sprintf("%d GiB", *snapshot.AllocatedStorage)
		}

		if snapshot.SnapshotCreateTime != nil {
			snapshotInfo["CreationTime"] = snapshot.SnapshotCreateTime.Format(time.RFC3339)
		}

		if snapshot.Encrypted != nil {
			if *snapshot.Encrypted {
				snapshotInfo["Encrypted"] = "true"
			} else {
				snapshotInfo["Encrypted"] = "false"
			}
		}

		snapshots = append(snapshots, snapshotInfo)
	}

	return snapshots, nil
}
