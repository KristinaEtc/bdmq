package parser

import (
	"net"
)

// Parser parses
type Reader struct {
}

func NewReader(conn net.Conn) *Reader {
	return nil
}

func (r *Reader) ParseData(numOfBytes int) ([]byte, error) {
	return nil, nil
}
