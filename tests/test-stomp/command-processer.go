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
	"receiveFrame":     receiveFrameProcesser,
	"waitMessage":      waitMessageProcesser,
	"waitFrame":        waitFrameProcesser,
	"waitMessageMulti": waitMessageMultiProcesser,
	"waitFrameMulti":   waitFrameMultiProcesser,
	"sleep":            sleepProcesser,
	"dump":             dumpProcesser,
}

func process(n *stomp.NodeStomp) error {

	// openning file with commands and reading from it

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

		// deadind line by line, one line must contain only one command

		commandDeclaration = scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Errorf("globalOpt.FileWithCommands: %s", err.Error())
			return err
		}
		if len(commandDeclaration) == 0 {
			log.Errorf("Empty line in File with commands=%s; ignored", globalOpt.FileWithCommands)
			continue
		}

		// parse a command declaration

		cmd, signature, err := parseCommand(commandDeclaration)
		if err != nil {
			log.Errorf("Parse command from file: %s; ignored", err.Error())
			continue
		}

		// finding this command in known commands and run correspoding processer

		var function processorFunc
		var ok bool

		if function, ok = commands[cmd]; !ok {
			log.Errorf("Command [%s] from command file [%s] not found", cmd, globalOpt.FileWithCommands)
			continue
		}
		err = function(signature, n)
		if err != nil {
			log.Errorf("%s", err.Error())
			os.Exit(1)
		}

	}

	return nil
}
