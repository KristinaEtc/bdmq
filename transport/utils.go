package transport

type LinkDescFromJSON struct {
<<<<<<< Updated upstream
	LinkID         string
	Address        string
	Mode           string
	Handler        string
	BufSize        int
	FrameProcessor string
=======
	LinkID     string
	Address    string
	Mode       string
	Handler    string
	BufSize    int
	QuequeName string
>>>>>>> Stashed changes
}

/*
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
*/
