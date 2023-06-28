package transferobjects

import "time"

type ExceptionContainerResponse struct {
	ExceptionContainer
	HVersion     string        `json:"_version,omitempty"`
	HTags        []interface{} `json:"_tags,omitempty"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
	CreatedBy    string        `json:"created_by,omitempty"`
	TieBreakerID string        `json:"tie_breaker_id,omitempty"`
	UpdatedAt    time.Time     `json:"updated_at,omitempty"`
	UpdatedBy    string        `json:"updated_by,omitempty"`
}

type ExceptionContainer struct {
	Description   string   `json:"description,omitempty"`
	ID            string   `json:"id,omitempty"`
	ListID        string   `json:"list_id,omitempty"`
	Name          string   `json:"name,omitempty"`
	NamespaceType string   `json:"namespace_type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Type          string   `json:"type,omitempty"`
}
