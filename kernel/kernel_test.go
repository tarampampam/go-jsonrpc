package kernel

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
	rpcRequest "github.com/tarampampam/go-jsonrpc/request"
	rpcResponse "github.com/tarampampam/go-jsonrpc/response"
	rpcRouter "github.com/tarampampam/go-jsonrpc/router"
)

func TestKernel_ParseJSONToRequests(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		giveJSON      []byte
		wantIsBatch   bool
		wantLength    int
		resultCheckFn func(t *testing.T, in *[]rpcRequest.Request)
		wantErr       bool
	}{
		{
			name:        "call with positional parameters",
			giveJSON:    []byte(`{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`),
			wantIsBatch: false,
			wantLength:  1,
			resultCheckFn: func(t *testing.T, in *[]rpcRequest.Request) {
				assert.Equal(t, "2.0", (*in)[0].Version)
				assert.Equal(t, "subtract", (*in)[0].Method)
				assert.EqualValues(t, []interface{}{float64(42), float64(23)}, (*in)[0].Params)
				assert.Equal(t, 1, (*in)[0].ID)
			},
		},
		{
			name:        "call with positional parameters (id as a string)",
			giveJSON:    []byte(`{"jsonrpc": "2.0", "method": "subtract", "params": {"foo": 1}, "id": "some-uid"}`),
			wantIsBatch: false,
			wantLength:  1,
			resultCheckFn: func(t *testing.T, in *[]rpcRequest.Request) {
				assert.Equal(t, "2.0", (*in)[0].Version)
				assert.Equal(t, "subtract", (*in)[0].Method)
				assert.EqualValues(t, float64(1), (*in)[0].Params.(map[string]interface{})["foo"])
				assert.Equal(t, "some-uid", (*in)[0].ID)
			},
		},
		{
			name:        "call with positional parameters (params as an array)",
			giveJSON:    []byte(`{"jsonrpc": "2.0", "method": "subtract", "params": ["foo", "bar"]}`),
			wantIsBatch: false,
			wantLength:  1,
			resultCheckFn: func(t *testing.T, in *[]rpcRequest.Request) {
				assert.Equal(t, "2.0", (*in)[0].Version)
				assert.EqualValues(t, []interface{}{"foo", "bar"}, (*in)[0].Params)
				assert.Equal(t, "subtract", (*in)[0].Method)
				assert.Nil(t, (*in)[0].ID)
			},
		},
		{
			name:        "batch call with positional parameters",
			giveJSON:    []byte(`[{"jsonrpc": "2.0", "method": "subtract", "params": ["foo"], "id": 1}]`),
			wantIsBatch: true,
			wantLength:  1,
			resultCheckFn: func(t *testing.T, in *[]rpcRequest.Request) {
				assert.Equal(t, "2.0", (*in)[0].Version)
				assert.EqualValues(t, "foo", (*in)[0].Params.([]interface{})[0])
				assert.Equal(t, "subtract", (*in)[0].Method)
				assert.Equal(t, 1, (*in)[0].ID)
			},
		},
		{
			name:        "call with invalid JSON",
			giveJSON:    []byte(`{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`),
			wantIsBatch: false,
			wantErr:     true,
		},
		{
			name:        "batch call with invalid JSON",
			giveJSON:    []byte(`[{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]]`),
			wantIsBatch: true,
			wantErr:     true,
		},
		{
			name:        "call with an empty Array",
			giveJSON:    []byte(`[]`),
			wantIsBatch: true,
			wantLength:  0,
		},
		{
			name:        "call with an invalid Batch (but not empty)",
			giveJSON:    []byte(`[1]`),
			wantIsBatch: true,
			wantLength:  1,
		},
		{
			name:        "call with invalid Batch",
			giveJSON:    []byte(`[1,2,3]`),
			wantIsBatch: true,
			wantLength:  3,
		},
		{
			name:        "call with regular and wrong requests",
			giveJSON:    []byte(`[{"jsonrpc": "2.0", "method": "notify_hello", "params": [7]}, {"foo": "boo"}]`),
			wantIsBatch: true,
			wantLength:  2,
			resultCheckFn: func(t *testing.T, in *[]rpcRequest.Request) {
				assert.Equal(t, "2.0", (*in)[0].Version)
				assert.Equal(t, "notify_hello", (*in)[0].Method)
				assert.Nil(t, (*in)[0].ID)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			kernel := New(rpcRouter.New())
			gotRequests, gotIsBatch, err := kernel.ParseJSONToRequests(tt.giveJSON)

			assert.Equal(t, gotIsBatch, tt.wantIsBatch)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, gotRequests)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, gotRequests)

				assert.Len(t, *gotRequests, tt.wantLength)

				if tt.resultCheckFn != nil {
					tt.resultCheckFn(t, gotRequests)
				}
			}
		})
	}
}

