package main

import (
	"encoding/json"
	"net/http"
)

func inviteDataApiHandler(w http.ResponseWriter, re *http.Request) {
    jsonBytes, err := json.Marshal(whitelist)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
    if _, err := w.Write(jsonBytes); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}	
}
