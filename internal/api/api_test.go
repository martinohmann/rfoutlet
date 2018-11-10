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
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func createControl() *outlet.Control {
	t, _ := gpio.NewNullTransmitter()
	sm := outlet.NewNullStateManager()
	c := &outlet.Config{
		OutletGroups: []*outlet.OutletGroup{
			{
				Identifier: "foo",
				Outlets: []*outlet.Outlet{
					{
						Identifier: "bar",
						Protocol:   1,
					},
					{
						Identifier: "baz",
						Protocol:   1,
						State:      outlet.StateOn,
					},
				},
			},
		},
	}
	return outlet.NewControl(c, sm, t)
}

func TestInvalidJson(t *testing.T) {
	a := api.New(createControl())

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
			body: `{"error":"invalid action \"\""}`,
		},
		{
			code: http.StatusBadRequest,
			data: api.OutletRequest{
				GroupId:  1,
				OutletId: 2,
			},
			body: `{"error":"invalid outlet group offset 1"}`,
		},
		{
			code: http.StatusBadRequest,
			data: api.OutletRequest{
				GroupId:  0,
				OutletId: 2,
			},
			body: `{"error":"invalid outlet offset 2 in group 0"}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				GroupId:  0,
				OutletId: 0,
				Action:   "on",
			},
			body: `{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				GroupId:  0,
				OutletId: 0,
				Action:   "toggle",
			},
			body: `{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}`,
		},
	}

	for _, tt := range tests {
		a := api.New(createControl())

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
				GroupId: 1,
			},
			body: `{"error":"invalid outlet group offset 1"}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletGroupRequest{
				GroupId: 0,
				Action:  "on",
			},
			body: `{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1},{"identifier":"baz","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}]}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletGroupRequest{
				GroupId: 0,
				Action:  "off",
			},
			body: `{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":2},{"identifier":"baz","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":2}]}`,
		},
		{
			code: http.StatusBadRequest,
			data: api.OutletGroupRequest{
				GroupId: 0,
				Action:  "foo",
			},
			body: `{"error":"invalid action \"foo\""}`,
		},
	}

	for _, tt := range tests {
		a := api.New(createControl())

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
	a := api.New(createControl())

	r := gin.New()
	r.POST("/api/status", a.StatusRequestHandler)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/status", nil)

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `[{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":0},{"identifier":"baz","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}]}]`, rr.Body.String())
}
