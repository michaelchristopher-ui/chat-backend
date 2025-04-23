package httpadapter

import (
	"io"
	"net/http"
)

// Adapter defines an interface for HTTP-related functions
//
//go:generate mockgen -source=adapter.go -package=httpadapter -destination=adapter_mock.go
type Adapter interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}
