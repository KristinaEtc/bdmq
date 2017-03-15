package transport

const (
	ServerMode int = iota
	ClientMode
)

type CommandID int

const (
	registerControl CommandID = iota
	unregisterControl
	registerActive
	unregisterActive
	stopNode
	sendMessageNode
	sendMessageByIDNode
)

type NodeCommand struct {
	Cmd CommandID
}

func (nC *NodeCommand) GetCommandID() CommandID {
	return nC.Cmd
}

type NodeCommandControlLink struct {
	NodeCommand
	ctrl *LinkControl
}

type NodeCommandActiveLink struct {
	NodeCommand
	active *LinkActive
}

type NodeCommandSendMessage struct {
	NodeCommand
	quequeName string
	msg        string
}

type cmdActiveLink struct {
	cmd commandActiveLink
	msg []byte
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
