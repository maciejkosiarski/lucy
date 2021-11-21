package installer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const configFile = "installer.yaml"
const snapCommand = "snap"
const aptCommand = "apt"

type Config struct {
	Snap []Snap
	Apt  []Apt
}

type Snap struct {
	Package  string
	Args     []string
	Commands []Command
}

type Apt struct {
	Package  string
	Args     []string
	Commands []Command
}

type Command struct {
	Cmd  string
	Args []string
}

type CommandResult struct {
	Package string
	Status  int //0 - success, 1 - error
}

func InstallPackages(listName string) {
	config := provideConfig()

	list := findListByName(listName, config)

	// packagesQuanitity := 0
	// packagesQuanitity = packagesQuanitity + len(list.Apt)
	// packagesQuanitity = packagesQuanitity + len(list.Snap)

	// installed := make(chan CommandResult, packagesQuanitity)

	for _, apt := range list.Apt {
		installApt(apt)
	}

	for _, snap := range list.Snap {
		installSnap(snap)
	}

	// for result := range installed {
	// 	if result.Status == 0 {
	// 		formatError(result.Package + " installed!")
	// 	} else {
	// 		formatSuccess(result.Package + " error detected!")
	// 	}
	// }
}

func installSnap(snap Snap) {
	args := []string{"install", snap.Package}
	args = append(args, snap.Args...)

	cmd := exec.Command(snapCommand, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// results <- CommandResult{snap.Package, 1}
		log.Fatal(formatErrorLog(err))
	}

	fmt.Println(string(output))

	// results <- CommandResult{snap.Package, 0}

	for _, cmd := range snap.Commands {
		runCmd(cmd)
	}
}

func installApt(apt Apt) {
	args := []string{"install", apt.Package}
	args = append(args, apt.Args...)

	cmd := exec.Command(aptCommand, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// results <- CommandResult{apt.Package, 1}
		log.Fatal(formatErrorLog(err))
	}

	fmt.Println(string(output))

	// results <- CommandResult{apt.Package, 0}

	for _, cmd := range apt.Commands {
		runCmd(cmd)
	}
}

func runCmd(command Command) {
	cmd := exec.Command(command.Cmd, command.Args...)
	stdout := cmd.Stdout
	fmt.Printf("stdout: %v\n", stdout)
}

func findListByName(name string, config map[string]Config) Config {

	list, exist := config[name]

	if !exist {
		err := fmt.Errorf("%v list does't exist", name)
		log.Fatal(formatErrorLog(err))

	}

	return list
}

func PrintAvailableLists() {
	config := provideConfig()

	i := 1
	for listName, value := range config {
		fmt.Printf("%v. %v:\n", i, listName)
		i++
		for index, item := range value.Snap {
			if index == 0 {
				fmt.Printf("	### snap ###\n")
			}
			fmt.Printf("	%v. %v\n", index+1, item)
		}
		for index, item := range value.Apt {
			if index == 0 {
				fmt.Printf("	### apt ###\n")
			}
			fmt.Printf("	%v. %v\n", index+1, item)
		}
	}
}

func provideConfig() map[string]Config {
	exPath, err := os.Executable()

	if err != nil {
		log.Fatal("Opening a installer config file failed. Exiting")
	}

	exPath = filepath.Dir(exPath)

	configFile, err := ioutil.ReadFile(exPath + "/" + configFile)

	if err != nil {
		log.Fatal(formatErrorLog(err))
	}

	var config map[string]Config
	err = yaml.Unmarshal([]byte(configFile), &config)

	if err != nil {
		log.Fatal(formatErrorLog(err))
	}

	return config
}

func formatErrorLog(err error) string {
	return fmt.Sprintf("\033[38;5;196m%s\033[39;49m", err)
}

// func formatSuccess(msg string) string {
// 	return fmt.Sprintf("\033[38;5;118m%s\033[39;49m", msg)
// }

// func formatError(msg string) string {
// 	return fmt.Sprintf("\033[38;5;196m%s\033[39;49m", msg)
// }
