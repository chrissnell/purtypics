package editor

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON serializes v as JSON with the correct Content-Type. If encoding
// fails after the header is written we can't recover the response, but we do
// log it so the failure isn't silent.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("json encode: %v", err)
	}
}
