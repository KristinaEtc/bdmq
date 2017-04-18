package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"fmt"

	"reflect"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
)

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

func dumpProcesser(signature string, n *stomp.NodeStomp) error {
	log.Debugf("[command] dump=[%s]", signature)
	return nil
}

func sleepProcesser(signature string, n *stomp.NodeStomp) error {
	log.Infof("[command] sleep=[%s]", signature)

	sec, err := strconv.Atoi(signature)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(sec))
	return nil
}

func read(ch chan frame.Frame, prescriptiveFrame *frame.Frame) error {

	timeout := time.NewTicker(time.Second * 7)

	log.Debug("read")

	for {
		select {
		case fr := <-ch:

			frameReceived++
			if globalOpt.ShowFrames {
				log.Infof("Received: %s", fr.Dump())
			}
			if reflect.DeepEqual(fr, prescriptiveFrame) != true {
				log.Debugf("Frames not equal: \n%+v\n%+v\n", *prescriptiveFrame, fr)
				log.Debugf("Frames not equal: \n%+v\n%+v\n", *prescriptiveFrame.Header, fr.Header)
				return errors.New("Frames not equal")
			}
			log.Debug("Right frame")
			return nil

		case _ = <-timeout.C:
			return errors.New("timeout")
		}
	}
}

func subscribeProcesser(signature string, n *stomp.NodeStomp) error {
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

func unsubscribeProcesser(signature string, n *stomp.NodeStomp) error {
	log.Debugf("[command] unsubscribe=[%s]", signature)
	//n.Unsubscribe(stringParam(signature))
	return nil
}

func sendMessageProcesser(signature string, n *stomp.NodeStomp) error {

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

func sendFrameProcesser(signature string, n *stomp.NodeStomp) error {

	params := strings.Split(signature, ",")

	log.Debugf("[params] sendFrame=[%v]", params)

	var topic, frameType string

	topic = stringParam(params[0])
	frameType = stringParam(params[1])

	var headers = make([]string, 0)
	for i := 2; i < len(params); i++ {
		h := stringParam(params[i])
		headers = append(headers, h)
	}

	frame := frame.New(
		frameType,
		"message", "test1",
	)

	//log.Debugf("%v", frame.Header)

	n.SendFrame(topic, *frame)
	log.Infof("Frame [%v] was sent from [%s]", frame, topic)
	return nil
}

func sendMessageMultiProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[params] sendMessageMulti=[%v]", params)
	return nil
}

func receiveFrameProcesser(signature string, n *stomp.NodeStomp) error {

	params := strings.Split(signature, ",")

	log.Debugf("[params] receiveFrame=[%v]", params)

	var topic, frameType string

	topic = stringParam(params[0])
	ch, err := n.Subscribe(topic)
	if err != nil {
		log.Errorf("topic=[%s]: %s", topic, err.Error())
		return err
	}

	frameType = stringParam(params[1])

	var headers = make([]string, 0)
	for i := 2; i < len(params); i++ {
		h := stringParam(params[i])
		headers = append(headers, h)
	}

	frame := frame.New(
		frameType,
		"message", "test1",
	)

	log.Debugf("%v", frame.Header)

	err = read(ch, frame)
	if err != nil {
		return err
	}
	return nil
}

func sendFrameMultiProcesser(signature string, n *stomp.NodeStomp) error {

	/*

		frameType, err := strconv.Atoi(params[3])
		if err != nil {
			return errors.New("Wrong number parameter.")
		}

		var fr *frame.Frame
		var headers = make([]string)

		for i := 0; i < numOfFrames; i++ {
			headers = append(headers, params[2])
			log.Debug("headers=%v", headers)
			fr = frame.New(
				params[1],
				headers)

		}

	*/

	/*	params := strings.Split(signature, ",")
		for i, p := range params {
			params[i] = stringParam(p)
		}
		log.Debugf("[params] sendFrameMulti=[%v]", params)

		frame := frame.New(
			params[1],
			params[2])

		// TODO:add return error to SendFrame()
		n.SendFrame(params[0], *frame)*/
	return nil
}

func waitMessageProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[command] waitMessage=[%s]", params)
	return nil
}

func waitFrameProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[command] waitFrame=[%s]", params)
	return nil
}

func waitMessageMultiProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[command] waitMessageMulti=[%s]", params)
	return nil
}

func waitFrameMultiProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[command] waitFrameMulti=[%s]", params)

	return nil
}
