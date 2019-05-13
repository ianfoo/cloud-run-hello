package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}
	routes()
	logrus.WithField("port", port).Info("server listening")
	logrus.WithError(http.ListenAndServe(":"+port, nil)).Error("cannot serve")
}

func routes() {
	http.Handle("/", http.FileServer(http.Dir("pages")))
	http.HandleFunc("/ping", logAccess(ping))
	http.HandleFunc("/hello", logAccess(hello))
	http.HandleFunc("/time", logAccess(reportTime))
}

func logAccess(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		logrus.WithFields(logrus.Fields{
			"elapsed": time.Since(start),
			"path":    r.URL.Path,
		}).Info("served request")
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	http.Error(w, http.StatusText(status), status)
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "there"
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func reportTime(w http.ResponseWriter, r *http.Request) {
	formats := map[string]string{
		"ansic":       time.ANSIC,
		"unixdate":    time.UnixDate,
		"rubydate":    time.RubyDate,
		"rfc822":      time.RFC822,
		"rfc822z":     time.RFC822Z,
		"rfc850":      time.RFC850,
		"rfc1123":     time.RFC1123,
		"rfc1123z":    time.RFC1123Z,
		"rfc3339":     time.RFC3339,
		"rfc3339nano": time.RFC3339Nano,
		"kitchen":     time.Kitchen,
	}
	now := time.Now()
	format := r.URL.Query().Get("format")
	format = strings.TrimSpace(strings.ToLower(format))
	if format == "" || format == "all" {
		resp := make(map[string]string, len(formats))
		for k, v := range formats {
			resp[k] = now.Format(v)
		}
		writeJSON(w, resp)
		return
	}
	formatStr, ok := formats[format]
	if !ok {
		http.Error(w, "Invalid format", http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]string{format: now.Format(formatStr)})
}

func writeJSON(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
