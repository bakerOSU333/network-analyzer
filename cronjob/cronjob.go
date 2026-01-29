package cronjob

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// calculate timestring in cronjob format, ex: * * * * *
func calculateTimeString(timeInMinutes int) (timeStringInCronjob string) {
	// initialize the first and second position
	firstPosition := "*"
	secondPosition := "*"

	// if 24 hours
	if timeInMinutes == 1440 {
		firstPosition = "0"
		secondPosition = "0"
	}

	// if 60 minutes
	if timeInMinutes == 60 {
		firstPosition = "0"
	}

	// if less than 60 minutes and not 1 minute
	if timeInMinutes < 60 && timeInMinutes != 1 {
		firstPosition += "/" + strconv.Itoa(timeInMinutes)
	}

	// if more than 60 minutes and below 24 hours
	if timeInMinutes > 60 && timeInMinutes < 1440 {
		secondPosition += "/" + strconv.Itoa(timeInMinutes / 60)
	}

	finalString := firstPosition + " " + secondPosition + " * * *"

	return finalString
}

// ask for user time input, return time in minutes
func askForTimeInput() (timeInMinutes int) {
	var time int

	fmt.Println("Please choose how often you want to check your network latency")
	fmt.Println("Minimum: 1 minute, Maximum: 1 day")

	for {
		var inputTime string
		var timeMark string
		var min int

		fmt.Println("Possible input: 1 - 60 mins, 1 - 24 hrs. We don't support decimal hrs and mins")
		fmt.Printf("Your chosen time is: ")

		// get user input
		fmt.Scanf("%s %s", &inputTime, &timeMark)

		// if user input in mins, convert that to min
		if strings.Contains(timeMark, "mins") {
			minString := strings.Split(inputTime, "mins")[0]
			min, _ = strconv.Atoi(minString)
			// if user input in hrs, conver that to min
		} else if strings.Contains(timeMark, "hrs") {
			minString := strings.Split(inputTime, "hrs")[0]
			hours, _ := strconv.Atoi(minString)
			min = hours * 60
			// if user input is not hrs or mins, prompt again
		} else {
			fmt.Println("It's either minutes or hours")
			continue
		}

		// if min is between 1 and 60 or min is between 1 hr and 24 hrs and the number is whole
		if min >= 1 && min <= 60 || min > 60 && min % 60 == 0 && min <= 1440 {
			time = min
			break
			// if the time is not within the limit, prompt user again
		} else {
			fmt.Println("Please remember the maximum time")
		}
	}

	return time
}

// save the cronjob to the system, has 2 modes: add and remove
// add mode adds the cronjob to the system
// remove mode removes the cronjob from the system
func SaveCronJob(timeStringInCronjob string, WORKING_DIR string, mode string) error {
	// initialize the path of txt file
	cronTxtPath := WORKING_DIR + "/cronjob/cron.txt"

	// if the file exists, then delete it
	if _, err := os.Stat(cronTxtPath); err == nil {
		e := os.Remove(cronTxtPath)
		if e != nil {
			return e
		}
	}

	// create a new file
	file, openFileErr := os.OpenFile(cronTxtPath, os.O_RDWR | os.O_CREATE, 0777)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// list existing cronjobs
	crontabJobs, setupCronjobErr := exec.Command("crontab", "-l").Output()
	if setupCronjobErr != nil {
		return setupCronjobErr
	}

	if mode == "remove" {
		// convert []byte to string array
		cronjobArray := strings.Split(string(crontabJobs), "\n")

		// scanning target to remove
		scanningCronJob := WORKING_DIR + "/scanning"
		envEnvironment := "WORKING_DIR=" + WORKING_DIR

		// remove the line that contains the scanning target
		for index, cronjob := range cronjobArray {
			if strings.Contains(cronjob, scanningCronJob) {
				cronjobArray = append(cronjobArray[:index], cronjobArray[index + 1:]...)
			}

			if strings.Contains(cronjob, envEnvironment) {
				cronjobArray = append(cronjobArray[:index], cronjobArray[index + 1:]...)
			}
		}

		// convert string array to []byte
		joinedString := strings.Join(cronjobArray, "\n")

		// assign the []byte to the existing cronjobs
		crontabJobs = []byte(joinedString)
	}

	if mode == "add" {
		_, writeEnv := file.WriteString("WORKING_DIR=" + WORKING_DIR + "\n")
		if writeEnv != nil {
			return writeEnv
		}
	}

	// write existing cronjobs to the file
	_, writeExistingCronjobErr := file.Write(crontabJobs)
	if writeExistingCronjobErr != nil {
		return writeExistingCronjobErr
	}

	if mode == "add" {
		// create a new cronjob string
		cronjob := timeStringInCronjob + " " + WORKING_DIR + "/scanning >> /tmp/scanning.out 2>> /tmp/scanning.err" + "\n"

		// write the new cronjob to the file
		_, writeNewCronjobErr := file.WriteString(cronjob)
		if writeNewCronjobErr != nil {
			return writeNewCronjobErr
		}
	}

	// set up the cronjob to the system
	setupCronjob := exec.Command("crontab", cronTxtPath)
	setupCronjobErr = setupCronjob.Run()
	if setupCronjobErr != nil {
		return setupCronjobErr
	}

	// if no error then return nil
	return nil
}

// set up cronjob
func SetUpCronJob(WORKING_DIR string) error {
	// ask for time in minutes
	timeInMinutes := askForTimeInput()
	fmt.Printf("You chose %d minutes\n", timeInMinutes)

	// calculate time string in cronjob
	timeStringInCronjob := calculateTimeString(timeInMinutes)
	fmt.Printf("Your cronjob timestring is: %s\n", timeStringInCronjob)
	fmt.Println("Please allow the script to set the cronjob by clicking allow")

	// save the cronjob
	err := SaveCronJob(timeStringInCronjob, WORKING_DIR, "add")

	return err
}