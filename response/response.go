package response

import (
	"sync"

	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

type (
	// Response represents a JSON-RPC response returned by the server.
	Response struct {
		Version string           `json:"jsonrpc"`          // required, string "2.0" only
		Result  interface{}      `json:"result,omitempty"` // required only for success response (not for error)
		Error   *rpcErrors.Error `json:"error,omitempty"`  // required only for error response (not for success)
		ID      interface{}      `json:"id,omitempty"`     // optional for errors, string|int
	}

	// Responses collects set of responses
	Responses struct {
		mutex sync.RWMutex
		Items []Response
	}
)

// NewResponses creates new responses collection
func NewResponses() *Responses {
	return &Responses{
		mutex: sync.RWMutex{},
		Items: make([]Response, 0),
	}
}

// Add appends single response into collection with locking.
func (r *Responses) Add(response Response) {
	r.mutex.Lock()
	r.add(response)
	r.mutex.Unlock()
}

func (r *Responses) add(response Response) {
	r.Items = append(r.Items, response)
}
