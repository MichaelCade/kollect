package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/michaelcade/kollect/pkg/kollect"
)

func main() {
	storageOnly := flag.Bool("storage", false, "Collect only storage-related objects")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Path to the kubeconfig file")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		var err error

		if *storageOnly {
			data, err = kollect.CollectStorageData(*kubeconfig)
		} else {
			data, err = kollect.CollectData(*kubeconfig)
		}

		if err != nil {
			log.Printf("Error collecting data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("Error encoding data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
