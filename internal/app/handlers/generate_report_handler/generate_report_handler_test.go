package generate_report_handler

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

func (m *mockLinkService) GenerateReport(batchIDs []int) ([]byte, error) {
	return []byte("fake pdf data"), nil
}

func (m *mockLinkService) CheckLinks(ctx context.Context, urls []string) (*model.LinkBatch, error) {
	return nil, nil
}

func TestGenerateReportHandler_ServeHTTP(t *testing.T) {
	handler := NewGenerateReportHandler(&mockLinkService{})

	reqBody := GenerateReportRequest{
		LinksList: []int{1, 2},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/generate-report", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("Expected Content-Type application/pdf, got %s", contentType)
	}

	if w.Body.Len() == 0 {
		t.Error("Expected PDF data, got empty response")
	}
}