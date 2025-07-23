package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type HealthResponse struct {
	MinResponseTime int  `json:"minResponseTime"`
	Failing         bool `json:"failing"`
}

func CheckService(statusChan chan map[string]int) {

	default_uri := os.Getenv("PAYMENT_PROCESSOR_DEFAULT_URL")
	fallback_uri := os.Getenv("PAYMENT_PROCESSOR_FALLBACK_URL")

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	healthy := make(map[string]int)

	check := func(name, url string) {

		resp, err := client.Get(fmt.Sprintf("%s/payments/service-health", url))
		if err != nil {
			healthy[name] = -1
			return
		}
		defer resp.Body.Close()

		var data HealthResponse

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			healthy[name] = -1
			return
		}

		if data.Failing {
			healthy[name] = -1
		} else {
			healthy[name] = data.MinResponseTime
		}

	}

	check("default", default_uri)
	check("fallback", fallback_uri)

	statusChan <- healthy

}
