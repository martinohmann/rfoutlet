package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/api"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

var (
	transmitter, _ = gpio.NewNullTransmitter()
	sm             = outlet.NewNullStateManager()
	config         = &outlet.Config{
		OutletGroups: []*outlet.OutletGroup{
			{
				Identifier: "foo",
				Outlets: []*outlet.Outlet{
					{
						Identifier: "bar",
						Protocol:   1,
					},
				},
			},
		},
	}
)

func createControl() *outlet.Control {
	return outlet.NewControl(config, sm, transmitter)
}

func TestStatusRequest(t *testing.T) {
	a := api.New(createControl())
	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	a.HandleStatusRequest(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `[{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":0}]}]`, rr.Body.String())
}

func TestValidateRequest(t *testing.T) {
	a := api.New(createControl())
	f := a.ValidateRequest(func(w http.ResponseWriter, r *http.Request, action string, groupId int) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	tests := []struct {
		method   string
		action   string
		code     int
		body     string
		postForm url.Values
	}{
		{
			method: "GET",
			code:   http.StatusMethodNotAllowed,
			body:   `{"error":"Method Not Allowed"}`,
		},
		{
			method: "POST",
			action: "foo",
			code:   http.StatusBadRequest,
			body:   `{"error":"foo is not a valid action"}`,
		},
		{
			method: "POST",
			action: "on",
			code:   http.StatusBadRequest,
			body:   `{"error":"group_id field missing"}`,
		},
		{
			method:   "POST",
			action:   "on",
			code:     http.StatusBadRequest,
			postForm: url.Values{"group_id": []string{"foo"}},
			body:     `{"error":"strconv.Atoi: parsing \"foo\": invalid syntax"}`,
		},
		{
			method:   "POST",
			action:   "off",
			code:     http.StatusOK,
			postForm: url.Values{"group_id": []string{"0"}},
			body:     `ok`,
		},
	}

	for _, tt := range tests {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest(tt.method, fmt.Sprintf("/api/outlet/%s", tt.action), nil)
		if err != nil {
			t.Fatal(err)
		}

		req.PostForm = tt.postForm

		f(rr, req)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}

func TestOutletRequest(t *testing.T) {
	a := api.New(createControl())

	tests := []struct {
		action   string
		groupId  int
		code     int
		body     string
		postForm url.Values
	}{
		{
			code: http.StatusBadRequest,
			body: `{"error":"outlet_id field missing"}`,
		},
		{
			code:     http.StatusBadRequest,
			postForm: url.Values{"outlet_id": []string{"2"}},
			groupId:  1,
			body:     `{"error":"invalid outlet group offset 1"}`,
		},
		{
			code:     http.StatusBadRequest,
			groupId:  0,
			postForm: url.Values{"outlet_id": []string{"2"}},
			body:     `{"error":"invalid outlet offset 2 in group 0"}`,
		},
		{
			code:     http.StatusOK,
			groupId:  0,
			action:   "on",
			postForm: url.Values{"outlet_id": []string{"0"}},
			body:     `{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}`,
		},
		{
			code:     http.StatusOK,
			groupId:  0,
			action:   "off",
			postForm: url.Values{"outlet_id": []string{"0"}},
			body:     `{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":2}`,
		},
		{
			code:     http.StatusOK,
			groupId:  0,
			action:   "toggle",
			postForm: url.Values{"outlet_id": []string{"0"}},
			body:     `{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}`,
		},
	}

	for _, tt := range tests {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("POST", fmt.Sprintf("/api/outlet/%s", tt.action), nil)
		if err != nil {
			t.Fatal(err)
		}

		req.PostForm = tt.postForm

		a.HandleOutletRequest(rr, req, tt.action, tt.groupId)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}

func TestOutletGroupRequest(t *testing.T) {
	a := api.New(createControl())

	tests := []struct {
		action   string
		groupId  int
		code     int
		body     string
		postForm url.Values
	}{
		{
			code:    http.StatusBadRequest,
			groupId: 1,
			body:    `{"error":"invalid outlet group offset 1"}`,
		},
		{
			code:    http.StatusOK,
			groupId: 0,
			action:  "on",
			body:    `{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":1}]}`,
		},
		{
			code:    http.StatusOK,
			groupId: 0,
			action:  "off",
			body:    `{"identifier":"foo","outlets":[{"identifier":"bar","pulse_length":0,"protocol":1,"code_on":0,"code_off":0,"state":2}]}`,
		},
	}

	for _, tt := range tests {
		rr := httptest.NewRecorder()

		req, err := http.NewRequest("POST", fmt.Sprintf("/api/outlet_group/%s", tt.action), nil)
		if err != nil {
			t.Fatal(err)
		}

		req.PostForm = tt.postForm

		a.HandleOutletGroupRequest(rr, req, tt.action, tt.groupId)

		assert.Equal(t, tt.code, rr.Code)
		assert.Equal(t, tt.body, rr.Body.String())
	}
}
