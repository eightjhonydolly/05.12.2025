package check_links_handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/eightjhonydolly/05.12.2025/internal/domain/links/service"
	"github.com/eightjhonydolly/05.12.2025/internal/domain/model"
)

type LinkService interface {
	CheckLinks(ctx context.Context, urls []string) (*model.LinkBatch, error)
}

type CheckLinksHandler struct {
	linkService LinkService
}

func NewCheckLinksHandler(linkService service.LinkService) *CheckLinksHandler {
	return &CheckLinksHandler{
		linkService: linkService,
	}
}

func (h *CheckLinksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req CheckLinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid JSON in check-links request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Checking %d links", len(req.Links))
	batch, err := h.linkService.CheckLinks(r.Context(), req.Links)
	if err != nil {
		log.Printf("Error checking links: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Created batch %d with %d links", batch.ID, len(batch.Links))

	resp := CheckLinksResponse{
		Links:    make(map[string]string),
		LinksNum: batch.ID,
	}

	for _, link := range batch.Links {
		resp.Links[link.URL] = string(link.Status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
