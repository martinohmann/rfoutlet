package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	uuid "github.com/satori/go.uuid"
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

// OutletGroupRequest type definition
type OutletGroupRequest struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

// OutletScheduleIntervalAddRequest type definition
type OutletScheduleIntervalAddRequest struct {
	ID       string            `json:"id"`
	Interval schedule.Interval `json:"interval"`
}

// OutletScheduleIntervalUpdateRequest type definition
type OutletScheduleIntervalUpdateRequest struct {
	ID       string            `json:"id"`
	Interval schedule.Interval `json:"interval"`
}

// OutletScheduleIntervalDeleteRequest type definition
type OutletScheduleIntervalDeleteRequest struct {
	ID         string `json:"id"`
	IntervalID string `json:"intervalId"`
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

// OutletGroupRequestHandler performs actions on an outlet group
func (a *API) OutletGroupRequestHandler(c *gin.Context) {
	var data OutletGroupRequest

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

// OutletScheduleIntervalDeleteRequestHandler adds a schedule interval for outlet
func (a *API) OutletScheduleIntervalDeleteRequestHandler(c *gin.Context) {
	var data OutletScheduleIntervalDeleteRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := a.ctx.GetOutlet(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, interval := range o.Schedule {
		if interval.ID == data.IntervalID {
			o.Schedule = append(o.Schedule[:i], o.Schedule[i+1:]...)
			break
		}
	}

	a.ctx.State.Schedules[o.ID] = o.Schedule

	a.control.SaveState()

	c.JSON(http.StatusOK, o)
}

// OutletScheduleIntervalAddRequestHandler adds a schedule interval for outlet
func (a *API) OutletScheduleIntervalAddRequestHandler(c *gin.Context) {
	var data OutletScheduleIntervalAddRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := a.ctx.GetOutlet(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Interval.ID == "" {
		data.Interval.ID = uuid.NewV4().String()
	}

	o.Schedule = append(o.Schedule, data.Interval)
	a.ctx.State.Schedules[o.ID] = o.Schedule

	a.control.SaveState()

	c.JSON(http.StatusOK, o)
}

// OutletScheduleIntervalUpdateRequestHandler adds a schedule interval for outlet
func (a *API) OutletScheduleIntervalUpdateRequestHandler(c *gin.Context) {
	var data OutletScheduleIntervalUpdateRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o, err := a.ctx.GetOutlet(data.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, interval := range o.Schedule {
		if interval.ID == data.Interval.ID {
			o.Schedule[i] = data.Interval

			a.ctx.State.Schedules[o.ID] = o.Schedule

			a.control.SaveState()

			c.JSON(http.StatusOK, o)
			return
		}
	}

	c.JSON(http.StatusBadRequest,
		gin.H{"error": fmt.Errorf("interval with identifier %q does not exist", data.Interval.ID)})
	return
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
