package main

import (
	"net/http"
	"os"

	"github.com/dtsong/battleforge/backend/internal/httpapi"
	"github.com/dtsong/battleforge/backend/internal/observability"
)

func main() {
	logger := observability.NewLogger()

	addr := getAddr()
	logger.Infof("starting battleforge-api on %s", addr)

	router := httpapi.NewRouter(logger)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatalf("server failed: %v", err)
	}
}

func getAddr() string {
	if v := os.Getenv("BATTLEFORGE_API_ADDR"); v != "" {
		return v
	}
	// default dev address
	return ":8080"
}
