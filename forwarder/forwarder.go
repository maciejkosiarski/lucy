package forwarder

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const portForwardsConfig = "forwarder.yaml"

type Command struct {
	Cmd  string
	Args []string
}

func ForwardPorts() {
	commands := providePortForwardCommands()
	endedCommands := make(chan Command)

	for _, cmds := range commands {
		for _, cmd := range cmds {
			go runPortForward(cmd, endedCommands)
			fmt.Printf("%v %v running...\n", cmd.Cmd, cmd.Args)
		}
	}

	for cmd := range endedCommands {
		go runPortForward(cmd, endedCommands)
		fmt.Printf(formatSuccessLog(fmt.Sprintf("%v %v reforwarded!\n", cmd.Cmd, cmd.Args)))
	}
}

func runPortForward(cmd Command, endedCommands chan Command) {
	command := exec.Command(cmd.Cmd, cmd.Args...)

	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	if err := command.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}

	go listenOutput(stdout, "%s")
	go listenOutput(stderr, "\033[38;5;196m%s\033[39;49m")

	err := command.Wait()

	if err != nil {
		fmt.Printf("\033[38;5;196m%v %v %s\033[39;49m\n", cmd.Cmd, cmd.Args, err)
	} else {
		endedCommands <- cmd
	}
}

func listenOutput(ioReader io.ReadCloser, format string) {
	scanner := bufio.NewScanner(ioReader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Printf(format+"\n", scanner.Text())
	}
}

func providePortForwardCommands() map[string][]Command {
	exPath, err := os.Executable()

	if err != nil {
		log.Fatal("Opening a installer config file failed. Exiting")
	}

	exPath = filepath.Dir(exPath)

	configFile, err := ioutil.ReadFile(exPath + "/" + portForwardsConfig)

	if err != nil {
		log.Fatal(formatErrorLog(err))
	}

	var commands map[string][]Command
	err = yaml.Unmarshal([]byte(configFile), &commands)

	if err != nil {
		log.Fatal(formatErrorLog(err))
	}

	return commands
}

func formatSuccessLog(log string) string {
	return fmt.Sprintf("\033[38;5;118m%s\033[39;49m", log)
}

func formatErrorLog(err error) string {
	return fmt.Sprintf("\033[38;5;196m%s\033[39;49m", err)
}
