# Kollect

Kollect is a tool designed to collect various Kubernetes resources and provide an overview of the cluster's state. It gathers data on VolumeSnapshotClasses, Nodes, Namespaces, StatefulSets, PersistentVolumes, PersistentVolumeClaims, and StorageClasses.

## Features

- Collects data on VolumeSnapshotClasses, Nodes, Namespaces, StatefulSets, PersistentVolumes, PersistentVolumeClaims, and StorageClasses.
- Provides a comprehensive overview of the Kubernetes cluster's state.
- Easy to configure and use.

## Installation

To install Kollect, you need to have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

Clone the repository:

```sh
git clone https://github.com/yourusername/kollect.git
cd kollect
```

Build the project:

```sh
go build -o kollect
```

## Usage

To use Kollect, you need to provide the path to your Kubernetes configuration file.

```sh
./kollect --kubeconfig /path/to/your/kubeconfig
```

## Example

```sh
./kollect --kubeconfig ~/.kube/config
```

## Contributing


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.


