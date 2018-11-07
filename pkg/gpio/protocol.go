package gpio

// HighLow type definition
type HighLow struct {
	High, Low uint
}

// Protocol type definition
type Protocol struct {
	PulseLength     uint
	Sync, Zero, One HighLow
}

// Protocols defines known remote control protocols. These are exported to give
// users the ability to add more protocols if needed.
var Protocols = []Protocol{
	{350, HighLow{1, 31}, HighLow{1, 3}, HighLow{3, 1}},
	{650, HighLow{1, 10}, HighLow{1, 2}, HighLow{2, 1}},
	{100, HighLow{30, 71}, HighLow{4, 11}, HighLow{9, 6}},
	{380, HighLow{1, 6}, HighLow{1, 3}, HighLow{3, 1}},
	{500, HighLow{6, 14}, HighLow{1, 2}, HighLow{2, 1}},
}
