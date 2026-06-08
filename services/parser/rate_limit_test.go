package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitMiddlewareBlocksSecondRequestInWindow(t *testing.T) {
	handler := (&Handler{}).ParseHandler
	middleware := newRateLimitMiddleware(rateLimitConfig{
		enabled: true,
		limit:   1,
		window:  60 * time.Second,
	})
	wrapped := middleware(http.HandlerFunc(handler))
	body := []byte(`{"text":["M393A4K40DB3-CWE","S/N: SN-12345"]}`)

	firstRequest := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader(body))
	firstRequest.RemoteAddr = "127.0.0.1:5000"
	firstResponse := httptest.NewRecorder()
	wrapped.ServeHTTP(firstResponse, firstRequest)

	if firstResponse.Code != http.StatusOK {
		t.Fatalf("expected first request status 200, got %d", firstResponse.Code)
	}

	secondRequest := httptest.NewRequest(http.MethodPost, "/parse", bytes.NewReader(body))
	secondRequest.RemoteAddr = "127.0.0.1:5001"
	secondResponse := httptest.NewRecorder()
	wrapped.ServeHTTP(secondResponse, secondRequest)

	if secondResponse.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request status 429, got %d", secondResponse.Code)
	}
}

func TestRateLimitMiddlewareSkipsHealthEndpoint(t *testing.T) {
	handler := (&Handler{}).HealthHandler
	middleware := newRateLimitMiddleware(rateLimitConfig{
		enabled: true,
		limit:   1,
		window:  60 * time.Second,
	})
	wrapped := middleware(http.HandlerFunc(handler))

	firstRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	firstRequest.RemoteAddr = "127.0.0.1:5000"
	firstResponse := httptest.NewRecorder()
	wrapped.ServeHTTP(firstResponse, firstRequest)

	secondRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	secondRequest.RemoteAddr = "127.0.0.1:5001"
	secondResponse := httptest.NewRecorder()
	wrapped.ServeHTTP(secondResponse, secondRequest)

	if firstResponse.Code != http.StatusOK || secondResponse.Code != http.StatusOK {
		t.Fatalf("expected health endpoint to bypass rate limit, got %d and %d", firstResponse.Code, secondResponse.Code)
	}
}
