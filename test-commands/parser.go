package test

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("test-parser")

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

// Process implements all test logic: parse line by line commands from file
// and call corresponding functions
func Process(exec CommandExecutor, fileWithCommands string) error {

	// openning file with commands and reading from it

	fd, err := os.Open(fileWithCommands)
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
			log.Errorf("file with commands: %s", err.Error())
			return err
		}
		if len(commandDeclaration) == 0 {
			log.Errorf("Empty line in File with commands=%s; ignored", fileWithCommands)
			continue
		}

		// parse a command declaration

		cmd, signature, err := parseCommand(commandDeclaration)
		if err != nil {
			log.Errorf("Parse command from file: %s; ignored", err.Error())
			continue
		}

		// finding this command in known commands and run correspoding processer
		err = exec.execute(cmd, signature)
		if err != nil {
			log.Errorf("Execute: %s [%s]", err.Error(), cmd)
			os.Exit(1)
		}
	}

	return nil
}
