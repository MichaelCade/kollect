package main

import (
	"bufio"
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/azure"
	"github.com/michaelcade/kollect/pkg/cost"
	"github.com/michaelcade/kollect/pkg/docker"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/kollect"
	"github.com/michaelcade/kollect/pkg/snapshots"
	"github.com/michaelcade/kollect/pkg/terraform"
	"github.com/michaelcade/kollect/pkg/vault"
	"github.com/michaelcade/kollect/pkg/veeam"
	"golang.org/x/term"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	dataMutex sync.Mutex
	data      interface{}
	//go:embed web/*
	staticFiles embed.FS
)

func main() {
	storageOnly := flag.Bool("storage", false, "Collect only storage-related objects (Kubernetes Only)")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	browser := flag.Bool("browser", false, "Open the web interface in a browser (can be used alone to import data)")
	dockerHost := flag.String("docker-host", "", "Docker host (e.g. unix:///var/run/docker.sock or tcp://host:2375)")
	output := flag.String("output", "", "Output file to save the collected data")
	inventoryType := flag.String("inventory", "", "Type of inventory to collect (kubernetes/aws/azure/gcp/terraform/vault/docker/veeam)")
	baseURL := flag.String("veeam-url", "", "Veeam server URL")
	username := flag.String("veeam-username", "", "Veeam username")
	password := flag.String("veeam-password", "", "Veeam password")
	terraformStateFile := flag.String("terraform-state", "", "Path to a local Terraform state file")
	terraformS3Bucket := flag.String("terraform-s3", "", "S3 bucket containing Terraform state (format: bucket/key)")
	terraformS3Region := flag.String("terraform-s3-region", "", "AWS region for S3 bucket (defaults to AWS_REGION env var)")
	terraformAzureContainer := flag.String("terraform-azure", "", "Azure storage container (format: storageaccount/container/blob)")
	terraformGCSBucket := flag.String("terraform-gcs", "", "GCS bucket and object (format: bucket/object)")
	kubeContext := flag.String("kube-context", "", "Kubernetes context to use")
	snapshotFlag := flag.Bool("snapshots", false, "Collect snapshots from all available platforms")
	vaultAddr := flag.String("vault-addr", "", "Vault server address")
	vaultToken := flag.String("vault-token", "", "Vault token")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()
	if *help {
		fmt.Println("Usage: kollect [flags]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println("\nTo pretty-print JSON output, you can use `jq`:")
		fmt.Println("  ./kollect | jq")
		return
	}

	if *snapshotFlag {
		fmt.Println("Collecting snapshots from all available platforms...")
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		snapshotData, err := snapshots.CollectAllSnapshots(context.Background(), kubeconfigPath)
		if err != nil {
			fmt.Printf("Error collecting snapshots: %v\n", err)
			os.Exit(1)
		}

		outputData := *output
		if outputData != "" {
			jsonData, err := json.MarshalIndent(snapshotData, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling data: %v\n", err)
				os.Exit(1)
			}

			err = os.WriteFile(outputData, jsonData, 0644)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Snapshot data saved to %s\n", outputData)
		} else {
			jsonData, err := json.MarshalIndent(snapshotData, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling data: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(string(jsonData))
		}

		return
	}

	if *browser && *inventoryType == "" && *output == "" {
		fmt.Println("Starting browser interface. Use the import function to load data.")
		startWebServer(map[string]interface{}{}, true, "", "", "")
		return
	}

	if *inventoryType == "" && !*snapshotFlag && !(*browser && *output == "") {
		fmt.Println("Error: You must specify an inventory type with --inventory")
		fmt.Println("Available inventory types: kubernetes, aws, azure, gcp, veeam, terraform, vault, docker")
		fmt.Println("Or use --browser alone to start web interface for importing data")
		fmt.Println("Or use --snapshots to collect snapshot data from all available platforms")
		os.Exit(1)
	}

	ctx := context.Background()

	var err error
	switch *inventoryType {
	case "aws":
		data, err = aws.CollectAWSData(ctx)
	case "azure":
		data, err = azure.CollectAzureData(ctx)
	case "gcp":
		data, err = gcp.CollectGCPData(ctx)
	case "kubernetes":
		if *kubeContext != "" {
			data, err = collectData(ctx, *storageOnly, *kubeconfig, *kubeContext)
		} else {
			data, err = collectData(ctx, *storageOnly, *kubeconfig)
		}
	case "terraform":
		if *terraformStateFile != "" {
			data, err = terraform.CollectTerraformData(ctx, *terraformStateFile)
		} else if *terraformS3Bucket != "" {
			parts := strings.SplitN(*terraformS3Bucket, "/", 2)
			if len(parts) != 2 {
				fmt.Println("Error: terraform-s3 flag must be in format 'bucket/key'")
				os.Exit(1)
			}
			region := *terraformS3Region
			if region == "" {
				region = os.Getenv("AWS_REGION")
				if region == "" {
					region = "us-east-1"
				}
			}
			data, err = terraform.CollectTerraformDataFromS3(ctx, parts[0], parts[1], region)
		} else if *terraformAzureContainer != "" {
			parts := strings.SplitN(*terraformAzureContainer, "/", 3)
			if len(parts) != 3 {
				fmt.Println("Error: terraform-azure flag must be in format 'storageaccount/container/blob'")
				os.Exit(1)
			}
			data, err = terraform.CollectTerraformDataFromAzure(ctx, parts[0], parts[1], parts[2])
		} else if *terraformGCSBucket != "" {
			parts := strings.SplitN(*terraformGCSBucket, "/", 2)
			if len(parts) != 2 {
				fmt.Println("Error: terraform-gcs flag must be in format 'bucket/object'")
				os.Exit(1)
			}
			data, err = terraform.CollectTerraformDataFromGCS(ctx, parts[0], parts[1])
		} else {
			fmt.Println("Error: You must specify a Terraform state source with one of:")
			fmt.Println("  --terraform-state (local file)")
			fmt.Println("  --terraform-s3 (AWS S3)")
			fmt.Println("  --terraform-azure (Azure Blob Storage)")
			fmt.Println("  --terraform-gcs (Google Cloud Storage)")
			os.Exit(1)
		}
	case "veeam":
		if *baseURL == "" {
			*baseURL = os.Getenv("VBR_SERVER_URL")
		}
		if *baseURL == "" {
			serverAddress := promptUser("Enter VBR Server IP or DNS name: ")
			*baseURL = fmt.Sprintf("https://%s:9419", serverAddress)
		}
		if *username == "" {
			*username = getEnv("VBR_USERNAME", "Enter VBR Username: ")
		}
		if *password == "" {
			*password = getSensitiveInput("Enter VBR Password: ")
		}
		if !strings.HasPrefix(*baseURL, "http://") && !strings.HasPrefix(*baseURL, "https://") {
			*baseURL = "http://" + *baseURL
		}

		data, err = veeam.CollectVeeamData(ctx, *baseURL, *username, *password, true)
	case "docker":
		data, err = docker.CollectDockerData(ctx, *dockerHost)

	default:
		log.Fatalf("Invalid inventory type: %s", *inventoryType)
	case "vault":
		vaultAddress := *vaultAddr
		vaultTokenVal := *vaultToken
		vaultIgnoreSSL := false

		if vaultAddress == "" {
			vaultAddress = os.Getenv("VAULT_ADDR")
			if vaultAddress == "" {
				vaultAddress = promptUser("Enter Vault server address (e.g. http://localhost:8200): ")
			}
		}

		if vaultTokenVal == "" {
			vaultTokenVal = os.Getenv("VAULT_TOKEN")
			if vaultTokenVal == "" {
				vaultTokenVal = getSensitiveInput("Enter Vault token: ")
			}
		}

		data, err = vault.CollectVaultData(ctx, vaultAddress, vaultTokenVal, vaultIgnoreSSL)
	}
	if err != nil {
		log.Fatalf("Error collecting data: %v", err)
	}

	if *output != "" {
		err = saveToFile(data, *output)
		if err != nil {
			log.Fatalf("Error saving data to file: %v", err)
		}
		fmt.Printf("Data saved to %s\n", *output)
		return
	}

	printData(data)

	if *browser {
		startWebServer(data, true, *baseURL, *username, *password)
	}

}

func collectAllSnapshots(ctx context.Context) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	credentials := checkCredentials(ctx)

	if credentials["kubernetes"] {
		fmt.Println("Collecting Kubernetes snapshots...")
		k8sSnapshots, err := kollect.CollectSnapshotData(ctx, filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			fmt.Printf("Warning: Error collecting Kubernetes snapshots: %v\n", err)
		} else if k8sSnapshots != nil {
			fmt.Println("Successfully collected Kubernetes snapshots")
			results["kubernetes"] = k8sSnapshots
		}
	}

	if credentials["aws"] {
		fmt.Println("Collecting AWS snapshots...")
		awsSnapshots, err := aws.CollectSnapshotData(ctx)
		if err != nil {
			fmt.Printf("Warning: Error collecting AWS snapshots: %v\n", err)
		} else if awsSnapshots != nil {
			fmt.Println("Successfully collected AWS snapshots")
			results["aws"] = awsSnapshots
		}
	}

	if credentials["azure"] {
		fmt.Println("Collecting Azure snapshots...")
		azureSnapshots, err := azure.CollectSnapshotData(ctx)
		if err != nil {
			fmt.Printf("Warning: Error collecting Azure snapshots: %v\n", err)
		} else if azureSnapshots != nil {
			fmt.Println("Successfully collected Azure snapshots")
			results["azure"] = azureSnapshots
		}
	}

	if credentials["gcp"] {
		fmt.Println("Collecting GCP snapshots...")
		gcpSnapshots, err := gcp.CollectSnapshotData(ctx)
		if err != nil {
			fmt.Printf("Warning: Error collecting GCP snapshots: %v\n", err)
		} else if gcpSnapshots != nil {
			fmt.Println("Successfully collected GCP snapshots")
			results["gcp"] = gcpSnapshots
		}
	}

	return results, nil
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getSensitiveInput(prompt string) string {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return promptUser("(echo enabled) " + prompt)
	}
	return string(password)
}

func getEnv(envVar, prompt string) string {
	value := os.Getenv(envVar)
	if value == "" {
		value = promptUser(prompt)
	}
	return value
}

func collectData(ctx context.Context, storageOnly bool, kubeconfigPath string, contextName ...string) (interface{}, error) {
	if storageOnly {
		return kollect.CollectStorageData(ctx, kubeconfigPath)
	}

	if len(contextName) > 0 && contextName[0] != "" {
		return kollect.CollectDataWithContext(ctx, kubeconfigPath, contextName[0])
	}

	return kollect.CollectData(ctx, kubeconfigPath)
}

func saveToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	prettyData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(prettyData)
	return err
}

