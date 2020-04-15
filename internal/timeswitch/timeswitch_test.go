package timeswitch

import (
	"context"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/controller/commands"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestTimeSwitch(t *testing.T) {
	now := time.Now()
	plus1 := now.Add(time.Hour)

	tests := []struct {
		name             string
		outlets          []*outlet.Outlet
		expectedCommands []command.Command
	}{
		{
			name: "no outlets",
		},
		{
			name: "outlet with no schedule",
			outlets: []*outlet.Outlet{
				{Schedule: schedule.New()},
			},
		},
		{
			name: "outlet with disabled interval",
			outlets: []*outlet.Outlet{
				{
					Schedule: schedule.NewWithIntervals([]schedule.Interval{
						{
							Enabled:  false,
							Weekdays: []time.Weekday{now.Weekday()},
							From:     schedule.NewDayTime(now.Hour(), now.Minute()),
							To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
						},
					}),
				},
			},
		},
		{
			name: "disabled outlet should be enabled",
			outlets: []*outlet.Outlet{
				{
					State: outlet.StateOff,
					Schedule: schedule.NewWithIntervals([]schedule.Interval{
						{
							Enabled:  true,
							Weekdays: []time.Weekday{now.Weekday()},
							From:     schedule.NewDayTime(now.Hour(), now.Minute()),
							To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
						},
					}),
				},
			},
			expectedCommands: []command.Command{
				commands.StateCorrectionCommand{
					DesiredState: outlet.StateOn,
					Outlet: &outlet.Outlet{
						State: outlet.StateOff,
						Schedule: schedule.NewWithIntervals([]schedule.Interval{
							{
								Enabled:  true,
								Weekdays: []time.Weekday{now.Weekday()},
								From:     schedule.NewDayTime(now.Hour(), now.Minute()),
								To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
							},
						}),
					},
				},
			},
		},
		{
			name: "enabled outlet should not be enabled again",
			outlets: []*outlet.Outlet{
				{
					State: outlet.StateOn,
					Schedule: schedule.NewWithIntervals([]schedule.Interval{
						{
							Enabled:  true,
							Weekdays: []time.Weekday{now.Weekday()},
							From:     schedule.NewDayTime(now.Hour(), now.Minute()),
							To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
						},
					}),
				},
			},
		},
		{
			name: "enabled outlet should be switched off",
			outlets: []*outlet.Outlet{
				{
					State: outlet.StateOn,
					Schedule: schedule.NewWithIntervals([]schedule.Interval{
						{
							Enabled:  true,
							Weekdays: []time.Weekday{now.Weekday()},
							From:     schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
							To:       schedule.NewDayTime(now.Hour(), now.Minute()),
						},
					}),
				},
			},
			expectedCommands: []command.Command{
				commands.StateCorrectionCommand{
					DesiredState: outlet.StateOff,
					Outlet: &outlet.Outlet{
						State: outlet.StateOn,
						Schedule: schedule.NewWithIntervals([]schedule.Interval{
							{
								Enabled:  true,
								Weekdays: []time.Weekday{now.Weekday()},
								From:     schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
								To:       schedule.NewDayTime(now.Hour(), now.Minute()),
							},
						}),
					},
				},
			},
		},
		{
			name: "disabled outlet should not be switched off again",
			outlets: []*outlet.Outlet{
				{
					State: outlet.StateOff,
					Schedule: schedule.NewWithIntervals([]schedule.Interval{
						{
							Enabled:  true,
							Weekdays: []time.Weekday{now.Weekday()},
							From:     schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
							To:       schedule.NewDayTime(now.Hour(), now.Minute()),
						},
					}),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reg := outlet.NewRegistry()
			reg.RegisterOutlets(test.outlets...)

			queue := make(chan command.Command)
			defer close(queue)

			timeSwitch := New(reg, queue)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()

			var commands []command.Command
			doneCh := make(chan struct{})

			go func() {
				defer close(doneCh)
				for {
					select {
					case <-ctx.Done():
						return
					case cmd, ok := <-queue:
						if !ok {
							return
						}
						commands = append(commands, cmd)
					}
				}
			}()

			timeSwitch.check()

			<-doneCh

			assert.Equal(t, test.expectedCommands, commands)
		})
	}
}
