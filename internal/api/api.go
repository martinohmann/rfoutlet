package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
)

const (
	actionOn     = "on"
	actionOff    = "off"
	actionToggle = "toggle"
)

// API type definition
type API struct {
	ctx     *context.Context
	control *control.Control
}

// OutletRequest type definition
type OutletRequest struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// GroupRequest type definition
type GroupRequest struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// IntervalRequest type definition
type IntervalRequest struct {
	ID       string            `json:"id"`
	Interval schedule.Interval `json:"interval"`
}

// New create a new API
func New(ctx *context.Context, control *control.Control) *API {
	return &API{ctx: ctx, control: control}
}

// StatusRequestHandler returns the outlet groups with the status of all
// contained outlets
func (a *API) StatusRequestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, a.ctx.Groups)
}

// GroupRequestHandler performs actions on an outlet group
func (a *API) GroupRequestHandler(c *gin.Context) {
	var data GroupRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	og, err := a.ctx.GetGroup(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, o := range og.Outlets {
		if err := a.performAction(o, data.Action); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, og)
}

// OutletRequestHandler performs actions on an outlet
func (a *API) OutletRequestHandler(c *gin.Context) {
	var data OutletRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := a.ctx.GetOutlet(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.performAction(o, data.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

// IntervalRequestHandler handles interval create/update/delete requests
func (a *API) IntervalRequestHandler(c *gin.Context) {
	var data IntervalRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := a.ctx.GetOutlet(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch c.Request.Method {
	case http.MethodPut:
		err = a.control.AddInterval(o, data.Interval)
	case http.MethodPost:
		err = a.control.UpdateInterval(o, data.Interval)
	case http.MethodDelete:
		err = a.control.DeleteInterval(o, data.Interval)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
}

func (a *API) performAction(o *context.Outlet, action string) error {
	switch action {
	case actionOn:
		return a.control.SwitchState(o, state.SwitchStateOn)
	case actionOff:
		return a.control.SwitchState(o, state.SwitchStateOff)
	case actionToggle:
		return a.control.Toggle(o)
	}

	return fmt.Errorf("invalid action %q", action)
}
