package transport

import "errors"

var (
	ErrEmptyLinkRepository = errors.New("No Link Descriptions in Node map.")
	ErrQuitLinkRequested   = errors.New("Quit link requested.")
)
