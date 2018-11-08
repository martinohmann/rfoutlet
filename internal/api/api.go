package api

import (
	"log"
	"net/http"
	"os"

	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	actionOn     = "on"
	actionOff    = "off"
	actionToggle = "toggle"
)

// APIHandlerFunc function type definition
type APIHandlerFunc func(http.ResponseWriter, *http.Request, string, int)

var (
	validActions = []string{actionOn, actionOff, actionToggle}
	logger       *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "api: ", log.LstdFlags|log.Lshortfile)
}

// API type definition
type API struct {
	control *outlet.Control
}

// New create a new API
func New(control *outlet.Control) *API {
	return &API{control: control}
}

// ValidateRequest ensures that a request is not malformed. If valid, the
// request is passed to the handler func, an error is returned otherwise
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

// HandleStatusRequest returns the outlet groups with the status of all
// contained outlets
func (a *API) HandleStatusRequest(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, a.control.OutletGroups(), http.StatusOK)
}

// HandleOutletGroupRequest performs actions on the outlet group identified by
// groupId
func (a *API) HandleOutletGroupRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletGroup, err := a.control.OutletGroup(groupId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.handleSwitchRequest(w, outletGroup, action)
}

// HandleOutletRequest performs actions on the outlet identified by groupId and outletId
func (a *API) HandleOutletRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletId, err := parseIntField(r, "outlet_id")
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	outlet, err := a.control.Outlet(groupId, outletId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.handleSwitchRequest(w, outlet, action)
}

// handleSwitchRequest performs actions on a given switcher and returns the new
// state
func (a *API) handleSwitchRequest(w http.ResponseWriter, s outlet.Switcher, action string) {
	var err error

	switch action {
	case actionOn:
		err = a.control.SwitchOn(s)
	case actionOff:
		err = a.control.SwitchOff(s)
	case actionToggle:
		err = a.control.ToggleState(s)
	}

	if err != nil {
		renderJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, s, http.StatusOK)
}
