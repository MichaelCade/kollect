# Kollect

Kollect is a tool for collecting and displaying Kubernetes cluster data. It provides a web interface to visualize various Kubernetes resources and allows exporting the collected data as a JSON file.

## Features

- Collects data from Kubernetes clusters
- Displays data in a web interface
- Supports exporting data as a JSON file

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

- `-storage`: Collect only storage-related objects
- `-kubeconfig`: Path to the kubeconfig file (default: `$HOME/.kube/config`)
- `-browser`: Open the web interface in a browser
- `-output`: Output file to save the collected data
- `-inventory`: Type of inventory to collect (kubernetes, aws, azure, gcp) (default: kubernetes)
- `-help`: Show help message

### Examples

Collect all data and open the web interface in a browser:

```sh
./kollect -browser
```

Collect only storage-related objects and save the data to a file:

```sh
./kollect -storage -output data.json
```

## Web Interface

The web interface provides a visual representation of the collected data. It includes tables for various Kubernetes resources such as Nodes, Namespaces, Pods, Deployments, StatefulSets, Services, PersistentVolumes, PersistentVolumeClaims, StorageClasses, and VolumeSnapshotClasses.

### Export Data

You can export the collected data as a JSON file by clicking the "Export Data" button in the web interface.

## Development

To contribute to Kollect, follow these steps:

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Make your changes
4. Commit your changes (`git commit -am 'Add new feature'`)
5. Push to the branch (`git push origin feature-branch`)
6. Create a new Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Material UI](https://material-ui.com/) for the web interface components
- [HTMX](https://htmx.org/) and [Hyperscript](https://hyperscript.org/) for dynamic content loading

## Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/michaelcade/kollect/issues).