package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	ActionOn     = "on"
	ActionOff    = "off"
	ActionToggle = "toggle"
)

type APIHandlerFunc func(http.ResponseWriter, *http.Request, string, int)

var validActions = []string{ActionOn, ActionOff, ActionToggle}

type API struct {
	config *Config
}

func NewAPI(config *Config) *API {
	return &API{
		config: config,
	}
}

func (a *API) ValidateRequest(f APIHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			renderJSONError(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		action, err := parseAction(r)
		if err != nil {
			renderJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		groupId, err := parseIntField(r, "group_id")
		if err != nil {
			renderJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		f(w, r, action, groupId)
	})
}

func (a *API) HandleStatusRequest(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, a.config.OutletGroups, http.StatusOK)
}

func (a *API) HandleOutletGroupRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletGroup, err := a.config.OutletGroup(groupId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch action {
	case ActionOn:
		_ = outletGroup.SwitchOn()
	case ActionOff:
		_ = outletGroup.SwitchOff()
	case ActionToggle:
		_ = outletGroup.ToggleState()
	}

	log.Printf("%s, action: %s, group_id: %d\n", r.RequestURI, action, groupId)

	renderJSON(w, outletGroup, http.StatusOK)
}

func (a *API) HandleOutletRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletGroup, err := a.config.OutletGroup(groupId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	outletId, err := parseIntField(r, "outlet_id")
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	outlet, err := outletGroup.Outlet(outletId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch action {
	case ActionOn:
		_ = outlet.SwitchOn()
	case ActionOff:
		_ = outlet.SwitchOff()
	case ActionToggle:
		_ = outlet.ToggleState()
	}

	log.Printf("%s, action: %s, group_id: %d, outlet_id: %d\n", r.RequestURI, action, groupId, outletId)

	renderJSON(w, outlet, http.StatusOK)
}

func parseIntField(r *http.Request, fieldName string) (int, error) {
	rawValue := r.FormValue(fieldName)
	if rawValue == "" {
		return 0, fmt.Errorf("%s field missing", fieldName)
	}

	id, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func parseAction(r *http.Request) (string, error) {
	urlParts := strings.Split(r.URL.Path, "/")

	if len(urlParts) < 3 {
		return "", errors.New("invalid url path")
	}

	action := urlParts[3]
	if !isValidAction(action) {
		return "", fmt.Errorf("%s is not a valid action", action)
	}

	return action, nil
}

func isValidAction(action string) bool {
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}

func renderJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	responseBody, err := json.Marshal(payload)

	if err != nil {
		renderJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}

func renderJSONError(w http.ResponseWriter, msg string, statusCode int) {
	payload := make(map[string]string)
	payload["error"] = msg

	responseBody, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}
