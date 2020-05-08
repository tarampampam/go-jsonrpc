package response

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	rpcErrors "github.com/tarampampam/go-jsonrpc/errors"
)

func TestJsonMarshalingResult(t *testing.T) {
	t.Parallel()

	res, err := json.Marshal(Response{Version: "1.2", Result: "foo", ID: "bar"})

	assert.Nil(t, err)
	assert.JSONEq(t, `{"jsonrpc":"1.2", "result":"foo", "id":"bar"}`, string(res))
}

func TestJsonMarshalingError(t *testing.T) {
	t.Parallel()

	res, err := json.Marshal(Response{Version: "1.2", Error: &rpcErrors.Error{Code: 1, Message: "foo"}, ID: "bar"})

	assert.Nil(t, err)
	assert.JSONEq(t, `{"jsonrpc":"1.2", "error":{"code":1, "message":"foo"}, "id":"bar"}`, string(res))
}

func TestNewResponsesAndAdd(t *testing.T) {
	t.Parallel()

	responses := NewResponses()
	responses.Add(Response{ID: 1})

	assert.Equal(t, 1, responses.Items[0].ID)
}

func TestNewResponsesAddConcurrent(t *testing.T) {
	t.Parallel()

	responses := NewResponses()
	wg := sync.WaitGroup{}
	addFunc := func() {
		responses.Add(Response{})

		wg.Done()
	}

	for i := 1; i <= 100; i++ {
		wg.Add(1)

		go addFunc()
	}

	wg.Wait()
	assert.Len(t, responses.Items, 100)
}
