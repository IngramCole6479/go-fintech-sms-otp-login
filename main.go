package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	c := newClient()
	http.HandleFunc("/otp", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			To string `json:"to"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.To == "" {
			http.Error(w, "to is required", 400)
			return
		}
		if err := c.otp(in.To, "login-"+in.To); err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "code sent"})
	})
	http.HandleFunc("/risk-action", func(w http.ResponseWriter, r *http.Request) {
		var in loginRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", 400)
			return
		}
		msg, err := authorizeRiskAction(in, c.verify)
		if err != nil {
			http.Error(w, err.Error(), 403)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": msg})
	})
	log.Println("SMS login service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