func printData(data interface{}) {
	prettyData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting data: %v", err)
	}
	fmt.Println(string(prettyData))
}

func checkCredentials(ctx context.Context) map[string]bool {
	results := make(map[string]bool)

	k8sConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err == nil {
		clientset, err := kubernetes.NewForConfig(k8sConfig)
		if err == nil {
			_, err = clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{Limit: 1})
			results["kubernetes"] = err == nil
		}
	} else {
		results["kubernetes"] = false
	}

	awsHasCredentials, _ := aws.CheckCredentials(ctx)
	results["aws"] = awsHasCredentials

	azureHasCredentials, _ := azure.CheckCredentials(ctx)
	results["azure"] = azureHasCredentials

	gcpHasCredentials, _ := gcp.CheckCredentials(ctx)
	results["gcp"] = gcpHasCredentials

	dockerHasCredentials, _ := docker.CheckCredentials(ctx, "")
	results["docker"] = dockerHasCredentials

	dataMutex.Lock()
	veeamConnected := false
	if d, ok := data.(veeam.VeeamData); ok {
		veeamConnected = d.ServerInfo != nil && len(d.ServerInfo) > 0
	}
	dataMutex.Unlock()
	results["veeam"] = veeamConnected

	_, err = exec.LookPath("terraform")
	results["terraform"] = err == nil

	return results
}

