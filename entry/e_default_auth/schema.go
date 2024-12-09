package e_default_auth

import "github.com/goccy/go-json"

// DefaultAuth mapped from table <default_auth>
type DefaultAuth struct {
	ID   int32           `json:"id"`
	Type string          `json:"type"`
	Auth json.RawMessage `json:"auth"`
}
