package jsonrpc

import (
	"errors"

	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

type nothingMethod struct{}

func (*nothingMethod) GetParamsType() interface{}                        { return nil }
func (*nothingMethod) GetName() string                                   { return "nothing" }
func (*nothingMethod) Handle(_ interface{}) (interface{}, jsonrpc.Error) { return 1, nil }

type unnamedMethod struct{}

func (*unnamedMethod) GetParamsType() interface{}                        { return nil }
func (*unnamedMethod) GetName() string                                   { return "" }
func (*unnamedMethod) Handle(_ interface{}) (interface{}, jsonrpc.Error) { return nil, nil }

type erroredMethod struct{}

func (*erroredMethod) GetParamsType() interface{} { return nil }
func (*erroredMethod) GetName() string            { return "nothing" }
func (*erroredMethod) Handle(_ interface{}) (interface{}, jsonrpc.Error) {
	return nil, rpcErrors.New(rpcErrors.Code(1))
}

type (
	withParamsValidationMethod struct{}
	validatableParams          struct {
		Value bool `json:"my_value"`
	}
)

func (p *validatableParams) Validate() error {
	if !p.Value {
		return errors.New("foo")
	}

	return nil
}
func (*withParamsValidationMethod) GetParamsType() interface{} { return &validatableParams{} }
func (*withParamsValidationMethod) GetName() string            { return "validate" }
func (*withParamsValidationMethod) Handle(params interface{}) (interface{}, jsonrpc.Error) {
	return true, nil
}
