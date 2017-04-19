package main

import (
	"errors"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"fmt"

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
	debug.PrintStack()
	return nil
}

func sleepProcesser(signature string, n *stomp.NodeStomp) error {
	log.Debugf("[command] sleep=[%s]", signature)

	sec, err := strconv.Atoi(signature)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * time.Duration(sec))
	return nil
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

func read(ch chan frame.Frame, prescriptiveFrame frame.Frame) error {

	timeout := time.NewTicker(time.Second * 10)

	log.Debug("read")
	for {
		select {
		case fr := <-ch:

			frameReceived++
			if globalOpt.ShowFrames {
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

func receiveFrameProcesser(signature string, n *stomp.NodeStomp) error {

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

		err := read(ch, *fr)
		if err != nil {
			return err
		}
	}
	return nil
}
