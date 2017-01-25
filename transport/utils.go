package transport

import "encoding/json"

// LinkDescFromJSON contains options for connecting to an address.
// It is used for parsing from JSON config (all elements are public).
type LinkDescFromJSON struct {
	LinkID  string
	Address string
	Mode    string
	Handler string
}

func (n *Node) parseFullLinkConf(data []byte) ([]LinkDescFromJSON, error) {

	var l = []LinkDescFromJSON{}
	err := json.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// Be aware: parseLinkConf returns a pointer to LinkDescFromJSON,
// parseFullLinkConf returns a list of structs LinkDescFromJSON.
func (n *Node) parseLinkConf(data []byte) (*LinkDescFromJSON, error) {

	var l = LinkDescFromJSON{}
	err := json.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

/*
func (l *LinkDesc) initHandler(hType string) error {
	h := transport.HandlerFunc{}
}
*/

/*
err := lDesc.initHandler(l.HandlerType)
		if err != nil {
			log.WithField("ID=", l.LinkID).Warnf("Error while creating handler: %s. Will be used no handlers!", err.Error())
		}
*/
