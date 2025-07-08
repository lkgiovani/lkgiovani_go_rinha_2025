package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lkgiovani_go_rinha_2025/internal/http"
	nethttp "net/http"
	"strings"
)

type PaymentRequest struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

func HandlePaymentsSummary(hc *http.HttpCodec) {

	path := string(hc.GetParser().Path)

	fmt.Printf("=== DEBUG PAYMENTS SUMMARY ===\n")
	fmt.Printf("Original path received: %s\n", path)

	if queryStart := strings.Index(path, "?"); queryStart != -1 {
		queryParams := path[queryStart:]
		fmt.Printf("Query parameters found: %s\n", queryParams)

		params := strings.Split(queryParams[1:], "&")
		for _, param := range params {
			fmt.Printf("Parameter: %s\n", param)
		}
	} else {
		fmt.Printf("No query parameters found in path\n")
	}

	fmt.Printf("=== END DEBUG ===\n")

	response := map[string]interface{}{
		"default": map[string]interface{}{
			"totalAmount":   0,
			"totalRequests": 0,
		},
		"fallback": map[string]interface{}{
			"totalAmount":   0,
			"totalRequests": 0,
		},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
		http.WriteInternalServerErrorResponse(hc, "Failed to marshal response")
		return
	}

	http.WriteOKResponse(hc, "application/json", jsonData)
}

func HandleCreatePayment(hc *http.HttpCodec) {
	body := hc.GetBody()

	var paymentRequest PaymentRequest
	if err := json.Unmarshal(body, &paymentRequest); err != nil {
		http.WriteInternalServerErrorResponse(hc, "Invalid JSON body")
		return
	}

	paymentData, err := json.Marshal(paymentRequest)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to marshal payment data")
		return
	}

	req, err := nethttp.NewRequest("POST", "http://localhost:8001/payments", bytes.NewBuffer(paymentData))
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Rinha-Token", "123")

	client := &nethttp.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to make request to payment processor")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to read response")
		return
	}

	fmt.Printf("Payment processor response: %s\n", string(respBody))

	response := map[string]interface{}{
		"message":            "payment processed successfully",
		"correlationId":      paymentRequest.CorrelationId,
		"amount":             paymentRequest.Amount,
		"processor_response": string(respBody),
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to marshal response")
		return
	}

	http.WriteOKResponse(hc, "application/json", jsonData)
}
