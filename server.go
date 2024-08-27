package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ResponseError struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func main() {
	// Setup a logger to a file
	f, err := os.OpenFile("mockhttpserver.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Println("--- Starting Mock HTTP Server ---")

	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	SetupServerHandlers()
	StartServer()
}

func SetupServerHandlers() {
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/echo/details", echoDetailsHandler)
	http.HandleFunc("/write", writeHandler)
}

func ReadBody(r *http.Request) (*string, *ResponseError) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, &ResponseError{
			Error:   "Failed to read request body",
			Details: err.Error(),
		}
	}
	bodyString := string(body)
	return &bodyString, nil
}

func MarshalJSON(v interface{}) (*string, *ResponseError) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, &ResponseError{
			Error:   "Failed to convert to JSON",
			Details: err.Error(),
		}
	}
	jsonString := string(jsonData)
	return &jsonString, nil
}

func StartServer() {
	// Get the server port from environment variables
	port := os.Getenv("MOCK_SERVER_PORT")
	if port == "" {
		port = "8080" // Default to port 8080
	}

	// Start the HTTP server
	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println("Failed to start server:", err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ReadBody(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := MarshalJSON(err)
		if err != nil {
			http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
			return
		}
		http.Error(w, *jsonData, http.StatusInternalServerError)
		return
	}

	// Echo the request body
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(*body))
}

func echoDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the request details
	type RequestDetails struct {
		Method     string                 `json:"method"`
		URL        string                 `json:"url"`
		Proto      string                 `json:"proto"`
		Host       string                 `json:"host"`
		RemoteAddr string                 `json:"remote_addr"`
		RequestURI string                 `json:"request_uri"`
		UserAgent  string                 `json:"user_agent"`
		Referer    string                 `json:"referer"`
		Header     map[string][]string    `json:"header"`
		BodyText   string                 `json:"bodyText"`
		BodyJson   map[string]interface{} `json:"bodyJson"`
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	textBody := string(bodyBytes)

	jsonBody := make(map[string]interface{})
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
			log.Printf("Failed to decode JSON body: %v", err)
			jsonBody["error"] = "Failed to decode JSON body"
			jsonBody["details"] = err.Error()
		}
	}

	// Create an instance of the RequestDetails struct
	details := RequestDetails{
		Method:     r.Method,
		URL:        r.URL.String(),
		Proto:      r.Proto,
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
		UserAgent:  r.UserAgent(),
		Referer:    r.Referer(),
		Header:     r.Header,
		BodyText:   textBody,
		BodyJson:   jsonBody,
	}

	// Convert the struct to JSON
	jsonData, err := json.Marshal(details)
	if err != nil {
		http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
		return
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonData)
}
