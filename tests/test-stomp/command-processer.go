package main

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/KristinaEtc/bdmq/stomp"
)

type processorFunc func(string, *stomp.NodeStomp) error

func parseCommand(commandDeclaration string) (command, signature string, err error) {

	splitted := strings.SplitN(commandDeclaration, "(", 2)
	if len(splitted) < 1 {
		return "", "", errors.New("Wrong command declaration [(]")
	}
	command = splitted[0]
	if len(splitted) < 2 {
		return "", "", errors.New("Wrong command declaration [)]")
	}
	signature = splitted[1][:len(splitted[1])-1]

	return
}

var commands = map[string]processorFunc{
	"subscribe":        subscribeProcesser,
	"unsubscribe":      unsubscribeProcesser,
	"sendMessage":      sendMessageProcesser,
	"sendFrame":        sendFrameProcesser,
	"sendMessageMulti": sendMessageMultiProcesser,
	"sendFrameMulti":   sendFrameMultiProcesser,
	"waitMessage":      waitMessageProcesser,
	"waitFrame":        waitFrameProcesser,
	"waitMessageMulti": waitMessageMultiProcesser,
	"waitFrameMulti":   waitFrameMultiProcesser,
	"sleep":            sleepProcesser,
	"dump":             dumpProcesser,
}

func process(n *stomp.NodeStomp) error {

	fd, err := os.Open(globalOpt.FileWithCommands)
	if err != nil {
		log.Errorf("Error: [%s]", err.Error())
		return err
	}

	defer func() {
		closed := fd.Close()
		if closed != nil {
			log.Errorf("FileWithCommand: %s", closed.Error())
		}
	}()

	var commandDeclaration string

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		commandDeclaration = scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Errorf("globalOpt.FileWithCommands: %s", err.Error())
			return err
		}
		if len(commandDeclaration) == 0 {
			log.Errorf("Empty line in File with commands=%s; ignored", globalOpt.FileWithCommands)
			continue
		}

		var cmdFound bool

		cmd, signature, err := parseCommand(commandDeclaration)
		if err != nil {
			log.Errorf("Parse command from file: %s; ignored", err.Error())
			continue
		}
		log.Debugf("cmd=[%s], signature=[%s]", cmd, signature)
		for command, function := range commands {
			if strings.Compare(command, cmd) == 0 {
				log.Debugf("Got it=[%s]", cmd)
				cmdFound = true
				function(signature, n)
				break
			}
		}

		if cmdFound != true {
			log.Errorf("Command [%s] from command file [%s] not found", cmd, globalOpt.FileWithCommands)
		}
	}

	return nil
}