func TestKernel_HandleJSONRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		giveMethods      []jsonrpc.Method
		giveJSON         string
		giveErrorHandler ErrorHandler
		resultCheckFn    func(t *testing.T, out []byte)
		wantResultJSON   string
	}{
		{
			name:           "empty string",
			giveJSON:       "",
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}}`,
		},
		{
			name:           "rpc call with an empty Array",
			giveJSON:       "[]",
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}}`,
		},
		{
			name:     "rpc call with an invalid Batch (but not empty)",
			giveJSON: "[1]",
			wantResultJSON: `[{"jsonrpc": "2.0", "error": {
								"code": -32600, "message": "Invalid Request", "data": "wrong version"
							}}]`,
		},
		{
			name:     "rpc call with invalid Batch",
			giveJSON: `[1, 2 ,3]`,
			wantResultJSON: `[
								{"jsonrpc": "2.0", "error": {
									"code": -32600, "message": "Invalid Request", "data": "wrong version"
								}},
								{"jsonrpc": "2.0", "error": {
									"code": -32600, "message": "Invalid Request", "data": "wrong version"
								}},
								{"jsonrpc": "2.0", "error": {
									"code": -32600, "message": "Invalid Request", "data": "wrong version"
								}}
							]`,
		},
		{
			name:           "call for unknown method (with ID)",
			giveJSON:       `{"jsonrpc": "2.0", "method": "unknown", "params": true, "id": "1"}`,
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		},
		{
			name:           "call for unknown method (without ID)",
			giveJSON:       `{"jsonrpc": "2.0", "method": "unknown", "params": true, "id": null}`,
			wantResultJSON: "",
		},
		{
			name:           "rpc call with positional parameters",
			giveJSON:       `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			giveMethods:    []jsonrpc.Method{&subtractMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "result": 19, "id": 1}`,
		},
		{
			name:           "rpc call with positional parameters (2)",
			giveJSON:       `{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			giveMethods:    []jsonrpc.Method{&subtractMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "result": -19, "id": 2}`,
		},
		{
			name:           "rpc call 'subtract' method with invalid params",
			giveJSON:       `{"jsonrpc": "2.0", "method": "subtract", "params": [42], "id": 1}`,
			giveMethods:    []jsonrpc.Method{&subtractMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32602, "message": "Invalid params"}, "id": 1}`,
		},
		{
			name:           "rpc call method that returns nothing",
			giveJSON:       `{"jsonrpc": "2.0", "method": "nothing", "params": [42], "id": 1}`,
			giveMethods:    []jsonrpc.Method{&nothingMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "result": null, "id": 1}`,
		},
		{
			name:           "rpc call 'subtract' method with invalid params type",
			giveJSON:       `{"jsonrpc": "2.0", "method": "subtract", "params": "42, 23", "id": 1}`,
			giveMethods:    []jsonrpc.Method{&subtractMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32602, "message": "Invalid params"}, "id": 1}`,
		},
		{
			name: "rpc call with named parameters",
			giveJSON: `{
							"jsonrpc": "2.0",
							"method": "subtract_object",
							"params": {"first": 42, "second": 23},
							"id": "one"
						}`,
			giveMethods:    []jsonrpc.Method{&subtractObjectMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "result": {"result": 19}, "id": "one"}`,
		},
		{
			name:           "a Notification",
			giveJSON:       `{"jsonrpc": "2.0", "method": "subtract", "params": [1,2]}`,
			giveMethods:    []jsonrpc.Method{&subtractObjectMethod{}},
			wantResultJSON: "",
		},
		{
			name:        "rpc call with invalid Request object",
			giveJSON:    `{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			giveMethods: []jsonrpc.Method{&subtractObjectMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "error": {
								"code": -32600, "message": "Invalid Request", "data": "empty method"
							}}`,
		},
		{
			name: "rpc call Batch, invalid JSON",
			giveJSON: `[
							{"jsonrpc": "2.0", "method": "sum", "params": [1,2,4], "id": "1"},
							{"jsonrpc": "2.0", "method"
						]`,
			giveMethods:    []jsonrpc.Method{&subtractObjectMethod{}},
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}}`,
		},
		{
			name: "rpc call Batch",
			giveJSON: `[
							{"jsonrpc": "2.0", "method": "subtract_object", "params": {"first": 42, "second": 23}, "id": "1"},
							{"jsonrpc": "2.0", "method": "subtract", "params": [42,23], "id": "2"},
							{"jsonrpc": "2.0", "method": "subtract", "params": [3,2]},
							{"foo": "boo"},
							{"jsonrpc": "2.0", "method": "foo.get", "params": {"name": "myself"}, "id": "5"},
							{"jsonrpc": "2.0", "method": "get_data", "id": "9"}
						]`,
			giveMethods: []jsonrpc.Method{&subtractMethod{}, &subtractObjectMethod{}, &getDataMethod{}},
			resultCheckFn: func(t *testing.T, out []byte) {
				responses := make([]rpcResponse.Response, 0)
				assert.NoError(t, json.Unmarshal(out, &responses))
				assert.Len(t, responses, 5)

				var found int

				for _, response := range responses {
					if response.ID == "1" {
						assert.Equal(t, 19, int(response.Result.(map[string]interface{})["result"].(float64)))
						found++
					}
					if response.ID == "2" {
						assert.Equal(t, 19, int(response.Result.(float64)))
						found++
					}
					if response.ID == "9" {
						assert.Equal(t, "hello", response.Result.([]interface{})[0].(string))
						assert.Equal(t, 5, int(response.Result.([]interface{})[1].(float64)))
						found++
					}
					if response.ID == "5" {
						assert.Equal(t, rpcErrors.MethodNotFound, response.Error.Code)
						found++
					}
					if response.Error != nil && response.Error.Code == rpcErrors.InvalidRequest {
						found++
					}
				}

				assert.Equal(t, 5, found)
			},
		},
		{
			name:     "using errors handler",
			giveJSON: `{"jsonrpc": "2.0", "method": "nothing", "params": 1, "id": 1}`,
			giveErrorHandler: func(err jsonrpc.Error) *rpcErrors.Error {
				return &rpcErrors.Error{
					Code:    123,
					Message: "foo",
					Data:    "bar",
				}
			},
			wantResultJSON: `{"jsonrpc": "2.0", "error": {"code":123, "data":"bar", "message":"foo"}, "id": 1}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			router := rpcRouter.New()

			if tt.giveMethods != nil {
				for _, method := range tt.giveMethods {
					assert.NoError(t, router.RegisterMethod(method))
				}
			}

			kernel := New(router)

			if tt.giveErrorHandler != nil {
				kernel.InvokingErrorHandler = tt.giveErrorHandler
			}

			result := kernel.HandleJSONRequest([]byte(tt.giveJSON))

			if tt.resultCheckFn != nil {
				tt.resultCheckFn(t, result)
			}

			if tt.wantResultJSON != "" {
				assert.JSONEq(t, tt.wantResultJSON, string(result))
			}
		})
	}
}
