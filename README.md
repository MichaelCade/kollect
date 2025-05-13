
# Kollect

Kollect is a tool for collecting and displaying data from Kubernetes clusters, AWS, and Azure resources. It provides a web interface to visualize various resources and allows exporting the collected data as a JSON file.

## Features

- Collects data from Kubernetes clusters
- Collects data from AWS resources (EC2, S3, RDS, DynamoDB, VPCs)
- Collects data from Azure resources (VMs, Storage Accounts, Blob Storage, Virtual Networks, SQL Databases, File Shares, CosmosDB)
- Displays data in a web interface
- Supports exporting data as a JSON file

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

- `--inventory`: Type of inventory to collect (kubernetes/aws/azure)
- `--storage`: Collect only storage-related objects (default: false)
- `--kubeconfig`: Path to the kubeconfig file (default: $HOME/.kube/config)
- `--browser`: Open the web interface in a browser (default: false)
- `--output`: Output file to save the collected data
- `--help`: Show help message

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

Collect data from a Kubernetes cluster and open the web interface:

```sh
./kollect --inventory kubernetes --browser
```

Collect data from AWS resources and save it to a file:

```sh
./kollect --inventory aws --output aws_data.json
```

## Development

### Project Structure

```
.DS_Store
.github/
    workflows/
        release.yaml
.gitignore
api/
    .DS_Store
    v1/
        k8sdata.go
cmd/
    .DS_Store
    kollect/
        main.go


go.mod




go.sum


LICENSE
pkg/
    .DS_Store
    aws/
        inventory.go
    azure/
        inventory.go
    kollect/
        kollect.go


README.md


test/
    kollect_test.go
web/
    .DS_Store
    index.html
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

```
go run cmd/kollect/main.go --inventory kubernetes --storage | jq
```

With an example of this as 

```
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
    },
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

## Contributing

We welcome contributions to Kollect! Please open an issue or submit a pull request on GitHub.

## License

Kollect is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
```