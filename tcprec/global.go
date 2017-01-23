package tcp

import (
	"encoding/json"

	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("tcprec.go")

// Global struct stores all created entities. It is a root of all
// work with this library.
// Global struct consist of Node list - all about transport layer,
// and Handler list - methods, which allow to interact with application layer
// library.
//
//Global realizes creator interface from transport library.
type Global struct {
	F        *transport.Factory
	Nodes    []Node
	Handlers *transport.Handlers
}

// InitHandlers creates Handlers, adds it to maps of Handlers
// and returns it.
func (g *Global) InitHandlers(h *transport.Handlers) {
	g.Handlers = h
}

// InitNodes parses config with nodes' options and creates nodes with such
// options.
//
// InitNodes is a method of creater interface from transport library.
func (g *Global) InitNodes(confJSON []byte) ([]Node, error) {

	log.Debug("(g *Global) InitNodes")

	err := parseConf(confJSON, &g.Nodes)
	if err != nil {
		log.Debugf("Initialize node from code config - err: %d", err.Error())
		return nil, err
	}

	return g.Nodes, nil
}

//utils
//
func parseConf(data []byte, nodes *[]Node) error {

	//pointers?
	err := json.Unmarshal(data, &nodes)
	if err != nil {
		return err
	}

	return nil
}
