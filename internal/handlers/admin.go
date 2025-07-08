package handlers

import (
	"encoding/json"
	"lkgiovani_go_rinha_2025/internal/http"
)

func HandleAdminPaymentsSummary(hc *http.HttpCodec) {
	response := map[string]interface{}{
		"total_payments": 0,
		"total_amount":   0.0,
		"summary":        []interface{}{},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to marshal response")
		return
	}

	http.WriteOKResponse(hc, "application/json", jsonData)
}
