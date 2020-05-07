package kernel

import (
	"math"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
	rpcRequest "github.com/tarampampam/go-jsonrpc/request"
	rpcResponse "github.com/tarampampam/go-jsonrpc/response"
)

type ErrorHandler func(err jsonrpc.Error) *rpcErrors.Error

type Kernel struct {
	router               jsonrpc.Router
	json                 jsoniter.API
	InvokingErrorHandler ErrorHandler
}

// DefaultErrorHandler just proxy error interface into error struct.
var DefaultErrorHandler = func(err jsonrpc.Error) *rpcErrors.Error { //nolint:gochecknoglobals
	if err != nil {
		return &rpcErrors.Error{
			Code:    rpcErrors.Code(err.GetCode()),
			Message: err.GetMessage(),
			Data:    err.GetData(),
		}
	}

	return nil
}

// New creates new kernel for a working with PRC requests using Router.
func New(router jsonrpc.Router) *Kernel {
	return &Kernel{
		router:               router,
		json:                 jsoniter.ConfigFastest,
		InvokingErrorHandler: DefaultErrorHandler,
	}
}

// HandleJSONRequest accepts json request and returns processed json response.
func (kernel *Kernel) HandleJSONRequest(inJSON []byte) []byte {
	responses := rpcResponse.NewResponses()

	// parse incoming json string into requests
	requests, isBatch, parseErr := kernel.ParseJSONToRequests(inJSON)

	// and in parsing fails - push error about this into responses stack
	if parseErr != nil {
		responses.Add(rpcResponse.Response{Version: jsonrpc.Version, Error: rpcErrors.New(rpcErrors.Parse)})

		isBatch = false
	} else {
		// empty batch request cannot be processed
		if isBatch && len(*requests) == 0 {
			responses.Add(rpcResponse.Response{Version: jsonrpc.Version, Error: rpcErrors.New(rpcErrors.InvalidRequest)})

			isBatch = false
		} else {
			wg := sync.WaitGroup{}

			// loop over all passed requests
			for _, request := range *requests {
				wg.Add(1)

				// execute request processing using goroutines
				go func(request rpcRequest.Request) {
					if response := kernel.processRequest(request); response != nil {
						responses.Add(*response)
					}

					wg.Done()
				}(request)
			}

			wg.Wait()
		}
	}

	var result []byte

	if isBatch {
		result, _ = kernel.json.Marshal(responses.Items)
	} else if len(responses.Items) == 1 { // @todo: ` && responses.Items[0].Result != nil` ???
		// if request was NOT batch - only one response should be in responses stack
		result, _ = kernel.json.Marshal(responses.Items[0])
	}

	return result
}

// processRequest accepts PRC request, invoke it (if it can be invoked) and return response on success or error.
// Notifications will be processed without response returning.
func (kernel *Kernel) processRequest(request rpcRequest.Request) *rpcResponse.Response {
	// for valid request we do
	if validationErr := request.Validate(); validationErr != nil {
		err := rpcErrors.New(rpcErrors.InvalidRequest)
		err.Data = validationErr.Error()

		invalidRequestErr := rpcResponse.Response{
			Version: jsonrpc.Version,
			Error:   err,
		}

		// for request with ID we must set response ID
		if request.ID != nil {
			invalidRequestErr.ID = request.ID
		}

		return &invalidRequestErr
	}

	// method invoking with error handling
	result, invokeErr := kernel.router.Invoke(request.Method, request.Params)

	// if request has ID (it was NOT notification)
	if request.ID != nil {
		// and error was not occurred
		if invokeErr == nil {
			// push method result into responses stack (result as pointer is important)
			return &rpcResponse.Response{Version: jsonrpc.Version, Result: &result, ID: request.ID}
		}

		// on error - push error response into responses stack
		return &rpcResponse.Response{
			Version: jsonrpc.Version,
			Error:   kernel.InvokingErrorHandler(invokeErr),
			ID:      request.ID,
		}
	}

	return nil
}

// ParseJSONToRequests accepts json string and convert it into requests slice.
func (kernel *Kernel) ParseJSONToRequests(inJSON []byte) (requests *[]rpcRequest.Request, isBatch bool, err error) {
	var (
		batch  = make([]interface{}, 0)
		single interface{}
	)

	// string (as slice of bytes) must be without any whitespaces at the starting
	if len(inJSON) > 0 && inJSON[0] == byte('[') {
		isBatch = true
		err = kernel.json.Unmarshal(inJSON, &batch)
	} else {
		err = kernel.json.Unmarshal(inJSON, &single)
	}

	if err != nil {
		return
	}

	result := make([]rpcRequest.Request, 0)
	requests = &result

	if isBatch {
		for _, request := range batch {
			result = append(result, kernel.parseRawRequest(request))
		}
	} else {
		result = append(result, kernel.parseRawRequest(single))
	}

	return
}

// parseRawRequest accepts something that must be an RPC request (in "raw" structure) and returns Request object.
// Request can be invalid!
func (kernel *Kernel) parseRawRequest(in interface{}) rpcRequest.Request {
	result := rpcRequest.Request{}

	if input, ok := in.(map[string]interface{}); ok {
		if property, propOk := input["jsonrpc"]; propOk {
			if version, ok := property.(string); ok {
				result.Version = version
			}
		}

		if property, propOk := input["method"]; propOk {
			if method, ok := property.(string); ok {
				result.Method = method
			}
		}

		if property, propOk := input["params"]; propOk {
			switch params := property.(type) {
			case []interface{}, map[string]interface{}:
				result.Params = params
			}
		}

		if property, propOk := input["id"]; propOk {
			switch id := property.(type) {
			case string:
				result.ID = id
			case int:
				result.ID = id
			case int64:
				result.ID = int(id)
			case float64:
				if value, fraction := math.Modf(id); fraction == 0 {
					result.ID = int(value)
				}
			}
		}
	}

	return result
}
