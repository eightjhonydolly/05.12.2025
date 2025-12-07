package generate_report_handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/service"
)

type LinkService interface {
	GenerateReport(batchIDs []int) ([]byte, error)
}

type GenerateReportHandler struct {
	linkService LinkService
}

func NewGenerateReportHandler(linkService service.LinkService) *GenerateReportHandler {
	return &GenerateReportHandler{
		linkService: linkService,
	}
}

func (h *GenerateReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON in generate-report request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Generating report for batches: %v", req.LinksList)
	pdfData, err := h.linkService.GenerateReport(req.LinksList)
	if err != nil {
		log.Printf("Error generating report: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Generated PDF report, size: %d bytes", len(pdfData))

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=links_report.pdf")
	w.Write(pdfData)
}
