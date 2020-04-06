package gpio

// ReceiverOption is the signature of funcs that are used to configure a
// *Receiver.
type ReceiverOption func(*Receiver)

// ReceiverProtocols configures the protocols the receiver should try to
// detect. If not overridden explicitly, DefaultProtocols is used.
func ReceiverProtocols(protocols []Protocol) ReceiverOption {
	return func(r *Receiver) {
		r.protocols = protocols
	}
}

// TransmitterOption is the signature of funcs that are used to configure a
// *Transmitter.
type TransmitterOption func(*Transmitter)

// TransmissionRetries configures how many times a code should be transmitted
// in a row. The higher the value, the more likely it is that an outlet
// actually received the code.
func TransmissionRetries(retries int) TransmitterOption {
	return func(t *Transmitter) {
		t.retries = retries
	}
}
