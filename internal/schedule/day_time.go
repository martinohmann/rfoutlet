package schedule

// DayTime type definition
type DayTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// NewDayTime create a new day time
func NewDayTime(hour, minute int) DayTime {
	return DayTime{hour, minute}
}

// Equal returns true if t is equal to other
func (t DayTime) Equal(other DayTime) bool {
	return t.Hour == other.Hour && t.Minute == other.Minute
}

// Before returns true if t is before other
func (t DayTime) Before(other DayTime) bool {
	if t.Hour < other.Hour {
		return true
	}

	if t.Hour > other.Hour {
		return false
	}

	return t.Minute < other.Minute
}

// After returns true if t is after other
func (t DayTime) After(other DayTime) bool {
	return !t.Equal(other) && !t.Before(other)
}

// Between returns true if t is between start and end or equal to start
func (t DayTime) Between(start, end DayTime) bool {
	return t.Equal(start) || (t.After(start) && t.Before(end))
}
