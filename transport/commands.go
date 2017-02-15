package transport

type commandNode int

func (c commandNode) String() string {
	switch c {
	case 0:
		return "registerControl"
	case 1:
		return "unregisterControl"
	case 2:
		return "registerActive"
	case 3:
		return "unregisterActive"
	case 4:
		return "stop"
	case 5:
		return "sendMessage"
	default:
		return "unknown"
	}
}

type NodeCommand struct {
	cmd    commandNode
	ctrl   *LinkControl
	active *LinkActive
	msg    string
}

const (
	registerControl commandNode = iota
	unregisterControl
	registerActive
	unregisterActive
	stopNode
	sendMessageNode
)

type cmdActiveLink struct {
	cmd commandActiveLink
	msg string
	err string
}

type commandActiveLink int

const (
	quitLinkActive commandActiveLink = iota
	errorReadingActive
	sendMessageActive
)

func (c commandActiveLink) String() string {
	switch c {
	case 0:
		return "quitLinkActive"
	case 1:
		return "errorReadingActive"
	case 2:
		return "sendMessage"
	default:
		return "unknown"
	}
}

type cmdContrlLink struct {
	cmd commandContrlLink
	err string
}

type commandContrlLink int

const (
	quitControlLink commandContrlLink = iota
	errorControlLinkAccept
	errorControlLinkRead

//	errorFromActiveLink
)

func (c commandContrlLink) String() string {
	switch c {
	case 0:
		return "quitControlLink"
	case 1:
		return "errorControlLinkAccept"
	case 2:
		return "errorControlLinkRead"
	default:
		return "unknown"
	}
}
