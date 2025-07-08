package handlers

import (
	"encoding/json"

	"lkgiovani_go_rinha_2025/internal/http"
	"lkgiovani_go_rinha_2025/internal/models"
)

func HandleHealth(hc *http.HttpCodec) {
	response := models.HealthResponse{
		Status: "OK",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.WriteInternalServerErrorResponse(hc, "Failed to marshal response")
		return
	}

	http.WriteOKResponse(hc, "application/json", jsonData)
}
