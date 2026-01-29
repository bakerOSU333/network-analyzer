package main

import (
	"fmt"
	"github.com/bakerOSU333/network-analyzer/chart"
	"github.com/bakerOSU333/network-analyzer/cronjob"
	"github.com/nexidian/gocliselect"
	"log"
	"os"
	"os/exec"
)

func RunTerminal(WORKING_DIR string) {
	// create a terminal menu for the user
	menu := gocliselect.NewMenu("What do you want to do?")

	// create option for the user
	menu.AddItem("Cronjob options", "cronjob options")
	menu.AddItem("Show process network usage chart", "network usage chart")
	menu.AddItem("Show network latency chart", "network latency chart")
	menu.AddItem("Quit", "quit")

	for {
		clearTerminal()

		// get the choice from the user
		choice := menu.Display()

		switch choice {
		case "cronjob options":
			cronJobOptions(WORKING_DIR)
		case "network usage chart":
			err := chart.CreateNetworkChart(WORKING_DIR)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		case "network latency chart":
			chart.CreateSpeedtestChart()
			chart.CreatePingChart()
			os.Exit(1)
		case "quit":
			fmt.Println("Goodbye! See you later")
			os.Exit(1)
		}
	}
}

func cronJobOptions(WORKING_DIR string) {
	menu := gocliselect.NewMenu("Cronjob options")

	menu.AddItem("Edit cronjob time", "edit cronjob")
	menu.AddItem("Remove cronjob completely", "remove cronjob")
	menu.AddItem("Come back", "come back")

	clearTerminal()
	choice := menu.Display()

	switch choice {
	case "edit cronjob":
		// remove the current cronjob
		cronjob.SaveCronJob("", WORKING_DIR, "remove")
		// add a new cronjob
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	case "remove cronjob":
		cronjob.SaveCronJob("", WORKING_DIR, "remove")
	case "come back":
		fmt.Println("Coming back")
	}
}

func clearTerminal() {
	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()
}