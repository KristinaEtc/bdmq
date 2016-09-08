// Copyright © 2015 Clement 'cmc' Rey <cr.rey.clement@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tcprec

// ----------------------------------------------------------------------------

// Error is the error type of the GAS package.
//
// It implements the error interface.
type Error int

const (
	// ErrMaxRetries is returned when the called function failed after the
	// maximum number of allowed tries.
	ErrMaxRetries Error = 0x01
	// ErrEmptyNodeSlice is returned when user trying to init connections with
	//empty slice
	ErrEmptyNodeSlice Error = 0x02
	// ErrWrongNodeMode is an error which appears when i user set wrong node mode
	ErrWrongNodeMode Error = 0x03
)

// ----------------------------------------------------------------------------

// Error returns the error as a string.
func (e Error) Error() string {
	switch e {
	case ErrMaxRetries:
		return "ErrMaxRetries"
	case ErrEmptyNodeSlice:
		return "could not init connections: no datalink's slice."
	case ErrWrongNodeMode:
		return "Wrong mode in node; will be ignored."
	default:
		return "unknown error"
	}
}
