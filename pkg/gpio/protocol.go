package gpio

// HighLow defines the number of high pulses followed by a number of low pulses
// to send.
type HighLow struct {
	High, Low uint
}

// Protocol defines the HighLow sequences to send to emit ones (One) and zeros
// (Zero) and the sync sequence (Sync) which signals the end of a code
// transmission.
type Protocol struct {
	Sync, Zero, One HighLow
}

// DefaultProtocols defines known remote control protocols. These are exported
// to give users the ability to add more protocols if needed. However, it is
// advised to use the ReceiverProtocols ReceiverOption to configure a *Receiver
// with custom protocols.
var DefaultProtocols = []Protocol{
	{HighLow{1, 31}, HighLow{1, 3}, HighLow{3, 1}},
	{HighLow{1, 10}, HighLow{1, 2}, HighLow{2, 1}},
	{HighLow{30, 71}, HighLow{4, 11}, HighLow{9, 6}},
	{HighLow{1, 6}, HighLow{1, 3}, HighLow{3, 1}},
	{HighLow{6, 14}, HighLow{1, 2}, HighLow{2, 1}},
}
