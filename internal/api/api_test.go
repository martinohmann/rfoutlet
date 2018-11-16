package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func createContext() *context.Context {
	c := &config.Config{
		GroupOrder: []string{"foo"},
		Groups: map[string]*config.Group{
			"foo": {
				Name:    "Foo",
				Outlets: []string{"bar", "baz", "qux"},
			},
		},
		Outlets: map[string]*config.Outlet{
			"bar": {Name: "Bar", Protocol: 1},
			"baz": {Name: "Baz", Protocol: 1},
			"qux": {Name: "Qux", Protocol: 1},
		},
	}

	s := state.New()
	s.SwitchStates["qux"] = state.SwitchStateOn

	ctx, _ := context.New(c, s)

	return ctx
}

func createAPI() *api.API {
	ctx := createContext()
	t, _ := gpio.NewNullTransmitter()

	return api.New(ctx, control.New(ctx, t))
}

func TestInvalidJson(t *testing.T) {
	a := createAPI()

	tests := []struct {
		handler      gin.HandlerFunc
		givenBody    string
		expectedBody string
	}{
		{
			handler:      a.OutletRequestHandler,
			expectedBody: `{"error":"EOF"}`,
		},
		{
			handler:      a.OutletGroupRequestHandler,
			expectedBody: `{"error":"EOF"}`,
		},
		{
			handler:      a.OutletGroupRequestHandler,
			givenBody:    `{foo`,
			expectedBody: `{"error":"invalid character 'f' looking for beginning of object key string"}`,
		},
		{
			handler:      a.OutletRequestHandler,
			givenBody:    `[{}]`,
			expectedBody: `{"error":"json: cannot unmarshal array into Go value of type api.OutletRequest"}`,
		},
		{
			handler:      a.OutletGroupRequestHandler,
			givenBody:    `[]`,
			expectedBody: `{"error":"json: cannot unmarshal array into Go value of type api.OutletGroupRequest"}`,
		},
	}

	for _, tt := range tests {
		r := gin.New()
		r.POST("/", tt.handler)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(tt.givenBody)))

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, tt.expectedBody, rr.Body.String())
	}
}

func TestOutletRequest(t *testing.T) {
	tests := []struct {
		action   string
		groupId  int
		code     int
		body     string
		postForm url.Values
		data     api.OutletRequest
	}{
		{
			code: http.StatusBadRequest,
			data: api.OutletRequest{
				ID: "nonexistent",
			},
			body: `{"error":"outlet with identifier \"nonexistent\" does not exist"}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				ID:     "bar",
				Action: "on",
			},
			body: `{"id":"bar","name":"Bar","schedule":null,"state":1}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				ID:     "qux",
				Action: "off",
			},
			body: `{"id":"qux","name":"Qux","schedule":null,"state":0}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				ID:     "qux",
				Action: "toggle",
			},
			body: `{"id":"qux","name":"Qux","schedule":null,"state":0}`,
		},
		{
			code: http.StatusInternalServerError,
			data: api.OutletRequest{
				ID:     "bar",
				Action: "foo",
			},
			body: `{"error":"invalid action \"foo\""}`,
		},
	}

	for _, tt := range tests {
		a := createAPI()

		r := gin.New()
		r.POST("/api/outlet", a.OutletRequestHandler)

		rr := httptest.NewRecorder()
		json, _ := json.Marshal(tt.data)
		req, _ := http.NewRequest("POST", "/api/outlet", bytes.NewBuffer(json))

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}

func TestOutletGroupRequest(t *testing.T) {
	tests := []struct {
		code int
		body string
		data api.OutletGroupRequest
	}{
		{
			code: http.StatusBadRequest,
			data: api.OutletGroupRequest{
				ID: "nonexistent",
			},
			body: `{"error":"group with identifier \"nonexistent\" does not exist"}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletGroupRequest{
				ID:     "foo",
				Action: "on",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":1},{"id":"baz","name":"Baz","schedule":null,"state":1},{"id":"qux","name":"Qux","schedule":null,"state":1}]}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletGroupRequest{
				ID:     "foo",
				Action: "off",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":0},{"id":"baz","name":"Baz","schedule":null,"state":0},{"id":"qux","name":"Qux","schedule":null,"state":0}]}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletGroupRequest{
				ID:     "foo",
				Action: "toggle",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":1},{"id":"baz","name":"Baz","schedule":null,"state":1},{"id":"qux","name":"Qux","schedule":null,"state":0}]}`,
		},
		{
			code: http.StatusInternalServerError,
			data: api.OutletGroupRequest{
				ID:     "foo",
				Action: "foo",
			},
			body: `{"error":"invalid action \"foo\""}`,
		},
	}

	for _, tt := range tests {
		a := createAPI()

		r := gin.New()
		r.POST("/api/outlet_group", a.OutletGroupRequestHandler)

		rr := httptest.NewRecorder()
		json, _ := json.Marshal(tt.data)
		req, _ := http.NewRequest("POST", "/api/outlet_group", bytes.NewBuffer(json))

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}

func TestStatusRequestHandler(t *testing.T) {
	a := createAPI()

	r := gin.New()
	r.POST("/api/status", a.StatusRequestHandler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/status", nil)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `[{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":0},{"id":"baz","name":"Baz","schedule":null,"state":0},{"id":"qux","name":"Qux","schedule":null,"state":1}]}]`, rr.Body.String())
}
