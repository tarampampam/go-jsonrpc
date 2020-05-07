// Docs: <https://www.jsonrpc.org/specification#error_object>
package errors

import (
	"fmt"
)

// Code is error code by JSON-RPC 2.0.
type Code int

const (
	Parse          Code = -32700
	InvalidRequest Code = -32600
	MethodNotFound Code = -32601
	InvalidParams  Code = -32602
	Internal       Code = -32603
)

// Error is a wrapper for a JSON interface value.
type Error struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// String returns error code in string representation (error message).
func (code Code) String() string {
	switch code {
	case Parse: // An error occurred on the server while parsing the JSON text
		return "Parse error"
	case InvalidRequest: // The JSON sent is not a valid Request object
		return "Invalid Request"
	case MethodNotFound: // The method does not exist / is not available
		return "Method not found"
	case InvalidParams: // Invalid method parameter(s)
		return "Invalid params"
	case Internal: // Internal JSON-RPC error
		return "Internal error"
	}

	return "Unrecognized error code"
}

// New creates new error.
func New(code Code) *Error {
	return &Error{
		Code:    code,
		Message: code.String(),
	}
}

// Error implements error interface.
func (err *Error) Error() string {
	var result = fmt.Sprintf("jsonrpc: code: %d, message: %s", err.Code, err.Message)

	if err.Data != nil {
		result += fmt.Sprintf(", data: %+v", err.Data)
	}

	return result
}

func (err *Error) GetCode() int         { return int(err.Code) }
func (err *Error) GetMessage() string   { return err.Message }
func (err *Error) GetData() interface{} { return err.Data }
