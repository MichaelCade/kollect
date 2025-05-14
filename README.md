# Kollect

Kollect is a tool for collecting and displaying data from Kubernetes clusters, AWS, Azure, Google Cloud, and Veeam resources. It provides a web interface to visualize various resources and allows exporting the collected data as a JSON file.

## Features

- Collects data from Kubernetes clusters (including KubeVirt VMs and CRDs)
- Collects data from AWS resources (EC2, S3, RDS, DynamoDB, VPCs)
- Collects data from Azure resources (VMs, Storage Accounts, Blob Storage, Virtual Networks, SQL Databases, File Shares, CosmosDB)
- Collects data from Google Cloud resources (Compute Instances, Storage Buckets, SQL Instances, VPCs)
- Collects data from Veeam Backup & Replication servers (Backup Jobs, Repositories, Proxies, Scale-out Repositories)
- Inventory data from a Terraform state file (.tfstate / .json) (Local, AWS S3, Azure Blob, Google Cloud Storage)
- Displays data in a web interface
- Supports exporting data as a JSON file

## Security & Credentials

**Important:** Kollect does not store, transmit, or share any credentials. The tool works by:

- Using your existing local configurations (kubeconfig, AWS/Azure/GCP profiles)
- Leveraging environment variables when available
- Prompting for credentials only when necessary (e.g., Veeam connections)
- Never persisting credentials to disk
- Only collecting inventory data - no backup content or sensitive configuration data

All collected data is stored locally and only visualized in your browser on your local machine.

## Diagram 
![](diagram.png)

