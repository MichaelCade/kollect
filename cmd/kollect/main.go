package main

import (
	"bufio"
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/azure"
	"github.com/michaelcade/kollect/pkg/gcp"
	"github.com/michaelcade/kollect/pkg/kollect"
	"github.com/michaelcade/kollect/pkg/terraform"
	"github.com/michaelcade/kollect/pkg/veeam"
	"golang.org/x/term"
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
	output := flag.String("output", "", "Output file to save the collected data")
	inventoryType := flag.String("inventory", "", "Type of inventory to collect (kubernetes/aws/azure/gcp/veeam/terraform)")
	baseURL := flag.String("veeam-url", "", "Veeam server URL")
	username := flag.String("veeam-username", "", "Veeam username")
	password := flag.String("veeam-password", "", "Veeam password")
	terraformStateFile := flag.String("terraform-state", "", "Path to a local Terraform state file")
	terraformS3Bucket := flag.String("terraform-s3", "", "S3 bucket containing Terraform state (format: bucket/key)")
	terraformS3Region := flag.String("terraform-s3-region", "", "AWS region for S3 bucket (defaults to AWS_REGION env var)")
	terraformAzureContainer := flag.String("terraform-azure", "", "Azure storage container (format: storageaccount/container/blob)")
	terraformGCSBucket := flag.String("terraform-gcs", "", "GCS bucket and object (format: bucket/object)")
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

	if *browser && *inventoryType == "" && *output == "" {
		fmt.Println("Starting browser interface. Use the import function to load data.")
		startWebServer(map[string]interface{}{}, true, "", "", "")
		return
	}

	if *inventoryType == "" {
		fmt.Println("Error: You must specify an inventory type with --inventory")
		fmt.Println("Available inventory types: kubernetes, aws, azure, gcp, veeam, terraform")
		fmt.Println("Or use --browser alone to start web interface for importing data")
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
		data, err = collectData(ctx, *storageOnly, *kubeconfig)
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
					region = "us-east-1" // Default region
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

		data, err = veeam.CollectVeeamData(ctx, *baseURL, *username, *password)
	default:
		log.Fatalf("Invalid inventory type: %s", *inventoryType)
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

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getSensitiveInput(prompt string) string {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add a newline after the password input
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

func collectData(ctx context.Context, storageOnly bool, kubeconfig string) (interface{}, error) {
	return kollect.CollectData(ctx, kubeconfig)
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

func startWebServer(initialData interface{}, openBrowser bool, baseURL, username, password string) {
	fsys, err := fs.Sub(staticFiles, "web")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(fsys))

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
			data, err = veeam.CollectVeeamData(ctx, baseURL, username, password)
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

	// Add the Terraform API endpoints
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

	log.Println("Server starting on port http://localhost:8080")
	if openBrowser {
		// Open the browser
		go func() {
			var err error
			switch runtime.GOOS {
			case "darwin":
				err = exec.Command("open", "http://localhost:8080").Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:8080").Start()
			default: // Linux and other Unix-like systems
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
