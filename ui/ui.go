package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

const (
	InstallPackages = "1"
	ForwardPorts    = "2"
)

const mainTitle = "Lucy app"
const menuPrompt = "What do you want to do: "
const whichListPrompt = "Which list of apps You want to install?\n"

var reader = bufio.NewReader(os.Stdin)

func PrintMenu() {
	fmt.Println("")

	gameTitle := figure.NewFigure(mainTitle, "doom", true)
	gameTitle.Print()

	fmt.Println("1. Install packages")
	fmt.Println("2. Port forwards")
	fmt.Println("q. Quit")
}

func FetchMenuInput() string {
	fmt.Print(menuPrompt)
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}

func FetchListName() string {
	fmt.Print(whichListPrompt)
	input, _ := reader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}
