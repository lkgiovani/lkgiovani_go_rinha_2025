package router

import (
	"lkgiovani_go_rinha_2025/internal/handlers"
	"lkgiovani_go_rinha_2025/internal/http"
	"strings"
)

func HandleRequest(hc *http.HttpCodec) {
	path := string(hc.GetParser().Path)
	method := string(hc.GetParser().Method)

	switch {
	case path == "/healthz":
		handlers.HandleHealth(hc)

	case strings.HasPrefix(path, "/payments-summary") && method == "GET":
		handlers.HandlePaymentsSummary(hc)

	case path == "/payments" && method == "POST":
		handlers.HandleCreatePayment(hc)

	default:
		http.WriteNotFoundResponse(hc)
	}
}
