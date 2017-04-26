package test

import (
	"runtime/debug"
	"strconv"
	"strings"

	"fmt"

	test "github.com/KristinaEtc/bdmq/test-commands"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var pwdCurr = "test-commands"
var log = slf.WithContext(pwdCurr)

func stringParam(str string) string {
	if strings.Compare("\"", str[len(str)-1:])+strings.Compare("\"", str[:1]) == 0 {
		strWithoutFirstQuot := str[1:]
		str = strWithoutFirstQuot[:len(strWithoutFirstQuot)-1]
	}
	return str
}

func boolParam(str string) (bool, error) {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return b, nil
}

func stopProcesser(n *transport.Node) test.ProcessorFunc {

	return func(signature string) error {

		log.Debugf("[command] stop=[%s]", signature)
		n.Stop()
		return nil
	}
}

func dumpProcesser(n *transport.Node) test.ProcessorFunc {

	return func(signature string) error {

		log.Debugf("[command] dump=[%s]", signature)
		debug.PrintStack()
		return nil
	}
}

func sendStringProcesser(n *transport.Node) test.ProcessorFunc {

	return func(signature string) error {

		log.Debugf("[command] sendMessageProcesser; signature=[%s]", signature)

		params := strings.SplitN(signature, ",", 2)

		if len(params) != 2 {
			log.Errorf("SendMessageProcesser:param=[%s]: must be 2 parametors", signature)
			return fmt.Errorf("SendMessageProcesser:param=[%s]: must be 2 parametors", signature)
		}

		var message, linkActiveID string
		message = stringParam(params[0])
		linkActiveID = stringParam(params[1])

		n.SendMessage(linkActiveID, message)
		log.Infof("Message [%s] was sent from [%s]", message, linkActiveID)
		return nil
	}
}

// Register call RegisterCommand for all commands declarated in this package
func Register(node *transport.Node, cmdRegistrar test.CommandRegistrar) {
	cmdRegistrar.RegisterCommand("sendString", sendStringProcesser(node))
	cmdRegistrar.RegisterCommand("dump", dumpProcesser(node))
	cmdRegistrar.RegisterCommand("stop", stopProcesser(node))
}
