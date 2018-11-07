package api

import (
	"log"
	"net/http"
	"os"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
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
	config      *outlet.Config
	transmitter gpio.CodeTransmitter
}

// New create a new API
func New(config *outlet.Config, transmitter gpio.CodeTransmitter) *API {
	return &API{
		config:      config,
		transmitter: transmitter,
	}
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
	renderJSON(w, a.config.OutletGroups, http.StatusOK)
}

// HandleOutletGroupRequest performs actions on the outlet group identified by
// groupId
func (a *API) HandleOutletGroupRequest(w http.ResponseWriter, r *http.Request, action string, groupId int) {
	outletGroup, err := a.config.OutletGroup(groupId)
	if err != nil {
		renderJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	handleSwitchRequest(w, outletGroup, a.transmitter, action)
}

// HandleOutletRequest performs actions on the outlet identified by groupId and outletId
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

	handleSwitchRequest(w, outlet, a.transmitter, action)
}

// handleSwitchRequest performs actions on a given switcher and returns the new
// state
func handleSwitchRequest(w http.ResponseWriter, s outlet.Switcher, t gpio.CodeTransmitter, action string) {
	var err error

	switch action {
	case actionOn:
		err = s.SwitchOn(t)
	case actionOff:
		err = s.SwitchOff(t)
	case actionToggle:
		err = s.ToggleState(t)
	}

	if err != nil {
		renderJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, s, http.StatusOK)
}
