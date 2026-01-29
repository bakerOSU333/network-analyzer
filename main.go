package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bakerOSU333/network-analyzer/cronjob"
	"github.com/bakerOSU333/network-analyzer/network"
	"github.com/bakerOSU333/network-analyzer/ping"
	"github.com/bakerOSU333/network-analyzer/speedtest"
)

func main() {
	// get working directory by env
	WORKING_DIR := os.Getenv("WORKING_DIR")

	// if env is not available, set the working directory to current working dir
	if WORKING_DIR == "" {
		pwdCommand := exec.Command("pwd")

		pwdByte, err := pwdCommand.Output()
		if err != nil {
			log.Fatal(err)
		}

		WORKING_DIR = strings.TrimSpace(string(pwdByte))
	}

	if len(os.Args) > 1 && os.Args[1] != "-a" {
		fmt.Println("Invalid command line")
		os.Exit(1)
	}

	// welcome message
	fmt.Println("Welcome to the Network Analyzer!")
	fmt.Println("----------------------------------")

	// get into advanced mode by using -a
	if len(os.Args) > 1 && os.Args[1] == "-a" {
		RunTerminal(WORKING_DIR)
	}
	// if user does not input custom -a flag, then switch to collecting data mode

	// list existing cronjobs
	cronjobList, listCronjobErr := exec.Command("crontab", "-l").Output()
	if listCronjobErr != nil {
		log.Fatal(listCronjobErr)
	}

	// if we already set up cronjob, then we will collect data
	if strings.Contains(string(cronjobList), "scanning") {
		fmt.Println("We already set up cronjob")
		fmt.Println("Perform automatic scanning, record network data and record Upload Speed and Download Speed")

		err := ping.RecordPingData(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}

		recordErr := network.RecordNetworkData(WORKING_DIR)
		if recordErr != nil {
			log.Fatal(recordErr)
		}

		speedtestErr := speedtest.RecordSpeedTestData(WORKING_DIR)
		if speedtestErr != nil {
			log.Fatal(speedtestErr)
		}

		os.Exit(1)
	} else {
		fmt.Println("Looks like we haven't set up the cronjob, let's do it now!")
		err := cronjob.SetUpCronJob(WORKING_DIR)
		if err != nil {
			log.Fatal(err)
		}
	}
}