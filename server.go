package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	store = make(map[string]string)
	mu    sync.RWMutex
)

type hea struct {
	Status string `json:"status"`
}

type request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func hello(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(res, "hello nigg\n")
}

func health(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	r := hea{Status: "ok"}
	if err := json.NewEncoder(res).Encode(r); err != nil {
		http.Error(res, "Something went wrong", http.StatusInternalServerError)
		return
	}

}

func echo(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "this method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body map[string]interface{}

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(body); err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}

}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(res, req)
		log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
	})

}
func setHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "this method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in request
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}
	mu.Lock()
	store[in.Key] = in.Value
	mu.Unlock()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, `{"status":"ok"}`)
}
func getHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := req.URL.Query().Get("key")
	mu.RLock()
	value, ok := store[key]
	mu.RUnlock()
	if !ok {
		http.Error(res, "key not found", http.StatusNotFound)
		return
	}
	res.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(res).Encode(map[string]string{"value": value}); err != nil {
		http.Error(res, "something went wrong", http.StatusInternalServerError)
		return
	}

}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/set", setHandler)
	mux.HandleFunc("/get", getHandler)

	if err := http.ListenAndServe(":8080", loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}

}
