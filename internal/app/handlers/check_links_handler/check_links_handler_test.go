package check_links_handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
)

type mockLinkService struct{}

func (m *mockLinkService) CheckLinks(ctx context.Context, urls []string) (*model.LinkBatch, error) {
	return &model.LinkBatch{
		ID: 1,
		Links: []model.LinkCheck{
			{URL: "google.com", Status: model.StatusAvailable},
		},
	}, nil
}

func (m *mockLinkService) GenerateReport(batchIDs []int) ([]byte, error) {
	return []byte("fake pdf"), nil
}

func TestCheckLinksHandler_ServeHTTP(t *testing.T) {
	handler := NewCheckLinksHandler(&mockLinkService{})

	reqBody := CheckLinksRequest{
		Links: []string{"google.com"},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/check-links", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp CheckLinksResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.LinksNum != 1 {
		t.Errorf("Expected LinksNum 1, got %d", resp.LinksNum)
	}

	if resp.Links["google.com"] != "available" {
		t.Errorf("Expected google.com to be available, got %s", resp.Links["google.com"])
	}
}