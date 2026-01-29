package ping

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PingStats struct {
	Min []string
	Avg []string
	Max []string
	Stddev []string
	TimeString []string
}

// function to scan ping for network latency then save the stats to ping/ping.txt
func RecordPingData(WORKING_DIR string) (scanningErr error) {
	// initialize path of the report file
	workingDirReport := WORKING_DIR + "/ping/ping.txt"

	// testing log
	fmt.Printf("Working dir ping: %s\n", workingDirReport)

	// scan using ping
	scanResult, pingScanningErr := exec.Command("/sbin/ping", "google.com", "-c", "10").Output()
	if pingScanningErr != nil {
		return pingScanningErr
	}

	// open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// create custom text for the report
	resultArray := strings.Split(string(scanResult), "\n")

	// ex output: round-trip min/avg/max/stddev = 26.483/30.290/37.375/3.926 ms | 2024-10-13 20:35:10
	finalString := resultArray[len(resultArray) - 2] + " | " + time.Now().Format("2006-01-02 15:04:05") + "\n"

	// write the text to the file
	_, writeFileErr := file.WriteString(finalString)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

// function to read the report and return the stats, ready for chart building
func ReadPingReport(reportPath string) (pingStats PingStats, err error) {
	report, openFileErr := os.ReadFile(reportPath)
	if openFileErr != nil {
		return PingStats{}, openFileErr
	}

	// extract data
	lines := strings.Split(string(report), "\n")
	lines = lines[:len(lines) - 1]

	pingStats = PingStats{}

	for _, line := range lines {
		lineSlice := strings.Split(line, "|")
		time := lineSlice[1]
		data := lineSlice[0]
		stats := strings.Split(strings.Split(data, "=")[1], "/")

		pingStats.Min = append(pingStats.Min, stats[0])
		pingStats.Avg = append(pingStats.Avg, stats[1])
		pingStats.Max = append(pingStats.Max, stats[2])
		pingStats.Stddev = append(pingStats.Stddev, strings.Split(stats[3], " ")[0])
		pingStats.TimeString = append(pingStats.TimeString, time)
	}

	return pingStats, nil
}