[GitDiagram provided the ability to create the above diagram](https://gitdiagram.com/michaelcade/kollect)

## Resources 
- [Kollect - The Cloud Inventory Project](https://youtu.be/dfuQFjl1Tnw)
- [Kollect - Veeam Inventory](https://youtu.be/yQ1vlndXTQY)
- [Kollect - A Cloud & Kubernetes Inventory tool](https://community.veeam.com/kubernetes%2Dkorner%2D90/kollect%2Da%2Dcloud%2Dkubernetes%2Dinventory%2Dtool%2D8885)
- [Kollect: A Modern Take on RVTools for Cloud Environments](https://community.veeam.com/kubernetes-korner-90/kollect-a-modern-take-on-rvtools-for-cloud-environments-9472?tid=9472&fid=90)

## Installation

To install Kollect, clone the repository and build the binary:

```sh
git clone https://github.com/michaelcade/kollect.git
cd kollect
go build -o kollect ./cmd/kollect
```

## Usage

Run the Kollect binary with the desired flags:

```sh
./kollect [flags]
```

### Flags

- `--inventory=<source>`: Specifies the inventory source. Valid values are `kubernetes`, `aws`, `azure`, `gcp`, `terraform`, and `veeam`.
- `--kubeconfig=<file>`: Specifies a custom kubeconfig file for Kubernetes.
- `--output=<file>`: Specifies an output file to save the collected data as JSON.
- `--browser`: Opens the collected data in a web browser.
- `--terraform-state=<file>`: Specifies a local Terraform state file to read.
- `--terraform-s3=<bucket/key>`: Specifies an AWS S3 bucket and key for Terraform state.
- `--terraform-s3-region=<region>`: Specifies the AWS region for the S3 bucket.
- `--terraform-azure=<storage-account/container/blob>`: Specifies an Azure Blob Storage location for Terraform state.
- `--terraform-gcs=<bucket/object>`: Specifies a Google Cloud Storage location for Terraform state.
- `--base-url=<url>`: Specifies the base URL for Veeam connections.
- `--username=<username>`: Specifies the username for Veeam connections.
- `--password=<password>`: Specifies the password for Veeam connections.

### Examples

Collect data from a Kubernetes cluster and display it in the terminal:

```sh
./kollect --inventory kubernetes
```

Collect data from AWS resources and display it in the terminal:

```sh
./kollect --inventory aws
```

Collect data from Azure resources and display it in the terminal:

```sh
./kollect --inventory azure
```

Collect data from Google Cloud resources and display it in the terminal: 

```sh
./kollect --inventory gcp
```

Collect data from Veeam Backup & Replication resources and display it in the terminal: 

```sh
./kollect --inventory veeam --base-url https://vbr-server.example.com:9419 --username admin --password password
```

We also have the ability to use the browser so you can import JSON format data:

```sh
./kollect --browser
```

Collect data from a Kubernetes cluster and open the web interface:

```sh
./kollect --inventory kubernetes --browser
```

Collect data from AWS resources and save it to a file:

```sh
./kollect --inventory aws --output aws_data.json
```

### Building the Project

To build the project, run the following command:

```sh
go build -o kollect ./cmd/kollect
```

### Running Tests (Under Review)

To run the tests, use the following command:

```sh
go test ./...
```

## Outputs 

You are able to export to JSON your data from the browser function or of course you will get an output in JSON format to the terminal on each run and each inventory of your desired platform. 

```sh
go run cmd/kollect/main.go --inventory kubernetes --storage | jq
```

With an example of this as:

```json
{
  "Nodes": null,
  "Namespaces": null,
  "Pods": null,
  "Deployments": null,
  "StatefulSets": null,
  "Services": null,
  "PersistentVolumes": [
    {
      "Name": "kasten-nfs-pv",
      "Capacity": "100Gi",
      "AccessModes": "ReadWriteMany",
      "Status": "Bound",
      "AssociatedClaim": "kasten-nfs-pvc",
      "StorageClass": "kasten-nfs",
      "VolumeMode": "Filesystem"
    },
    {
      "Name": "pvc-004b46dc-657e-44b3-ba7f-1d5c33f8278f",
      "Capacity": "10Gi",
      "AccessModes": "ReadWriteOnce",
      "Status": "Bound",
      "AssociatedClaim": "debian12-iso",
      "StorageClass": "ceph-block",
      "VolumeMode": "Filesystem"
    },
    {
      "Name": "pvc-053a8311-dca1-45f9-900b-59366243a985",
      "Capacity": "30Gi",
      "AccessModes": "ReadWriteOnce",
      "Status": "Bound",
      "AssociatedClaim": "ollama-volume-ollama-0",
      "StorageClass": "ceph-block",
      "VolumeMode": "Filesystem"
    }
  ],
  "StorageClasses": [
    {
      "Name": "ceph-block",
      "Provisioner": "rook-ceph.rbd.csi.ceph.com",
      "VolumeExpansion": "true"
    },
    {
      "Name": "ceph-bucket",
      "Provisioner": "rook-ceph.ceph.rook.io/bucket",
      "VolumeExpansion": "false"
    },
    {
      "Name": "ceph-filesystem",
      "Provisioner": "rook-ceph.cephfs.csi.ceph.com",
      "VolumeExpansion": "true"
    }
  ],
  "VolumeSnapshotClasses": [
    {
      "Name": "ceph-block-sc",
      "Driver": "rook-ceph.rbd.csi.ceph.com"
    },
    {
      "Name": "ceph-filesystem-sc",
      "Driver": "rook-ceph.cephfs.csi.ceph.com"
    }
  ],
  "VolumeSnapshots": null
}
```

## Recent Improvements

### UI and User Experience
- **Enhanced Loading Indicators**: Added visual feedback during cloud platform connections and data loading operations
- **Improved Error Handling**: Better error messages for connection issues, particularly for Azure subscription access problems
- **Consistent Visual Cues**: Loading spinners now display for all cloud platform operations

### Platform Integration Enhancements
- **Azure Integration Fixes**: 
  - Improved subscription detection and display
  - Added resilience for permission-related issues
  - Prevent application crashes when encountering access denied errors
  - Better error messages for subscription query failures

### Backend Improvements
- **Error Handling**: Gracefully handle authorization failures instead of crashing
- **Logging Improvements**: More detailed logs for troubleshooting API connections
- **Stability**: Enhanced error recovery for API requests to various cloud platforms

### Technical Details
- Added a mutation observer to automatically enhance dynamically created UI elements
- Updated loading indicator z-index to ensure visibility above modals
- Fixed Azure API error handling to continue collecting accessible resources when permission errors occur

## Contributing

We welcome contributions to Kollect! Please open an issue or submit a pull request on GitHub.

## License

Kollect is licensed under the MIT License. See the LICENSE file for more information.