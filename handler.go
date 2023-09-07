package main

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"
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

// embed ui files
//go:embed ui/dist/*
var uiContent embed.FS

func embeddedUIHandler(w http.ResponseWriter, r *http.Request) {
	path := "ui/dist" + r.URL.Path

	if r.URL.Path == "/" {
		path = "ui/dist/index.html"
	}

	data, err := uiContent.ReadFile(path)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(data)
	if strings.HasSuffix(r.URL.Path, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(r.URL.Path, ".css") {
		contentType = "text/css"
	}

	w.Header().Set("Content-Type", contentType)

	if _, err := w.Write(data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
