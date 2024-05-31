// cmd/kollect/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/michaelcade/kollect/pkg/kollect"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/path/to/your/kubeconfig")
	if err != nil {
		log.Fatalf("Failed to build kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		data, err := kollect.CollectData(clientset)
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
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