func startWebServer(initialData interface{}, openBrowser bool, baseURL, username, password string) {
	fsys, err := fs.Sub(staticFiles, "web")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(fsys))
	cost.InitPricing()
	http.Handle("/", fileServer)

	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		dataMutex.Lock()
		defer dataMutex.Unlock()
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Error encoding data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/costs", cost.HandleCostRequest)
	http.HandleFunc("/api/refresh-pricing", cost.HandleRefreshPricing)
	http.HandleFunc("/api/pricing-info", cost.HandlePricingInfo)

	http.HandleFunc("/api/check-credentials", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		results := checkCredentials(ctx)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	http.HandleFunc("/api/import", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		var importedData interface{}
		err := json.NewDecoder(r.Body).Decode(&importedData)
		if err != nil {
			log.Printf("Error decoding imported data: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dataMutex.Lock()
		data = importedData
		dataMutex.Unlock()
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})

	http.HandleFunc("/api/switch", func(w http.ResponseWriter, r *http.Request) {
		inventoryType := r.URL.Query().Get("type")
		ctx := context.Background()
		var err error
		switch inventoryType {
		case "aws":
			data, err = aws.CollectAWSData(ctx)
		case "azure":
			data, err = azure.CollectAzureData(ctx)
		case "kubernetes":
			data, err = collectData(ctx, false, filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		case "gcp":
			data, err = gcp.CollectGCPData(ctx)
		case "terraform":
			if r.URL.Query().Get("state-file") == "" {
				http.Error(w, "Terraform state file must be provided", http.StatusBadRequest)
				return
			}
			data, err = terraform.CollectTerraformData(ctx, r.URL.Query().Get("state-file"))
		case "veeam":
			if baseURL == "" || username == "" || password == "" {
				http.Error(w, "Veeam URL, username, and password must be provided", http.StatusBadRequest)
				return
			}
			data, err = veeam.CollectVeeamData(ctx, baseURL, username, password, true)
		default:
			http.Error(w, "Invalid inventory type", http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Printf("Error collecting data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dataMutex.Lock()
		data = data
		dataMutex.Unlock()
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})

	http.HandleFunc("/api/terraform/s3-state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
			Region string `json:"region"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.Bucket == "" || params.Key == "" {
			http.Error(w, "Bucket and key are required", http.StatusBadRequest)
			return
		}

		if params.Region == "" {
			params.Region = "us-east-1"
		}

		ctx := context.Background()
		tfData, err := terraform.CollectTerraformDataFromS3(ctx, params.Bucket, params.Key, params.Region)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving state from S3: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tfData)
	})

	http.HandleFunc("/api/terraform/azure-state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Account   string `json:"account"`
			Container string `json:"container"`
			Blob      string `json:"blob"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.Account == "" || params.Container == "" || params.Blob == "" {
			http.Error(w, "Account, container, and blob are required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		tfData, err := terraform.CollectTerraformDataFromAzure(ctx, params.Account, params.Container, params.Blob)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving state from Azure: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tfData)
	})

	http.HandleFunc("/api/terraform/gcs-state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Bucket string `json:"bucket"`
			Object string `json:"object"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.Bucket == "" || params.Object == "" {
			http.Error(w, "Bucket and object are required", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		tfData, err := terraform.CollectTerraformDataFromGCS(ctx, params.Bucket, params.Object)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving state from GCS: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tfData)
	})

	http.HandleFunc("/api/veeam/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			BaseUrl   string `json:"baseUrl"`
			Username  string `json:"username"`
			Password  string `json:"password"`
			IgnoreSSL bool   `json:"ignoreSSL"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.BaseUrl == "" || params.Username == "" || params.Password == "" {
			http.Error(w, "URL, username, and password are required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		veeamData, err := veeam.CollectVeeamData(ctx, params.BaseUrl, params.Username, params.Password, params.IgnoreSSL)
		if err != nil {
			log.Printf("Error connecting to Veeam server %s: %v", params.BaseUrl, err)
			http.Error(w, fmt.Sprintf("Error connecting to Veeam: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = veeamData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to Veeam server",
		})
	})

	http.HandleFunc("/api/vault/cli-status", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("vault", "version")
		output, err := cmd.CombinedOutput()

		status := map[string]interface{}{
			"installed": err == nil,
		}

		if err == nil {
			versionStr := strings.TrimSpace(string(output))
			versionParts := strings.Split(versionStr, " ")
			if len(versionParts) >= 2 {
				status["version"] = versionParts[1]
			}

			checkCmd := exec.Command("vault", "token", "lookup")
			checkErr := checkCmd.Run()
			status["authenticated"] = checkErr == nil
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	http.HandleFunc("/api/vault/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Type            string `json:"type"`
			Server          string `json:"server"`
			Token           string `json:"token"`
			Username        string `json:"username"`
			Password        string `json:"password"`
			AuthPath        string `json:"authPath"`
			Insecure        bool   `json:"insecure"`
			IncludePolicies bool   `json:"includePolicies"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		var vaultToken string
		vaultAddr := params.Server

		if params.Type == "token" {
			if params.Server == "" || params.Token == "" {
				http.Error(w, "Server and token are required for token authentication", http.StatusBadRequest)
				return
			}
			vaultToken = params.Token
		} else if params.Type == "userpass" {
			if params.Server == "" || params.Username == "" || params.Password == "" {
				http.Error(w, "Server, username, and password are required for userpass authentication", http.StatusBadRequest)
				return
			}

			http.Error(w, "Userpass authentication not yet implemented", http.StatusNotImplemented)
			return
		} else if params.Type == "cli" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				tokenFile := filepath.Join(homeDir, ".vault-token")
				if tokenBytes, err := os.ReadFile(tokenFile); err == nil {
					vaultToken = strings.TrimSpace(string(tokenBytes))
				}
			}

			if vaultToken == "" {
				vaultToken = os.Getenv("VAULT_TOKEN")
			}

			if vaultToken == "" {
				http.Error(w, "Could not find Vault token in environment or token file", http.StatusUnauthorized)
				return
			}

			if vaultAddr == "" {
				vaultAddr = os.Getenv("VAULT_ADDR")
				if vaultAddr == "" {
					vaultAddr = "http://127.0.0.1:8200"
				}
			}
		}

		hasCredentials, err := vault.CheckCredentials(ctx, vaultAddr, vaultToken, params.Insecure)
		if err != nil || !hasCredentials {
			http.Error(w, fmt.Sprintf("Error connecting to Vault: %v", err), http.StatusUnauthorized)
			return
		}

		vaultData, err := vault.CollectVaultData(ctx, vaultAddr, vaultToken, params.Insecure)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error collecting Vault data: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = vaultData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to Vault",
		})
	})

	http.HandleFunc("/api/kubernetes/contexts", func(w http.ResponseWriter, r *http.Request) {
		kubeconfigPath := r.URL.Query().Get("path")
		if kubeconfigPath == "" {
			kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Kubeconfig file not found: %s", kubeconfigPath), http.StatusBadRequest)
			return
		}

		config, err := clientcmd.LoadFromFile(kubeconfigPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error loading kubeconfig: %v", err), http.StatusInternalServerError)
			return
		}

		contexts := make([]map[string]string, 0)
		for name, context := range config.Contexts {
			ctxInfo := map[string]string{
				"name":      name,
				"cluster":   context.Cluster,
				"namespace": context.Namespace,
				"user":      context.AuthInfo,
				"current":   "false",
			}

			if name == config.CurrentContext {
				ctxInfo["current"] = "true"
			}

			contexts = append(contexts, ctxInfo)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"contexts":       contexts,
			"currentContext": config.CurrentContext,
		})
	})

	http.HandleFunc("/api/kubernetes/upload-kubeconfig", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("kubeconfig")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		tempDir := os.TempDir()
		tempFileName := fmt.Sprintf("kubeconfig_%d_%s", time.Now().UnixNano(), handler.Filename)
		tempFilePath := filepath.Join(tempDir, tempFileName)

		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create temporary file: %v", err), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		if _, err := io.Copy(tempFile, file); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"path":   tempFilePath,
		})
	})

	http.HandleFunc("/api/kubernetes/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			KubeconfigPath string `json:"kubeconfigPath"`
			Context        string `json:"context"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.KubeconfigPath == "" {
			params.KubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		if _, err := os.Stat(params.KubeconfigPath); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Kubeconfig file not found: %s", params.KubeconfigPath), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		kubeData, err := collectData(ctx, false, params.KubeconfigPath, params.Context)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error connecting to Kubernetes: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = kubeData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to Kubernetes cluster",
		})

	})

	http.HandleFunc("/api/aws/profiles", func(w http.ResponseWriter, r *http.Request) {
		profiles := []string{"default"}
		homeDir, err := os.UserHomeDir()
		if err == nil {
			awsDir := filepath.Join(homeDir, ".aws")

			credPath := filepath.Join(awsDir, "credentials")
			if _, err := os.Stat(credPath); err == nil {
				data, err := os.ReadFile(credPath)
				if err == nil {
					lines := strings.Split(string(data), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
							profile := line[1 : len(line)-1]
							if profile != "default" && !contains(profiles, profile) {
								profiles = append(profiles, profile)
							}
						}
					}
				}
			}

			configPath := filepath.Join(awsDir, "config")
			if _, err := os.Stat(configPath); err == nil {
				data, err := os.ReadFile(configPath)
				if err == nil {
					lines := strings.Split(string(data), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "[profile ") && strings.HasSuffix(line, "]") {
							profile := line[9 : len(line)-1]
							if profile != "default" && !contains(profiles, profile) {
								profiles = append(profiles, profile)
							}
						}
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"profiles": profiles,
		})
	})

	http.HandleFunc("/api/aws/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Type      string `json:"type"`
			Profile   string `json:"profile"`
			AccessKey string `json:"accessKey"`
			SecretKey string `json:"secretKey"`
			Region    string `json:"region"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.Type == "credentials" {
			os.Setenv("AWS_ACCESS_KEY_ID", params.AccessKey)
			os.Setenv("AWS_SECRET_ACCESS_KEY", params.SecretKey)
			if params.Region != "" {
				os.Setenv("AWS_REGION", params.Region)
			}
		} else if params.Type == "profile" && params.Profile != "" {
			os.Setenv("AWS_PROFILE", params.Profile)
		}

		ctx := r.Context()
		hasCredentials, err := aws.CheckCredentials(ctx)
		if err != nil || !hasCredentials {
			http.Error(w, fmt.Sprintf("Error connecting to AWS: %v", err), http.StatusBadRequest)
			return
		}

		awsData, err := aws.CollectAWSData(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error collecting AWS data: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = awsData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to AWS",
		})
	})

	http.HandleFunc("/api/azure/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		subscriptions := []map[string]string{}
		defaultSub := ""
		subscriptionsCount := 0

		cmd := exec.Command("az", "account", "list", "--query", "[].{name:name, id:id, isDefault:isDefault}", "--output", "json")
		output, err := cmd.Output()

		if err == nil {
			var azSubs []map[string]interface{}
			if err := json.Unmarshal(output, &azSubs); err == nil {
				subscriptionsCount = len(azSubs)

				log.Printf("Found %d subscriptions in Azure CLI output", subscriptionsCount)

				for _, sub := range azSubs {
					if name, ok := sub["name"].(string); ok {
						if id, ok := sub["id"].(string); ok {
							subscription := map[string]string{
								"name":      name,
								"id":        id,
								"isDefault": "false",
							}

							if isDefault, ok := sub["isDefault"].(bool); ok && isDefault {
								subscription["isDefault"] = "true"
								defaultSub = id
							}

							subscriptions = append(subscriptions, subscription)
						}
					}
				}

				log.Printf("Processed %d Azure subscriptions", len(subscriptions))
			} else {
				log.Printf("Error parsing Azure CLI output: %v", err)
			}
		} else {
			log.Printf("Error running Azure CLI command: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"subscriptions":      subscriptions,
			"default":            defaultSub,
			"subscriptionsCount": subscriptionsCount,
		})
	})

	http.HandleFunc("/api/azure/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Type           string `json:"type"`
			Subscription   string `json:"subscription"`
			TenantId       string `json:"tenantId"`
			ClientId       string `json:"clientId"`
			ClientSecret   string `json:"clientSecret"`
			SubscriptionId string `json:"subscriptionId"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if params.Type == "service_principal" {
			os.Setenv("AZURE_TENANT_ID", params.TenantId)
			os.Setenv("AZURE_CLIENT_ID", params.ClientId)
			os.Setenv("AZURE_CLIENT_SECRET", params.ClientSecret)
			os.Setenv("AZURE_SUBSCRIPTION_ID", params.SubscriptionId)
		} else if params.Type == "cli" && params.Subscription != "" {
			cmd := exec.Command("az", "account", "set", "--subscription", params.Subscription)
			if err := cmd.Run(); err != nil {
				http.Error(w, fmt.Sprintf("Error setting Azure subscription: %v", err), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		hasCredentials, err := azure.CheckCredentials(ctx)
		if err != nil || !hasCredentials {
			http.Error(w, fmt.Sprintf("Error connecting to Azure: %v", err), http.StatusBadRequest)
			return
		}

		azureData, err := azure.CollectAzureData(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "Authorization") ||
				strings.Contains(err.Error(), "authorization") ||
				strings.Contains(err.Error(), "permission") ||
				strings.Contains(err.Error(), "access") {
				http.Error(w, "Permission denied: You don't have sufficient permissions for some Azure resources", http.StatusForbidden)
				return
			}

			http.Error(w, fmt.Sprintf("Error collecting Azure data: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = azureData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to Azure",
		})
	})

	http.HandleFunc("/api/gcp/projects", func(w http.ResponseWriter, r *http.Request) {
		projects := []map[string]string{}

		cmd := exec.Command("gcloud", "projects", "list", "--format=json")
		output, err := cmd.CombinedOutput()
		if err == nil {
			var gcpProjects []map[string]interface{}
			if err := json.Unmarshal(output, &gcpProjects); err == nil {
				for _, project := range gcpProjects {
					projectInfo := map[string]string{
						"name":      fmt.Sprintf("%v", project["name"]),
						"id":        fmt.Sprintf("%v", project["projectId"]),
						"isDefault": "false",
					}

					projects = append(projects, projectInfo)
				}
			}
		}

		cmd = exec.Command("gcloud", "config", "get-value", "project")
		output, err = cmd.CombinedOutput()
		if err == nil {
			defaultProject := strings.TrimSpace(string(output))
			for i, project := range projects {
				if project["id"] == defaultProject {
					projects[i]["isDefault"] = "true"
					break
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"projects": projects,
		})
	})

	http.HandleFunc("/api/snapshots", func(w http.ResponseWriter, r *http.Request) {
		platform := r.URL.Query().Get("platform")

		ctx := r.Context()
		kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

		var data map[string]interface{}
		var err error

		if platform == "all" {
			data, err = snapshots.CollectAllSnapshots(ctx, kubeconfigPath)
		} else {
			data, err = snapshots.CollectPlatformSnapshots(ctx, platform, kubeconfigPath)
		}

		if err != nil {
			log.Printf("Error collecting snapshots: %v", err)
			http.Error(w, fmt.Sprintf("Error collecting snapshots: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.HandleFunc("/api/gcp/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Type    string                 `json:"type"`
			Project string                 `json:"project"`
			KeyData map[string]interface{} `json:"keyData"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		tempKeyFile := ""

		if params.Type == "service_account" && params.KeyData != nil {
			keyBytes, err := json.Marshal(params.KeyData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing service account key: %v", err), http.StatusBadRequest)
				return
			}

			tempFile, err := os.CreateTemp("", "gcp-key-*.json")
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating temporary key file: %v", err), http.StatusInternalServerError)
				return
			}
			defer tempFile.Close()

			if _, err := tempFile.Write(keyBytes); err != nil {
				http.Error(w, fmt.Sprintf("Error writing key data: %v", err), http.StatusInternalServerError)
				return
			}

			tempKeyFile = tempFile.Name()
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", tempKeyFile)
		} else if params.Type == "gcloud" && params.Project != "" {
			cmd := exec.Command("gcloud", "config", "set", "project", params.Project)
			if err := cmd.Run(); err != nil {
				http.Error(w, fmt.Sprintf("Error setting GCP project: %v", err), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		hasCredentials, err := gcp.CheckCredentials(ctx)
		if err != nil || !hasCredentials {
			if tempKeyFile != "" {
				os.Remove(tempKeyFile)
				os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
			}
			http.Error(w, fmt.Sprintf("Error connecting to GCP: %v", err), http.StatusBadRequest)
			return
		}

		gcpData, err := gcp.CollectGCPData(ctx)
		if err != nil {
			if tempKeyFile != "" {
				os.Remove(tempKeyFile)
				os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
			}
			http.Error(w, fmt.Sprintf("Error collecting GCP data: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = gcpData
		dataMutex.Unlock()

		if tempKeyFile != "" {
			os.Remove(tempKeyFile)
			os.Unsetenv("GOOGLE_APPLICATION_CREDENTIALS")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to GCP",
		})
	})
	http.HandleFunc("/api/docker/test-connection", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Host string `json:"host"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		version, err := docker.TestConnection(ctx, params.Host)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot connect to the Docker daemon at %s. Is the docker daemon running?", params.Host), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"version": version,
			"message": "Successfully connected to Docker",
		})
	})

	http.HandleFunc("/api/docker/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Host string `json:"host"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		hasCredentials, err := docker.CheckCredentials(ctx, params.Host)
		if err != nil || !hasCredentials {
			http.Error(w, fmt.Sprintf("Error connecting to Docker: %v", err), http.StatusBadRequest)
			return
		}

		dockerData, err := docker.CollectDockerData(ctx, params.Host)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error collecting Docker data: %v", err), http.StatusInternalServerError)
			return
		}

		dataMutex.Lock()
		data = dockerData
		dataMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Successfully connected to Docker",
		})
	})

	log.Println("Server starting on port http://localhost:8080")
	if openBrowser {
		go func() {
			var err error
			switch runtime.GOOS {
			case "darwin":
				err = exec.Command("open", "http://localhost:8080").Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()
			default:
				err = exec.Command("xdg-open", "http://localhost:8080").Start()
			}
			if err != nil {
				log.Printf("Warning: Failed to open browser: %v", err)
			}
		}()
	}
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
