package router

import (
	"errors"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/tarampampam/go-jsonrpc"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

// Router is default RPC router implementation.
type Router struct {
	mutex   sync.RWMutex
	methods map[string]jsonrpc.Method
	json    jsoniter.API
}

// New creates new router instance.
func New() *Router {
	return &Router{
		mutex:   sync.RWMutex{},
		methods: map[string]jsonrpc.Method{},
		json:    jsoniter.ConfigFastest,
	}
}

// RegisterMethod make a method registration for later invoking.
func (router *Router) RegisterMethod(method jsonrpc.Method) error {
	var methodName = method.GetName()

	if methodName == "" {
		return errors.New("jsonrpc: method name should not be empty")
	}

	router.mutex.Lock()
	router.methods[methodName] = method
	router.mutex.Unlock()

	return nil
}

// MethodIsRegistered returns `true` only if passed method is registered.
func (router *Router) MethodIsRegistered(methodName string) bool {
	router.mutex.RLock()
	_, ok := router.methods[methodName]
	router.mutex.RUnlock()

	return ok
}

// Invoke accepts method name and invoke registered method with same name. If requested method is not
// registered - error will be returned.
func (router *Router) Invoke(methodName string, params interface{}) (interface{}, jsonrpc.Error) {
	router.mutex.RLock()
	method, ok := router.methods[methodName]
	router.mutex.RUnlock()

	if !ok {
		return nil, rpcErrors.New(rpcErrors.MethodNotFound)
	}

	// this is crutch for request params binding into required structure
	methodParams := method.GetParamsType()
	if methodParams != nil {
		// pass params through "params type" object
		bytes, _ := router.json.Marshal(params)
		if err := router.json.Unmarshal(bytes, &methodParams); err != nil {
			return nil, rpcErrors.New(rpcErrors.InvalidParams)
		}

		// if params struct follows validator interface - make check using validation method
		if p, ok := methodParams.(jsonrpc.Validator); ok {
			if validationErr := p.Validate(); validationErr != nil {
				err := rpcErrors.New(rpcErrors.InvalidParams)
				err.Data = validationErr.Error()

				return nil, err
			}
		}
	}

	return method.Handle(methodParams)
}
