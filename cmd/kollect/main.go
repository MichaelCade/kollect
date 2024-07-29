// cmd/kollect/main.go
package main

import (
	"encoding/json"
	"flag"
	"github.com/michaelcade/kollect/pkg/kollect"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	storageOnly := flag.Bool("storage", false, "Collect only storage-related objects")
	flag.Parse()

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		var err error

		if *storageOnly {
			data, err = kollect.CollectStorageData(kubeconfig)
		} else {
			data, err = kollect.CollectData(kubeconfig)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
