package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lucaslimafernandes/rinha-backend-2025-go/models"
	"github.com/lucaslimafernandes/rinha-backend-2025-go/utils"
)

type HealthResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

var (
	healthStatus = make(map[string]int)
	statusMutex  = sync.RWMutex{}
)

func updateHealthLoop() {

	for {
		ch := make(chan map[string]int)

		go utils.CheckService(ch)
		result := <-ch

		statusMutex.Lock()
		healthStatus = result
		statusMutex.Unlock()

		time.Sleep(5 * time.Second)

	}
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

	if req.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	statusMutex.RLock()
	defer statusMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	if healthStatus["default"] >= 0 {
		fmt.Fprintf(w, "default")
	}

	var payment utils.Payment

	err := json.NewDecoder(req.Body).Decode(&payment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	switch {
	case healthStatus["default"] >= 0 && healthStatus["default"] < 99:
		err = utils.PaymentSend("default", payment)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case healthStatus["default"] == -1 && healthStatus["fallback"] >= 0:
		err = utils.PaymentSend("fallback", payment)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	case healthStatus["default"] > 99 && healthStatus["fallback"] >= 0:
		err = utils.PaymentSend("fallback", payment)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	default:
		err = utils.PaymentSend("default", payment)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
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

	if os.Getenv("CREATE_SCHEMA") == "true" {
		models.CreateTable()
	}

	go updateHealthLoop()

	http.HandleFunc("/healthy", healthy)

	// http.HandleFunc("/admin/purge-payments", purgePayments)
	http.HandleFunc("/purge-payments", purgePayments)
	http.HandleFunc("/payments", payment)
	http.HandleFunc("/payments-summary", paymentsSummary)

	log.Println("porta 9999...")
	log.Fatalln(http.ListenAndServe(":9999", nil))

}
