// Run this example: `go run . &`
//
// And then: `curl -d '{"jsonrpc":"2.0","method":"my.method","params":{"number":1,"is_optional":"foo"},"id":1}' -s http://127.0.0.1:8080/rpc`
// Will return: `{"jsonrpc":"2.0","result":{"number":1,"is_optional":"foo","string":"default"},"id":1}`

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
	rpcKernel "github.com/tarampampam/go-jsonrpc/kernel"
	rpcResponse "github.com/tarampampam/go-jsonrpc/response"
	rpcRouter "github.com/tarampampam/go-jsonrpc/router"
)

type (
	// Define RPC method as a structure
	myRpcMethod struct{} // Implements `jsonrpc.Method` interface

	// Define params for method above as a structure with "json" fields
	myRpcMethodParams struct { // Implements `jsonrpc.Validator` interface
		Number     int     `json:"number"`
		IsOptional *string `json:"is_optional"` // optional value (can be nil)
		JustString string  `json:"string"`
	}
)

// Implement `jsonrpc.Validator` interface for easy incoming params validation
func (p *myRpcMethodParams) Validate() error {
	if p.Number <= 0 {
		return errors.New("number must be positive")
	}

	return nil // all is ok
}

// Define params structure for a method
func (*myRpcMethod) GetParamsType() interface{} {
	return &myRpcMethodParams{
		JustString: "default", // this is default value, and it will be overwritten if defined in incoming params (JSON)
	}
}

// Define method name
func (*myRpcMethod) GetName() string { return "my.method" }
func (*myRpcMethod) Handle(p interface{}) (interface{}, jsonrpc.Error) {
	// check passed params type
	params, ok := p.(*myRpcMethodParams)
	if !ok {
		return nil, rpcErrors.New(rpcErrors.InvalidParams)
	}

	return params, nil // send params back (as a response) without any changes
}

// httpHandler creates HTTP handler function for RPC requests handling
func httpHandler(router jsonrpc.Router) http.HandlerFunc {
	// create kernel using our router
	kernel := rpcKernel.New(router)

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(writer)

		// nil body means the request has no body, such as a GET request
		if request.Body == nil {
			_ = encoder.Encode(rpcResponse.Response{
				Version: jsonrpc.Version,
				Error:   rpcErrors.New(rpcErrors.InvalidRequest),
			})

			return
		}

		// read full request body
		body, _ := ioutil.ReadAll(request.Body)

		// and using RPC kernel handle them, and write response
		_, _ = writer.Write(kernel.HandleJSONRequest(body))
	}
}

// createRouter creates router with registered RPC methods.
func createRouter() jsonrpc.Router {
	// create router instance
	router := rpcRouter.New()

	// register our method
	if err := router.RegisterMethod(new(myRpcMethod)); err != nil {
		panic(err)
	}

	return router
}

func main() {
	// create router
	router := createRouter()

	// register our HTTP handler for RPS requests processing
	http.HandleFunc("/rpc", httpHandler(router))

	// and start HTTP server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
