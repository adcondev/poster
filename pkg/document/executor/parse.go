package executor

import (
	"encoding/json"
	"fmt"

	"github.com/adcondev/pos-printer/pkg/document/schema"
)

// PrintJob representa los datos de un documento de impresión
type PrintJob struct {
	Data schema.Document `json:"data"`
}

// ParseDocument parsea un documento JSON
func ParseDocument(data []byte) (*schema.Document, error) {
	var doc schema.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Validación básica
	if doc.Version == "" {
		// TODO: Review an smart way to handle versioning
		doc.Version = "1.0"
	}

	if len(doc.Commands) == 0 {
		return nil, fmt.Errorf("document must contain at least one command")
	}

	return &doc, nil
}
