package commands

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTMLHandlerRejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/html", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	htmlHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestURLHandlerRejectsInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	urlHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHTMLHandlerRequiresHTML(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/html", strings.NewReader(`{"html":""}`))
	rec := httptest.NewRecorder()

	htmlHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestURLHandlerRequiresURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(`{"url":""}`))
	rec := httptest.NewRecorder()

	urlHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
