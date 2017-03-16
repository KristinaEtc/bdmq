package transport

import "errors"

var (
	errEmptyLinkRepository = errors.New("No Link Descriptions in Node map")
	errQuitLinkRequested   = errors.New("Quit link requested")
)
