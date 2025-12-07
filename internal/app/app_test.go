package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_Integration_CheckLinksAndGenerateReport(t *testing.T) {
	app, err := NewApp("")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	server := httptest.NewServer(app.server.Handler)
	defer server.Close()


	checkReq := map[string]interface{}{
		"links": []string{"httpbin.org", "invalid.test"},
	}
	reqBody, _ := json.Marshal(checkReq)

	resp, err := http.Post(server.URL+"/api/check-links", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Check links request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var checkResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&checkResp); err != nil {
		t.Fatalf("Failed to decode check response: %v", err)
	}

	linksNum, ok := checkResp["links_num"].(float64)
	if !ok {
		t.Fatal("links_num not found in response")
	}


	reportReq := map[string]interface{}{
		"links_list": []int{int(linksNum)},
	}
	reqBody, _ = json.Marshal(reportReq)

	resp, err = http.Post(server.URL+"/api/generate-report", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Generate report request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("Expected Content-Type application/pdf, got %s", contentType)
	}
}

func TestApp_Integration_InvalidRequests(t *testing.T) {
	app, err := NewApp("")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}

	server := httptest.NewServer(app.server.Handler)
	defer server.Close()


	resp, err := http.Post(server.URL+"/api/check-links", "application/json", bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}


	reportReq := map[string]interface{}{
		"links_list": []int{999},
	}
	reqBody, _ := json.Marshal(reportReq)

	resp, err = http.Post(server.URL+"/api/generate-report", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 even for empty report, got %d", resp.StatusCode)
	}
}

func TestNewApp_ConfigError(t *testing.T) {

	app, err := NewApp("invalid/path")
	if err != nil {
		t.Fatalf("NewApp should not fail with invalid config path: %v", err)
	}

	if app == nil {
		t.Fatal("Expected app instance, got nil")
	}

	if app.config == nil {
		t.Fatal("Expected config to be initialized")
	}
}

func TestBootstrapHandler(t *testing.T) {
	handler := bootstrapHandler()
	if handler == nil {
		t.Fatal("Expected handler, got nil")
	}


	req := httptest.NewRequest("POST", "/api/check-links", bytes.NewReader([]byte(`{"links":["test.com"]}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Error("Handler should handle /api/check-links route")
	}
}

func TestApp_ListenAndServe_InvalidAddress(t *testing.T) {
	app, err := NewApp("")
	if err != nil {
		t.Fatalf("Failed to create app: %v", err)
	}


	app.config.Server.Host = "invalid-host-that-does-not-exist"
	app.config.Server.Port = "99999"

	err = app.ListenAndServe()
	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
}