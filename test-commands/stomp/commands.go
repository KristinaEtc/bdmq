package stomp

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
	test "github.com/KristinaEtc/bdmq/test-commands"
	"github.com/ventu-io/slf"
)

var frameReceived int

var log = slf.WithContext("test-stomp")

func stringParam(str string) string {

	var strUnquoted string
	var err error

	if strings.Compare("\"", str[len(str)-1:])+strings.Compare("\"", str[:1]) == 0 {
		strUnquoted, err = strconv.Unquote(str)
		if err != nil {
			return str
		}
	}
	return strUnquoted
}

func boolParam(str string) (bool, error) {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return b, nil
}

func stopProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

	return func(signature string) error {
		log.Debugf("[command] stop=[%s]", signature)
		n.Stop()
		return nil
	}
}

func dumpProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

	return func(signature string) error {
		log.Debugf("[command] dump=[%s]", signature)
		debug.PrintStack()
		return nil
	}
}

func subscribeProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

	return func(signature string) error {

		log.Debugf("[command] subscribe=[%s]", signature)

		topic := stringParam(signature)

		_, err := n.Subscribe(topic)
		if err != nil {
			log.Errorf("Could not subscribe: %s", err.Error())
			return err
		}

		//go read(ch)
		log.Infof("Subscribed to [%s] topic", topic)

		return nil
	}
}

func unsubscribeProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

	return func(signature string) error {

		log.Debugf("[command] unsubscribe=[%s]", signature)
		//n.Unsubscribe(stringParam(signature))
		return nil
	}
}

func sendMessageProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

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

func sendFrameProcesser(n *stomp.NodeStomp) test.ProcessorFunc {

	return func(signature string) error {

		params := strings.Split(signature, ",")

		log.Debugf("[params] sendFrame=[%v]", params)

		if len(params) < 5 {
			return errors.New("Must be at least 5 parameters")
		}

		var topic, frameType string

		// parsing parameters from string signature

		topic = stringParam(params[0])
		frameType = stringParam(params[1])
		numOfFrames, err := strconv.Atoi(params[2])
		if err != nil {
			return errors.New("Wrong [num of frames] parameter")
		}

		// a creating of headers of frames without last kye headers
		// he will be added with his number

		var headers = make([]string, 0)
		for i := 3; i <= len(params)-2; i++ {
			h := stringParam(params[i])
			headers = append(headers, h)
		}

		// sending frames with last header key with his serial number

		var fr *frame.Frame
		h := stringParam(params[len(params)-1])
		for i := 1; i <= numOfFrames; i++ {
			headersFull := append(headers, h+strconv.Itoa(i))
			fr = frame.New(
				frameType,
				headersFull...,
			)
			n.SendFrame(topic, *fr)
		}
		//	log.Infof("[%d] frame(s) was sent from [%s]", numOfFrames, topic)

		return nil
	}
}

func read(ch chan frame.Frame, prescriptiveFrame frame.Frame, showInput bool) error {

	timeout := time.NewTicker(time.Second * 10)

	log.Debug("read")

	for {
		select {
		case fr := <-ch:

			frameReceived++
			if showInput {
				log.Infof("Received: %s", fr.Dump())
			}
			if frame.CompareFrames(prescriptiveFrame, fr) != true {
				return errors.New("Frames not equal")
			}
			log.Debug("Right frame")
			return nil

		case _ = <-timeout.C:
			return errors.New("timeout")
		}
	}
}

func receiveFrameProcesser(n *stomp.NodeStomp, showInput bool) test.ProcessorFunc {

	return func(signature string) error {

		var frameReceived int

		params := strings.Split(signature, ",")

		log.Debugf("[params] receiveFrame=[%v]", params)

		if len(params) < 5 {
			log.Error("Wrong signature: must be at least 5 parameters")
			return errors.New("Wrong signature: must be at least 5 parameters")
		}

		var topic, frameType string

		topic = stringParam(params[0])
		ch, err := n.Subscribe(topic)
		if err != nil {
			log.Errorf("topic=[%s]: %s", topic, err.Error())
			return err
		}

		frameType = stringParam(params[1])
		numOfFrames, err := strconv.Atoi(params[2])
		if err != nil {
			return errors.New("Wrong [number of frames] parameter.")
		}

		var headers = make([]string, 0)
		for i := 3; i <= len(params)-2; i++ {
			h := stringParam(params[i])
			headers = append(headers, h)
		}

		var fr *frame.Frame
		h := stringParam(params[len(params)-1])

		for i := 1; i <= numOfFrames; i++ {
			headersFull := append(headers, h+strconv.Itoa(i))
			fr = frame.New(
				frameType,
				headersFull...,
			)

			err := read(ch, *fr, showInput)
			if err != nil {
				log.Errorf("FramesReceived=[%d]", frameReceived)
				return err
			}
		}
		log.Infof("FramesReceived=[%d]", frameReceived)
		return nil
	}
}

// Register call RegisterCommand for all commands declarated in this package
func Register(node *stomp.NodeStomp, cmdRegistrar test.CommandRegistrar, showInput bool) {
	cmdRegistrar.RegisterCommand("subscribe", subscribeProcesser(node))
	cmdRegistrar.RegisterCommand("receiveFrame", receiveFrameProcesser(node, showInput))
	cmdRegistrar.RegisterCommand("sendFrame", sendFrameProcesser(node))
	cmdRegistrar.RegisterCommand("unsubscribe", unsubscribeProcesser(node))
	cmdRegistrar.RegisterCommand("dump", dumpProcesser(node))
	cmdRegistrar.RegisterCommand("stop", stopProcesser(node))
}
