package transport

/*
import "github.com/ventu-io/slf"

var pwdCurr = "KristinaEtc/transport.go"
var log = slf.WithContext(pwdCurr)
*/

// Connection is a struct with data about active connection
type Connection struct {
	handler string
}

// ConnDescription stores all configurate information about connection
type ConnDescription struct {
	address  string
	nodeType string //client of server
}

// TCPNode is a struct with all information about node.
// It consist of Conn and ConnDesc structures with active connection and connection options respectively.
type TCPNode struct {
	Conn     Connection
	ConnDesc ConnDescription
}

//TCPNodeFactory realizes creater interface.
type TCPNodeFactory struct {
	*Factory
	nodes map[string]*TCPNode
}

func (t *TCPNodeFactory) CreateNode() TCPNode {
	node := TCPNode{}
	return node
}

type LinkDescription struct {
	Address string
	Name    string
}

func (t *TCPNode) Disconnect() {

}

func (t *TCPNode) Write() {

}

func (t *TCPNode) GetStatus() {

}
