package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

const (
	actionOn     = "on"
	actionOff    = "off"
	actionToggle = "toggle"
)

// API type definition
type API struct {
	control *outlet.Control
}

// OutletRequest type definition
type OutletRequest struct {
	GroupId  int    `json:"groupId"`
	OutletId int    `json:"outletId"`
	Action   string `json:"action"`
}

// OutletGroupRequest type definition
type OutletGroupRequest struct {
	GroupId int    `json:"groupId"`
	Action  string `json:"action"`
}

// New create a new API
func New(control *outlet.Control) *API {
	return &API{control: control}
}

// StatusRequestHandler returns the outlet groups with the status of all
// contained outlets
func (a *API) StatusRequestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, a.control.OutletGroups())
}

// OutletGroupRequestHandler performs actions on the outlet group identified by
// groupId
func (a *API) OutletGroupRequestHandler(c *gin.Context) {
	var data OutletGroupRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outletGroup, err := a.control.OutletGroup(data.GroupId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.handleSwitchRequest(c, outletGroup, data.Action)
}

// OutletRequestHandler performs actions on the outlet identified by groupId and outletId
func (a *API) OutletRequestHandler(c *gin.Context) {
	var data OutletRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outlet, err := a.control.Outlet(data.GroupId, data.OutletId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.handleSwitchRequest(c, outlet, data.Action)
}

// handleSwitchRequest performs actions on a given switcher and returns the new
// state
func (a *API) handleSwitchRequest(c *gin.Context, s outlet.Switcher, action string) {
	var err error

	switch action {
	case actionOn:
		err = a.control.SwitchOn(s)
	case actionOff:
		err = a.control.SwitchOff(s)
	case actionToggle:
		err = a.control.ToggleState(s)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid action %q", action)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}
