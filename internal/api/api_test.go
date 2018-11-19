package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/schedule"
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
	s.Schedules["qux"] = schedule.Schedule{
		{ID: "interval-1"},
	}

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
			handler:      a.GroupRequestHandler,
			expectedBody: `{"error":"EOF"}`,
		},
		{
			handler:      a.GroupRequestHandler,
			givenBody:    `{foo`,
			expectedBody: `{"error":"invalid character 'f' looking for beginning of object key string"}`,
		},
		{
			handler:      a.OutletRequestHandler,
			givenBody:    `[{}]`,
			expectedBody: `{"error":"json: cannot unmarshal array into Go value of type api.OutletRequest"}`,
		},
		{
			handler:      a.GroupRequestHandler,
			givenBody:    `[]`,
			expectedBody: `{"error":"json: cannot unmarshal array into Go value of type api.GroupRequest"}`,
		},
		{
			handler:      a.IntervalRequestHandler,
			givenBody:    `[]`,
			expectedBody: `{"error":"json: cannot unmarshal array into Go value of type api.IntervalRequest"}`,
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
		code int
		body string
		data api.OutletRequest
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
			body: `{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":false,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":0}`,
		},
		{
			code: http.StatusOK,
			data: api.OutletRequest{
				ID:     "bar",
				Action: "toggle",
			},
			body: `{"id":"bar","name":"Bar","schedule":null,"state":1}`,
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
		data api.GroupRequest
	}{
		{
			code: http.StatusBadRequest,
			data: api.GroupRequest{
				ID: "nonexistent",
			},
			body: `{"error":"group with identifier \"nonexistent\" does not exist"}`,
		},
		{
			code: http.StatusOK,
			data: api.GroupRequest{
				ID:     "foo",
				Action: "on",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":1},{"id":"baz","name":"Baz","schedule":null,"state":1},{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":false,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":1}]}`,
		},
		{
			code: http.StatusOK,
			data: api.GroupRequest{
				ID:     "foo",
				Action: "off",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":0},{"id":"baz","name":"Baz","schedule":null,"state":0},{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":false,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":0}]}`,
		},
		{
			code: http.StatusOK,
			data: api.GroupRequest{
				ID:     "foo",
				Action: "toggle",
			},
			body: `{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":1},{"id":"baz","name":"Baz","schedule":null,"state":1},{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":false,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":0}]}`,
		},
		{
			code: http.StatusInternalServerError,
			data: api.GroupRequest{
				ID:     "foo",
				Action: "foo",
			},
			body: `{"error":"invalid action \"foo\""}`,
		},
	}

	for _, tt := range tests {
		a := createAPI()

		r := gin.New()
		r.POST("/api/outlet_group", a.GroupRequestHandler)

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
	assert.Equal(t, `[{"id":"foo","name":"Foo","outlets":[{"id":"bar","name":"Bar","schedule":null,"state":0},{"id":"baz","name":"Baz","schedule":null,"state":0},{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":false,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":1}]}]`, rr.Body.String())
}

func TestIntervalRequest(t *testing.T) {
	tests := []struct {
		code   int
		method string
		body   string
		data   api.IntervalRequest
	}{
		{
			code:   http.StatusBadRequest,
			method: http.MethodPost,
			data: api.IntervalRequest{
				ID: "nonexistent",
			},
			body: `{"error":"outlet with identifier \"nonexistent\" does not exist"}`,
		},
		{
			code:   http.StatusInternalServerError,
			method: http.MethodPost,
			data: api.IntervalRequest{
				ID: "bar",
				Interval: schedule.Interval{
					ID: "foo",
				},
			},
			body: `{"error":"interval with identifier \"foo\" does not exist"}`,
		},
		{
			code:   http.StatusOK,
			method: http.MethodPut,
			data: api.IntervalRequest{
				ID: "bar",
				Interval: schedule.Interval{
					ID:       "interval-2",
					Enabled:  true,
					Weekdays: []time.Weekday{time.Monday, time.Tuesday},
				},
			},
			body: `{"id":"bar","name":"Bar","schedule":[{"id":"interval-2","enabled":true,"weekdays":[1,2],"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":0}`,
		},
		{
			code:   http.StatusOK,
			method: http.MethodPost,
			data: api.IntervalRequest{
				ID: "qux",
				Interval: schedule.Interval{
					ID:      "interval-1",
					Enabled: true,
				},
			},
			body: `{"id":"qux","name":"Qux","schedule":[{"id":"interval-1","enabled":true,"weekdays":null,"from":{"hour":0,"minute":0},"to":{"hour":0,"minute":0}}],"state":1}`,
		},
		{
			code:   http.StatusOK,
			method: http.MethodDelete,
			data: api.IntervalRequest{
				ID: "qux",
				Interval: schedule.Interval{
					ID: "interval-1",
				},
			},
			body: `{"id":"qux","name":"Qux","schedule":[],"state":1}`,
		},
	}

	for _, tt := range tests {
		a := createAPI()

		r := gin.New()
		r.PUT("/api/outlet/schedule", a.IntervalRequestHandler)
		r.POST("/api/outlet/schedule", a.IntervalRequestHandler)
		r.DELETE("/api/outlet/schedule", a.IntervalRequestHandler)

		rr := httptest.NewRecorder()
		json, _ := json.Marshal(tt.data)
		req, _ := http.NewRequest(tt.method, "/api/outlet/schedule", bytes.NewBuffer(json))

		r.ServeHTTP(rr, req)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}
