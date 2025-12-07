package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/repository"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
	"github.com/jung-kurt/gofpdf"
)

type LinkService interface {
	CheckLinks(ctx context.Context, urls []string) (*model.LinkBatch, error)
	GenerateReport(batchIDs []int) ([]byte, error)
}

type linkService struct {
	repo   repository.LinkRepository
	client *http.Client
}

func NewLinkService(repo repository.LinkRepository) LinkService {
	return &linkService{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *linkService) CheckLinks(ctx context.Context, urls []string) (*model.LinkBatch, error) {
	batch := &model.LinkBatch{
		ID:        s.repo.GetNextID(),
		Links:     make([]model.LinkCheck, len(urls)),
		CreatedAt: time.Now(),
	}

	for i, url := range urls {
		status := s.checkURL(ctx, url)
		batch.Links[i] = model.LinkCheck{
			URL:       url,
			Status:    status,
			CheckedAt: time.Now(),
		}
	}

	if err := s.repo.SaveBatch(batch); err != nil {
		return nil, fmt.Errorf("failed to save batch: %w", err)
	}

	return batch, nil
}

func (s *linkService) checkURL(ctx context.Context, url string) model.LinkStatus {
	originalURL := url
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", originalURL, err)
		return model.StatusNotAvailable
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("Failed to check URL %s: %v", originalURL, err)
		return model.StatusNotAvailable
	}
	defer resp.Body.Close()

	log.Printf("URL %s returned status %d", originalURL, resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return model.StatusAvailable
	}

	return model.StatusNotAvailable
}

func (s *linkService) GenerateReport(batchIDs []int) ([]byte, error) {
	batches, err := s.repo.GetBatches(batchIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get batches: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, batch := range batches {
		pdf.Cell(40, 10, fmt.Sprintf("Batch ID: %d", batch.ID))
		pdf.Ln(8)

		for _, link := range batch.Links {
			status := "Available"
			if link.Status == model.StatusNotAvailable {
				status = "Not Available"
			}
			pdf.Cell(40, 10, fmt.Sprintf("%s - %s", link.URL, status))
			pdf.Ln(6)
		}
		pdf.Ln(4)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}
