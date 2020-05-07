package errors

import (
	"encoding/json"
	"testing"

	"github.com/tarampampam/go-jsonrpc"

	"github.com/stretchr/testify/assert"
)

func TestImplementsInterface(t *testing.T) {
	t.Parallel()

	assert.Implements(t, (*jsonrpc.Error)(nil), new(Error))
}

func TestErrorCodeConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Code(-32700), Parse)
	assert.Equal(t, Code(-32600), InvalidRequest)
	assert.Equal(t, Code(-32601), MethodNotFound)
	assert.Equal(t, Code(-32602), InvalidParams)
	assert.Equal(t, Code(-32603), Internal)
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		giveCode    Code
		giveMessage string
		giveData    interface{}
		want        string
	}{
		{
			name:        "basic",
			giveCode:    123,
			giveMessage: "foo message",
			giveData:    []string{"foo", "bar"},
			want:        "code: 123, message: foo message, data: [foo bar]",
		},
		{
			name:        "without data",
			giveCode:    321,
			giveMessage: "bar message",
			giveData:    nil,
			want:        "code: 321, message: bar message",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := &Error{Code: tt.giveCode, Message: tt.giveMessage, Data: tt.giveData}

			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestJsonMarshaling(t *testing.T) {
	t.Parallel()

	res, err := json.Marshal(Error{Code: 1, Message: "foo", Data: []int{1, 2}})

	assert.Nil(t, err)
	assert.JSONEq(t, `{"code":1,"message":"foo","data":[1,2]}`, string(res))
}

func TestCode_String(t *testing.T) {
	t.Parallel()

	cases := []struct {
		giveCode   Code
		wantString string
	}{
		{giveCode: Parse, wantString: "Parse error"},
		{giveCode: InvalidRequest, wantString: "Invalid Request"},
		{giveCode: MethodNotFound, wantString: "Method not found"},
		{giveCode: InvalidParams, wantString: "Invalid params"},
		{giveCode: Internal, wantString: "Internal error"},
		{giveCode: Code(0), wantString: "Unrecognized error code"},
		{giveCode: Code(666), wantString: "Unrecognized error code"},
	}

	for _, tt := range cases {
		t.Run(tt.wantString, func(t *testing.T) {
			assert.Equal(t, tt.wantString, tt.giveCode.String())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	err := New(Parse)

	assert.Equal(t, Parse, err.Code)
	assert.Equal(t, Parse.String(), err.Message)
}

func TestErrorGetters(t *testing.T) {
	t.Parallel()

	err := Error{
		Code:    Code(123),
		Message: "foo",
		Data:    "bar",
	}

	assert.Equal(t, 123, err.GetCode())
	assert.Equal(t, "foo", err.GetMessage())
	assert.Equal(t, "bar", err.GetData())
}
