# Kollect

Kollect is a tool designed to collect various Kubernetes resources and provide an overview of the cluster's state. It gathers data on VolumeSnapshotClasses, Nodes, Namespaces, StatefulSets, PersistentVolumes, PersistentVolumeClaims, and StorageClasses.

## Features

- Collects data on VolumeSnapshotClasses, Nodes, Namespaces, StatefulSets, PersistentVolumes, PersistentVolumeClaims, and StorageClasses.
- Provides a comprehensive overview of the Kubernetes cluster's state.
- Easy to configure and use.
- Supports both storage-specific and general data collection modes.
- Serves collected data via a web interface.

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

Kollect uses the current context from `kubectl` on your machine. Simply run the following command:

```sh
./kollect
```

### Storage-Specific Data Collection

To collect storage-specific data, use the `--storage` flag:

```sh
./kollect --storage
```

## Example

```sh
./kollect
```

## Web Interface

Kollect serves a web interface to view the collected data. By default, the server starts on port 8080.

To access the web interface, open your browser and navigate to:

```
http://localhost:8080
```

## API Endpoints

- `/api/data`: Returns the collected data in JSON format.

## Getting Started

1. Ensure you have Go installed on your machine.
2. Clone the repository and navigate to the project directory.
3. Build the project using the provided build command.
4. Run Kollect to start collecting data.
5. Access the web interface to view the collected data.

## Contributing

## License


```