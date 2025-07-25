package main

import (
	"api-rinha/models"
	"api-rinha/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type HealthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func healthy(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Message: "API is healthy",
		Status:  "ok",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
	}

}

func purgePayments(w http.ResponseWriter, req *http.Request) {

	err := models.PurgeTable()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
	}

	response := HealthResponse{
		Message: "Purge payments",
		Status:  "ok",
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
	}

}

func payment(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var payment utils.Payment

	err := json.NewDecoder(req.Body).Decode(&payment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err = utils.PaymentSend(payment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// json.NewEncoder(w).Encode(payment)

}

func paymentsSummary(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fromStr := req.URL.Query().Get("from")
	toStr := req.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		http.Error(w, "Missing 'from' or 'to' parameters", http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(fromStr, "Z") {
		fromStr += "Z"
	}
	if !strings.HasSuffix(toStr, "Z") {
		toStr += "Z"
	}

	fromTime, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		http.Error(w, "Invalid 'from' datetime format", http.StatusBadRequest)
		return
	}

	toTime, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		http.Error(w, "Invalid 'to' datetime format", http.StatusBadRequest)
		return
	}

	summary, err := models.GetPaymentSummary(fromTime, toTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)

}

func main() {

	// PgBouncer
	os.Setenv("PGAPPNAME", "")
	os.Setenv("PGOPTIONS", "")

	models.DBConnect()
	models.CreateTable()

	go updateHealthLoop()

	http.HandleFunc("/healthy", healthy)

	http.HandleFunc("/admin/purge-payments", purgePayments)
	http.HandleFunc("/payments", payment)
	http.HandleFunc("/payments-summary", paymentsSummary)

	log.Println("porta 9999...")
	log.Fatalln(http.ListenAndServe(":9999", nil))

}
