package transport

import (
	"bufio"
	"net"

	"github.com/ventu-io/slf"
)

// LinkActive initialize LinkStater interface when
// connecting is ok
type ActiveLink struct {
	conn         net.Conn
	linkDesc     *LinkDesc
	LinkActiveID string
	handler      Handler
	//parser       *parser.Parser
	commandCh   chan cmdActiveLink
	linkControl *LinkControl
	log         slf.Logger
}

func (lA *ActiveLink) Id() string {
	return lA.LinkActiveID
}

func (lA *