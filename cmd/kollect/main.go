package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	"github.com/michaelcade/kollect/pkg/chatbot"
	"github.com/michaelcade/kollect/pkg/kollect"
	"github.com/michaelcade/kollect/pkg/veeam"
	"golang.org/x/term"
)

var (
	dataMutex sync.Mutex
	data      interface{}
)

func main() {
	storageOnly := flag.Bool("storage", false, "Collect only storage-related objects (Kubernetes Only)")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	browser := flag.Bool("browser", false, "Open the web interface in a browser")
	output := flag.String("output", "", "Output file to save the collected data")
	inventoryType := flag.String("inventory", "kubernetes", "Type of inventory to collect (kubernetes/aws/azure/veeam)")
	baseURL := flag.String("veeam-url", "", "Veeam server URL")
	username := flag.String("veeam-username", "", "Veeam username")
	password := flag.String("veeam-password", "", "Veeam password")
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

	ctx := context.Background()

	var err error
	switch *inventoryType {
	case "aws":
		data, err = aws.CollectAWSData(ctx)
	case "azure":
		data, err = azure.CollectAzureData(ctx)
	case "kubernetes":
		data, err = collectData(ctx, *storageOnly, *kubeconfig)
	case "veeam":
		// Load environment variables for Veeam
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

		// Ensure the baseURL includes the protocol scheme
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
	} else {
		printData(data)
	}
}

func promptUser(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func getSensitiveInput(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Move to the next line after password input
	if err != nil {
		log.Fatalf("Error reading password: %v", err)
	}
	return strings.TrimSpace(string(bytePassword))
}

func getEnv(envVar, prompt string) string {
	value := os.Getenv(envVar)
	if value == "" {
		value = promptUser(prompt)
	}
	return value
}

func collectData(ctx context.Context, storageOnly bool, kubeconfig string) (interface{}, error) {
	if storageOnly {
		return kollect.CollectStorageData(ctx, kubeconfig)
	}
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

func startWebServer(data interface{}, openBrowser bool, baseURL, username, password string) {
	// Serve the files from the web directory
	fileServer := http.FileServer(http.Dir("web"))

	// Serve the files
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
		case "google":
			// Placeholder for Google Cloud data collection
			data = map[string]string{"message": "Google Cloud data collection not implemented yet"}
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
	http.HandleFunc("/api/chatbot", chatbot.HandleChatbotRequest)

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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
