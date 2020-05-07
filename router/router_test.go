package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

func TestRouter_MethodRegistration(t *testing.T) {
	t.Parallel()

	router := New()
	method := &nothingMethod{}

	assert.Nil(t, router.RegisterMethod(method))
	assert.True(t, router.MethodIsRegistered(method.GetName()))

	assert.Contains(t, router.RegisterMethod(new(unnamedMethod)).Error(), "not be empty")
}

func TestRouter_Invoke(t *testing.T) {
	t.Parallel()

	router := New()
	method := &nothingMethod{}

	assert.Nil(t, router.RegisterMethod(method))

	res, err := router.Invoke(method.GetName(), nil)

	assert.Nil(t, err)
	assert.Equal(t, 1, res)
}

func TestRouter_InvokeUnknownMethod(t *testing.T) {
	t.Parallel()

	res, err := New().Invoke("foo", nil)

	assert.Nil(t, res)
	assert.Equal(t, int(rpcErrors.MethodNotFound), err.GetCode())
}

func TestRouter_InvokeWithMethodError(t *testing.T) {
	t.Parallel()

	router := New()
	method := &erroredMethod{}

	assert.Nil(t, router.RegisterMethod(method))

	res, err := router.Invoke(method.GetName(), nil)

	assert.Nil(t, res)
	assert.Equal(t, 1, err.GetCode())
}

func TestRouter_InvokeWithParamsValidation(t *testing.T) {
	t.Parallel()

	router := New()
	method := &withParamsValidationMethod{}

	assert.Nil(t, router.RegisterMethod(method))

	res, err := router.Invoke(method.GetName(), validatableParams{Value: true})
	assert.Nil(t, err)
	assert.Equal(t, true, res)

	res, err = router.Invoke(method.GetName(), validatableParams{Value: false})
	assert.Nil(t, res)
	assert.Equal(t, "foo", err.GetData())
}

func TestRouter_InvokeWithParamsValidationAndJsonMarshaling(t *testing.T) {
	t.Parallel()

	router := New()
	method := &withParamsValidationMethod{}

	assert.Nil(t, router.RegisterMethod(method))

	res, err := router.Invoke(method.GetName(), map[string]bool{"my_value": true})
	assert.Nil(t, err)
	assert.Equal(t, true, res)

	res, err = router.Invoke(method.GetName(), map[string]bool{"my_value": false})
	assert.Nil(t, res)
	assert.Equal(t, "foo", err.GetData())
}
