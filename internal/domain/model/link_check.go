package model

import "time"

type LinkStatus string

const (
	StatusAvailable    LinkStatus = "available"
	StatusNotAvailable LinkStatus = "not available"
)

type LinkCheck struct {
	URL       string
	Status    LinkStatus
	CheckedAt time.Time
}

type LinkBatch struct {
	ID        int
	Links     []LinkCheck
	CreatedAt time.Time
}