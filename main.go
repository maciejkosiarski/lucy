package main

import (
	"fmt"

	"github.com/maciejkosiarski/lucy/forwarder"
	"github.com/maciejkosiarski/lucy/installer"
	"github.com/maciejkosiarski/lucy/ui"
)

func main() {
	ui.PrintMenu()
	action := ui.FetchMenuInput()

	if action == ui.InstallPackages {
		installer.PrintAvailableLists()

		packagesToInstall := ui.FetchListName()

		installer.InstallPackages(packagesToInstall)
	}

	if action == ui.ForwardPorts {
		forwarder.ForwardPorts()
	}

	if action == "q" {
		fmt.Println("Bye!")
		return
	}
}
