package aws

import (
	"context"
	"fmt"

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

func FetchEC2Instances() ([]EC2InstanceInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := ec2.NewFromConfig(cfg)
	result, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe instances, %v", err)
	}

	var instances []EC2InstanceInfo
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, EC2InstanceInfo{
				Name:       aws.ToString(instance.KeyName),
				InstanceID: aws.ToString(instance.InstanceId),
				Type:       string(instance.InstanceType),
				State:      string(instance.State.Name),
				Region:     cfg.Region,
			})
		}
	}

	return instances, nil
}

func FetchS3Buckets() ([]S3BucketInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := s3.NewFromConfig(cfg)
	result, err := svc.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list buckets, %v", err)
	}

	var buckets []S3BucketInfo
	for _, bucket := range result.Buckets {
		region, err := svc.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to get bucket location for %s, %v", aws.ToString(bucket.Name), err)
		}

		// Check if the bucket has an Object Lock configuration (immutability)
		objectLockConfig, err := svc.GetObjectLockConfiguration(context.TODO(), &s3.GetObjectLockConfigurationInput{
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

func FetchRDSInstances() ([]RDSInstanceInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := rds.NewFromConfig(cfg)
	result, err := svc.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe DB instances, %v", err)
	}

	var instances []RDSInstanceInfo
	for _, instance := range result.DBInstances {
		instances = append(instances, RDSInstanceInfo{
			InstanceID: aws.ToString(instance.DBInstanceIdentifier),
			Engine:     aws.ToString(instance.Engine),
			Status:     aws.ToString(instance.DBInstanceStatus),
			Region:     cfg.Region,
		})
	}

	return instances, nil
}

func FetchDynamoDBTables() ([]DynamoDBTableInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	result, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list tables, %v", err)
	}

	var tables []DynamoDBTableInfo
	for _, tableName := range result.TableNames {
		describeResult, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: &tableName,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to describe table %s, %v", tableName, err)
		}
		tables = append(tables, DynamoDBTableInfo{
			TableName: tableName,
			Status:    string(describeResult.Table.TableStatus),
			Region:    cfg.Region,
		})
	}

	return tables, nil
}

func FetchVPCs() ([]VPCInfo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	svc := ec2.NewFromConfig(cfg)
	result, err := svc.DescribeVpcs(context.TODO(), &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to describe VPCs, %v", err)
	}

	var vpcs []VPCInfo
	for _, vpc := range result.Vpcs {
		vpcs = append(vpcs, VPCInfo{
			VPCID:  aws.ToString(vpc.VpcId),
			State:  string(vpc.State),
			Region: cfg.Region,
		})
	}

	return vpcs, nil
}

func CollectAWSData() (AWSData, error) {
	var data AWSData
	var err error

	data.EC2Instances, err = FetchEC2Instances()
	if err != nil {
		return AWSData{}, err
	}

	data.S3Buckets, err = FetchS3Buckets()
	if err != nil {
		return AWSData{}, err
	}

	data.RDSInstances, err = FetchRDSInstances()
	if err != nil {
		return AWSData{}, err
	}

	data.DynamoDBTables, err = FetchDynamoDBTables()
	if err != nil {
		return AWSData{}, err
	}

	data.VPCs, err = FetchVPCs()
	if err != nil {
		return AWSData{}, err
	}

	return data, nil
}
