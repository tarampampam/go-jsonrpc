package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_IsValid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name               string
		giveRequest        Request
		wantIsValid        bool
		wantErrorSubstring string
	}{
		{
			name:        "valid",
			giveRequest: Request{Version: "2.0", Method: "foo", Params: []int{1, 2}, ID: int(1)},
			wantIsValid: true,
		},
		{
			name:        "valid (with empty params)",
			giveRequest: Request{Version: "2.0", Method: "foo", Params: nil, ID: int(1)},
			wantIsValid: true,
		},
		{
			name:        "valid (with empty ID)",
			giveRequest: Request{Version: "2.0", Method: "foo", Params: nil, ID: nil},
			wantIsValid: true,
		},
		{
			name:        "valid (params as an object)",
			giveRequest: Request{Version: "2.0", Method: "foo", Params: Request{}, ID: nil},
			wantIsValid: true,
		},
		{
			name:        "valid (negative id value)",
			giveRequest: Request{Version: "2.0", Method: "foo", Params: nil, ID: -666},
			wantIsValid: true,
		},
		{
			name:               "wrong version",
			giveRequest:        Request{Version: "2.1", Method: "foo", Params: []int{1, 2}, ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong version",
		},
		{
			name:               "empty method",
			giveRequest:        Request{Version: "2.0", Method: "", Params: []int{1, 2}, ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "empty method",
		},
		{
			name:               "empty version",
			giveRequest:        Request{Version: "", Method: "foo", Params: []int{1, 2}, ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong version",
		},
		{
			name:               "wrong ID type",
			giveRequest:        Request{Version: "2.0", Method: "foo", Params: nil, ID: true},
			wantIsValid:        false,
			wantErrorSubstring: "wrong id type",
		},
		{
			name:               "wrong params (string)",
			giveRequest:        Request{Version: "2.0", Method: "foo", Params: "bar", ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong params type",
		},
		{
			name:               "wrong params (bool)",
			giveRequest:        Request{Version: "2.0", Method: "foo", Params: false, ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong params type",
		},
		{
			name:               "wrong params (uint32)",
			giveRequest:        Request{Version: "2.0", Method: "foo", Params: uint32(123), ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong params type",
		},
		{
			name:               "wrong params (float32)",
			giveRequest:        Request{Version: "2.0", Method: "foo", Params: float32(123), ID: int(1)},
			wantIsValid:        false,
			wantErrorSubstring: "wrong params type",
		},
		{
			name:        "wrong (empty) request",
			giveRequest: Request{},
			wantIsValid: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.giveRequest.Validate()

			if tt.wantIsValid {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)

				if tt.wantErrorSubstring != "" {
					assert.Contains(t, tt.wantErrorSubstring, result.Error())
				}
			}
		})
	}
}

func TestJsonMarshaling(t *testing.T) {
	t.Parallel()

	res, err := json.Marshal(Request{Version: "1.2", Method: "foo", Params: []int{1, 2}, ID: "bar"})

	assert.Nil(t, err)
	assert.JSONEq(t, `{"jsonrpc":"1.2", "method":"foo", "params":[1,2], "id":"bar"}`, string(res))
}
