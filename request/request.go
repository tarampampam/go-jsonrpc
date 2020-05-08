package request

import (
	"errors"

	"github.com/tarampampam/go-jsonrpc"
)

// Request represents a JSON-RPC request received by the server.
type Request struct {
	Version string      `json:"jsonrpc"`          // required, string "2.0" only
	Method  string      `json:"method"`           // required, string (any)
	Params  interface{} `json:"params,omitempty"` // optional, array|object
	ID      interface{} `json:"id"`               // optional for notifications only, string|int
}

// Validate makes request validation (request is correct and can be processed?).
func (request *Request) Validate() error {
	if request.Version != jsonrpc.Version {
		return errors.New("wrong version")
	}

	if request.Method == "" {
		return errors.New("empty method")
	}

	switch request.ID.(type) {
	case nil, string, int, int64:
		break
	default:
		return errors.New("wrong id type")
	}

	wrongParamsErr := errors.New("wrong params type")

	switch request.Params.(type) {
	case uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64:
		return wrongParamsErr
	case float32, float64, complex64, complex128:
		return wrongParamsErr
	case string, int, bool, uintptr:
		return wrongParamsErr
	case nil, []interface{}, map[interface{}]interface{}, interface{}:
		break
	default: // @todo: how to test this case?
		return wrongParamsErr
	}

	return nil
}
