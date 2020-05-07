package kernel

import (
	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

type (
	subtractMethod       struct{}
	subtractMethodParams []int
)

func (*subtractMethod) GetParamsType() interface{} { return &subtractMethodParams{} }
func (*subtractMethod) GetName() string            { return "subtract" }
func (*subtractMethod) Handle(params interface{}) (interface{}, jsonrpc.Error) {
	if _, ok := params.(*subtractMethodParams); !ok {
		return nil, rpcErrors.New(rpcErrors.InvalidParams)
	}

	p := *params.(*subtractMethodParams)

	if len(p) <= 1 {
		return nil, rpcErrors.New(rpcErrors.InvalidParams)
	}

	return p[0] - p[1], nil
}

type (
	subtractObjectMethod       struct{}
	subtractObjectMethodParams struct {
		First  int `json:"first"`
		Second int `json:"second"`
	}
	subtractObjectMethodResult struct {
		Result int `json:"result"`
	}
)

func (*subtractObjectMethod) GetParamsType() interface{} { return &subtractObjectMethodParams{} }
func (*subtractObjectMethod) GetName() string            { return "subtract_object" }
func (*subtractObjectMethod) Handle(params interface{}) (interface{}, jsonrpc.Error) {
	if _, ok := params.(*subtractObjectMethodParams); !ok {
		return nil, rpcErrors.New(rpcErrors.InvalidParams)
	}

	p := *params.(*subtractObjectMethodParams)

	return subtractObjectMethodResult{
		Result: p.First - p.Second,
	}, nil
}

type (
	nothingMethod struct{}
)

func (*nothingMethod) GetParamsType() interface{} { return nil }
func (*nothingMethod) GetName() string            { return "nothing" }
func (*nothingMethod) Handle(_ interface{}) (interface{}, jsonrpc.Error) {
	return nil, nil
}

type (
	getDataMethod       struct{}
	getDataMethodResult []interface{}
)

func (*getDataMethod) GetParamsType() interface{} { return nil }
func (*getDataMethod) GetName() string            { return "get_data" }
func (*getDataMethod) Handle(params interface{}) (interface{}, jsonrpc.Error) {
	return getDataMethodResult{"hello", 5}, nil
}
