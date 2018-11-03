package gpio

type highLow struct {
	high, low int
}

type protocol struct {
	pulseLength     int
	sync, zero, one highLow
}

var (
	protocols = []protocol{
		protocol{350, highLow{1, 31}, highLow{1, 3}, highLow{3, 1}},
		protocol{650, highLow{1, 10}, highLow{1, 2}, highLow{2, 1}},
		protocol{100, highLow{30, 71}, highLow{4, 11}, highLow{9, 6}},
		protocol{380, highLow{1, 6}, highLow{1, 3}, highLow{3, 1}},
		protocol{500, highLow{6, 14}, highLow{1, 2}, highLow{2, 1}},
	}
)
