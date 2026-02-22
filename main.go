package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"url-shortner/utils"
)

type HttpRequest struct {
	Method string
	URL    string
	Header http.Header
	Body   []byte
}

type HttpResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

type createDataRequest struct {
	URL string `json:"url"`
}

func main() {
	db := utils.ConnectDB()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		req := HttpRequest{
			Method: r.Method,
			URL:    r.URL.String(),
			Header: r.Header,
			Body:   body,
		}

		switch req.Method {

		case http.MethodGet:
			retrieveData := utils.GetURLData(db, req.URL[1:]) // Remove leading slash
			http.Redirect(w, r, retrieveData, http.StatusSeeOther)

		case http.MethodPost:
			var data createDataRequest
			if err := json.Unmarshal(req.Body, &data); err != nil {
				http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
				return
			}
			createdData := utils.CreateURLData(db, data.URL)
			fmt.Println("Created data for URL:", createdData)
			fmt.Println(string(req.Body))
			resp := HttpResponse{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       fmt.Appendf(nil, "%s \n", createdData),
			}
			resp.Header.Set("Content-Type", "text/plain")
			w.WriteHeader(resp.StatusCode)
			w.Write(resp.Body)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
