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

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgements

- [Material UI](https://material-ui.com/) for the web interface components
- [HTMX](https://htmx.org/) and [Hyperscript](https://hyperscript.org/) for dynamic content loading

## Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/michaelcade/kollect/issues).
```

### Summary of Changes

1. **Project Description**: Added a brief description of the project and its features.
2. **Installation**: Provided instructions for cloning the repository and building the binary.
3. **Usage**: Included detailed usage instructions and examples for running the binary with different flags.
4. **Web Interface**: Described the web interface and its features, including the export functionality.
5. **Development**: Added contribution guidelines for developers.
6. **License**: Mentioned the project license.
7. **Acknowledgements**: Acknowledged the use of Material UI, HTMX, and Hyperscript.
8. **Contact**: Provided a link to the GitHub repository for questions and feedback.

By following these updates, the 

README.md

 file will provide comprehensive information about the project, making it easier for users and contributors to understand and use Kollect.

MichaelCade: @workspace it does not seem like i can copy and paste the complete readme.md from this chat 

Workspace: Let's update the 

README.md

 file step by step. I'll provide the complete content for the updated 

README.md

 file, and you can copy and paste it into your file.

### Updated 

README.md



```md
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
```

### Summary of Changes

1. **Project Description**: Added a brief description of the project and its features.
2. **Installation**: Provided instructions for cloning the repository and building the binary.
3. **Usage**: Included detailed usage instructions and examples for running the binary with different flags.
4. **Web Interface**: Described the web interface and its features, including the export functionality.
5. **Development**: Added contribution guidelines for developers.
6. **License**: Mentioned the project license.
7. **Acknowledgements**: Acknowledged the use of Material UI, HTMX, and Hyperscript.
8. **Contact**: Provided a link to the GitHub repository for questions and feedback.
