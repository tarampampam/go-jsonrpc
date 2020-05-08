package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestsServing(t *testing.T) {
	t.Parallel()

	var cases = []struct {
		name     string
		giveBody string
		wantJSON string
	}{
		{
			name:     "regular method invoking",
			giveBody: `{"jsonrpc":"2.0", "method":"my.method", "params":{"number":1, "is_optional":"foo"}, "id":1}`,
			wantJSON: `{"jsonrpc":"2.0", "result":{"number":1, "is_optional":"foo", "string":"default"}, "id":1}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req, _  = http.NewRequest(http.MethodPost, "http://rpc", bytes.NewBuffer([]byte(tt.giveBody)))
				rr      = httptest.NewRecorder()
				router  = createRouter()
				handler = httpHandler(router)
			)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.JSONEq(t, tt.wantJSON, rr.Body.String())
		})
	}
}
