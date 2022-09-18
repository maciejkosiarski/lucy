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

type Forwarder struct {
	ConfigFileName string
	ClusterNames   []string
}

type Command struct {
	Cmd  string
	Args []string
}

func (f *Forwarder) ForwardPorts(cluster_name []string) {
	commands := f.providePortForwardCommands()
	endedCommands := make(chan Command)
	defer close(endedCommands)

	cmdsStarted := 0
	clStarted := 0
	for cn, cmds := range commands {
		if f.ifClusterSelected(cn) {
			clStarted++
			for _, cmd := range cmds {
				go f.runPortForward(cmd, endedCommands)
				cmdsStarted++
				fmt.Printf("[%s] %v %v running...\n", cn, cmd.Cmd, cmd.Args)
			}
		}
	}

	if clStarted == 0 {
		fmt.Printf("Cluster named: %s was not found\n", cluster_name)
		return
	}

	if cmdsStarted == 0 {
		fmt.Println("No commands to start!")
		return
	}

	for cmd := range endedCommands {
		go f.runPortForward(cmd, endedCommands)
		fmt.Println(f.formatSuccessLog(fmt.Sprintf("%v %v reforwarded!", cmd.Cmd, cmd.Args)))
	}

}

func (f *Forwarder) ifClusterSelected(cluster_name string) bool {
	if len(f.ClusterNames) == 0 {
		return true
	}
	for _, cn := range f.ClusterNames {
		if cn == cluster_name {
			return true
		}
	}
	return false
}

func (f *Forwarder) runPortForward(cmd Command, endedCommands chan Command) {
	command := exec.Command(cmd.Cmd, cmd.Args...)

	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()

	if err := command.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}

	go f.listenOutput(stdout, "-> %s")
	go f.listenOutput(stderr, "-> \033[38;5;196m%s\033[39;49m")

	err := command.Wait()

	if err != nil {
		fmt.Printf("\033[38;5;196m%v %v %s\033[39;49m\n", cmd.Cmd, cmd.Args, err)
	} else {
		endedCommands <- cmd
	}
}

func (f *Forwarder) listenOutput(ioReader io.ReadCloser, format string) {
	scanner := bufio.NewScanner(ioReader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Printf(format+"\n", scanner.Text())
	}
}

func (f *Forwarder) providePortForwardCommands() map[string][]Command {
	exPath, err := os.Executable()

	if err != nil {
		log.Fatal("Opening a installer config file failed. Exiting")
	}

	exPath = filepath.Dir(exPath)

	configFile, err := ioutil.ReadFile(exPath + "/" + f.ConfigFileName)

	if err != nil {
		log.Fatal(f.formatErrorLog(err))
	}

	var commands map[string][]Command
	err = yaml.Unmarshal([]byte(configFile), &commands)

	if err != nil {
		log.Fatal(f.formatErrorLog(err))
	}

	return commands
}

func (f *Forwarder) formatSuccessLog(log string) string {
	return fmt.Sprintf("\033[38;5;118m%s\033[39;49m", log)
}

func (f *Forwarder) formatErrorLog(err error) string {
	return fmt.Sprintf("\033[38;5;196m%s\033[39;49m", err)
}
