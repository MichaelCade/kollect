package chatbot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type ChatbotRequest struct {
	Query string `json:"query"`
}

type ChatbotResponse struct {
	Response string `json:"response"`
}

func HandleChatbotRequest(w http.ResponseWriter, r *http.Request) {
	var req ChatbotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := ProcessQuery(req.Query)
	if err != nil {
		log.Printf("Error processing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ChatbotResponse{Response: response}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func ProcessQuery(query string) (string, error) {
	// Fetch data from /api/data endpoint
	resp, err := http.Get("http://localhost:8080/api/data")
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error fetching data: %s", string(body))
		return "", fmt.Errorf("failed to fetch data: %s", string(body))
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding data: %v", err)
		return "", err
	}

	// Combine data into a single prompt
	prompt := "Here is the collected data:\n"
	prompt += "Data: " + formatData(data) + "\n"
	prompt += "Query: " + query

	// Generate response from the prompt using Ollama CLI
	ctx := context.Background()
	completion, err := generateResponseWithOllama(ctx, prompt)
	if err != nil {
		log.Printf("Error generating response with Ollama: %v", err)
		return "", err
	}

	return completion, nil
}

func generateResponseWithOllama(ctx context.Context, prompt string) (string, error) {
	// Write the prompt to a temporary file
	tmpfile, err := ioutil.TempFile("", "ollama-prompt-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(prompt); err != nil {
		tmpfile.Close()
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}
	tmpfile.Close()

	// Pass the temporary file to the Ollama CLI
	cmd := exec.CommandContext(ctx, "ollama", "run", "--prompt-file", tmpfile.Name())
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing Ollama CLI: %v", err)
		return "", fmt.Errorf("failed to generate response: %v", err)
	}
	return string(output), nil
}

func formatData(data interface{}) string {
	// Format the data as a string (you can customize this as needed)
	formattedData, _ := json.MarshalIndent(data, "", "  ")
	return string(formattedData)
}
