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
	default:
		return "unknown"
	}
}

type NodeCommand struct {
	command commandNode
	ctrl    LinkControl
	active  *LinkActive
}

const (
	registerControl commandNode = iota
	unregisterControl
	registerActive
	unregisterActive
	stopNode
)

type commandActiveLink int

const (
	quitLinkActive commandActiveLink = iota
	errorReading
)

func (c commandActiveLink) String() string {
	switch c {
	case 0:
		return "quit"
	case 1:
		return "errorReading"
	default:
		return "unknown"
	}
}
