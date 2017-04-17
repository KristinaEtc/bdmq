package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/KristinaEtc/bdmq/stomp"
)

func stringParam(str string) (string, error) {
	if strings.Compare("\"", str[len(str)-1:])+strings.Compare("\"", str[:1]) == 0 {
		strWithoutFirstQuot := str[1:]
		str = strWithoutFirstQuot[:len(strWithoutFirstQuot)-1]
	}
	return str, nil
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
	topic, err := stringParam(signature)
	if err != nil {
		log.Errorf("Wrong [%s] parameter: %s", signature, err)
	}
	n.Subscribe(topic)
	return nil
}

func unsubscribeProcesser(signature string, n *stomp.NodeStomp) error {
	log.Debugf("[command] unsubscribe=[%s]", signature)
	//n.Unsubscribe(stringParam(signature))
	return nil
}

func sendMessageProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[params] sendMessage=[%v]", params)
	return nil
}

func sendFrameProcesser(signature string, n *stomp.NodeStomp) error {
	/*
		params := strings.Split(signature, ",")
		for i, p := range params {
			params[i] = stringParam(p)
		}
		log.Debugf("[params] sendFrame=[%v]", params)

		if len(params) < 4 {
			return errors.New("Not enought parameters.")
		}

		numOfFrames, err := strconv.Atoi(params[3])
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

		// TODO:add return error to SendFrame()
		n.SendFrame(params[0], *frame)
	*/
	return nil
}

func sendMessageMultiProcesser(signature string, n *stomp.NodeStomp) error {
	params := strings.Split(signature, ",")
	log.Debugf("[params] sendMessageMulti=[%v]", params)
	return nil
}

func sendFrameMultiProcesser(signature string, n *stomp.NodeStomp) error {
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
