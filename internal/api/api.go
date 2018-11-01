package api

import (
	"net/http"

	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	ActionOn     = "on"
	ActionOff    = "off"
	ActionToggle = "toggle"
)

type APIHandlerFunc func(http.ResponseWriter, *http.Request, string, int)

var validActions = []string{ActionOn, ActionOff, ActionToggle}

type API struct {
	config *outlet.Config
}

func New(config *outlet.Config) *API {
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

	handleSwitchRequest(w, outletGroup, action)
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

	handleSwitchRequest(w, outlet, action)
}

func handleSwitchRequest(w http.ResponseWriter, s outlet.Switcher, action string) {
	var err error

	switch action {
	case ActionOn:
		err = s.SwitchOn()
	case ActionOff:
		err = s.SwitchOff()
	case ActionToggle:
		err = s.ToggleState()
	}

	if err != nil {
		renderJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, s, http.StatusOK)
}

func isValidAction(action string) bool {
	for _, validAction := range validActions {
		if action == validAction {
			return true
		}
	}
	return false
}
