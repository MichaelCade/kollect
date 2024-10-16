package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/michaelcade/kollect/pkg/aws"
	"github.com/michaelcade/kollect/pkg/kollect"
)

func main() {
	storageOnly := flag.Bool("storage", false, "Collect only storage-related objects")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	browser := flag.Bool("browser", false, "Open the web interface in a browser")
	output := flag.String("output", "", "Output file to save the collected data")
	inventoryType := flag.String("inventory", "kubernetes", "Type of inventory to collect (kubernetes/aws)")
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

	var data interface{}
	var err error

	switch *inventoryType {
	case "aws":
		data, err = aws.CollectAWSData()
	case "kubernetes":
		data, err = collectData(*storageOnly, *kubeconfig)
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

	if *browser {
		startWebServer(data)
	} else {
		printData(data)
	}
}

func collectData(storageOnly bool, kubeconfig string) (interface{}, error) {
	if storageOnly {
		return kollect.CollectStorageData(kubeconfig)
	}
	return kollect.CollectData(kubeconfig)
}

func saveToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}

func printData(data interface{}) {
	prettyData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting data: %v", err)
	}
	fmt.Println(string(prettyData))
}

func startWebServer(data interface{}) {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Error encoding data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server starting on port http://localhost:8080")

	// Open the browser
	go func() {
		err := exec.Command("open", "http://localhost:8080").Start()
		if err != nil {
			log.Fatalf("Failed to open browser: %v", err)
		}
	}()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
