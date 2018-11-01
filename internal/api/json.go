package api

import (
	"encoding/json"
	"net/http"
)

func renderJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	responseBody, err := json.Marshal(payload)

	if err != nil {
		renderJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, responseBody, statusCode)
}

func renderJSONError(w http.ResponseWriter, msg string, statusCode int) {
	payload := make(map[string]string)
	payload["error"] = msg

	responseBody, _ := json.Marshal(payload)

	writeJSON(w, responseBody, statusCode)
}

func writeJSON(w http.ResponseWriter, body []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	w.Write(body)
}
