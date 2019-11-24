// Package rfoutlet provides outlet control via cli and web interface for
// Raspberry PI 2/3.
//
// The transmitter and receiver logic has been ported from the great
// https://github.com/sui77/rc-switch C++ project to golang.
//
// rfoutlet comes with ready to use commands for transmitting and receiving
// remote control codes as well as a command for serving a web interface (see
// cmd/ directory). The pkg/ directory exposes the gpio package which contains
// the receiver and transmitter code.
package rfoutlet
