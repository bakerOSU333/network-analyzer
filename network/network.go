package network

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"slices"
)

type NetworkData struct {
	ProcessName string
	ReceivedMB []string
	SentMB []string
	Time []string
}

// function to record network data to file
func RecordNetworkData(WORKING_DIR string) error {
	workingDirReport := WORKING_DIR + "/network/network.txt"

	// open the file for appending
	file, openFileErr := os.OpenFile(workingDirReport, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if openFileErr != nil {
		return openFileErr
	}
	defer file.Close()

	// scan using nettop
	networkcmd, err := exec.Command("nettop", "-l", "1", "-P", "-x").Output()
	networkcmd = networkcmd[:len(networkcmd) - 1]
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(networkcmd), "\n")
	lines = lines[1:]

	// using regex to extract data
	re := regexp.MustCompile(`\d+:\d+:\d+\.\d+\s+([\w\s\(\)\.]+)\.(\d+)\s+(\d+)\s+(\d+)`)

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)

		// get the network consumption in bytes
		receivedBytes, err := strconv.ParseFloat(matches[3], 64)
		if err != nil {
			fmt.Println(err)
		}

		// get the sent network in bytes
		sentBytes, err := strconv.ParseFloat(matches[4], 64)
		if err != nil {
			fmt.Println(err)
		}

		// convert from bytes to MB
		receivedMB := receivedBytes / float64(1000000)
		sentMB := sentBytes / float64(1000000)

		// convert to string
		receivedMBString := fmt.Sprintf("%.5f", receivedMB)
		sentMBString := fmt.Sprintf("%.5f", sentMB)

		finalResult := matches[1] + "." + matches[2] + " | " + receivedMBString + " | " + sentMBString + " | " + time.Now().Format("2006-01-02 15:04:05")

		// write the string to the file
		_, writeFileErr := file.WriteString(finalResult + "\n")
		if writeFileErr != nil {
			return writeFileErr
		}
	}

	return nil
}

// function to read the network report and return the stats, ready for chart building
func ReadNetworkData(WORKING_DIR string) (networkMap map[string]NetworkData, err error) {
	filePath := WORKING_DIR + "/network/network.txt"

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// get all the lines in slice format
	lines := strings.Split(string(file), "\n")
	lines = lines[:len(lines) - 1]

	var networkDataMap = make(map[string]NetworkData)

	for _, line := range lines {
		slice := strings.Split(line, " | ")
		processName := slice[0]

		// if the process is already in the map, update its network data
		if networkData, ok := networkDataMap[processName]; ok {
			// update existing slice
			receivedMB := append(networkData.ReceivedMB, slice[1])
			sentMB := append(networkData.SentMB, slice[2])
			time := append(networkData.Time, slice[3])

			// replace the old network data by a new network data
			networkDataMap[processName] = NetworkData{
				ProcessName: processName,
				ReceivedMB: receivedMB,
				SentMB: sentMB,
				Time: time,
			}
			// if the process is not in the map, add it
		} else {
			networkDataMap[processName] = NetworkData{
				ProcessName: processName,
				ReceivedMB: []string{slice[1]},
				SentMB: []string{slice[2]},
				Time: []string{slice[3]},
			}
		}
	}

	networkDataMap = RemoveUnactivatedNetworkData(networkDataMap)

	return networkDataMap, nil
}

// sort the map in descending order
func SortNetworkDataMap(networkDataMap map[string]NetworkData, sortedByReceivedData bool) (keysSortedInDesc []string) {
	// initialize a slice of strings containing all the keys sorted in descending order
	keysDesc := make([]string, 0, len(networkDataMap))

	// add all the keys of the map to the slice
	for key := range networkDataMap {
		keysDesc = append(keysDesc, key)
	}

	// sort the keys based on the requirement
	sort.SliceStable(keysDesc, func(i, j int) bool {
		var (
			totalMBI float64
			totalMBJ float64
			MBSliceI []string
			MBSliceJ []string
		)

		// if sorted by incoming network is true, sort by incoming network
		if sortedByReceivedData {
			MBSliceI = networkDataMap[keysDesc[i]].ReceivedMB
			MBSliceJ = networkDataMap[keysDesc[j]].ReceivedMB
			// if sorted by incoming network is false, sort by outgoing network
		} else {
			MBSliceI = networkDataMap[keysDesc[i]].SentMB
			MBSliceJ = networkDataMap[keysDesc[j]].SentMB
		}

		for _, value := range MBSliceI {
			valueFloat, _ := strconv.ParseFloat(value, 64)
			totalMBI += valueFloat
		}

		for _, value := range MBSliceJ {
			valueFloat, _ := strconv.ParseFloat(value, 64)
			totalMBJ += valueFloat
		}

		avgMBI := totalMBI / float64(len(MBSliceI))
		avgMBJ := totalMBJ / float64(len(MBSliceJ))

		return avgMBI > avgMBJ
	})

	return keysDesc
}

// function to get the top N keys in descending order
func GetTopDesc(keysSorted []string, topNumber int) (topKeysInDesc []string) {
	topKeys := make([]string, 0, 3)

	for i := 0; i < topNumber; i++ {
		topKeys = append(topKeys, keysSorted[i])
	}

	return topKeys
}

func RemoveUnactivatedNetworkData(networkDataMap map[string]NetworkData) (networkDataMapCleaned map[string]NetworkData) {
	for key, networkdata := range networkDataMap {
		allZeroSent := CheckFullZero(networkdata.SentMB)
		allZeroReceived := CheckFullZero(networkdata.ReceivedMB)
		if allZeroSent && allZeroReceived {
			delete(networkDataMap, key)
		}
	}

	return networkDataMap
}

func CheckFullZero(slice []string) bool {
	for _, value := range slice {
		if value != "0.00000" {
			return false
		}
	}

	return true
}

// some processes might have different length of time, this function will make them have the same length
func EqualizeTopKey(networkDataMap map[string]NetworkData, TopDesc []string, processNameLongestTime string) (networkDataMapCleaned map[string]NetworkData) {
	for index, time := range networkDataMap[processNameLongestTime].Time {
		for _, processName := range TopDesc {
			if index == len(networkDataMap[processName].Time) {
				networkData := networkDataMap[processName]
				networkData.Time = slices.Insert(networkData.Time, index, time)
				networkData.ReceivedMB = slices.Insert(networkData.ReceivedMB, index, "0.00000")
				networkData.SentMB = slices.Insert(networkData.SentMB, index, "0.00000")
				networkDataMap[processName] = networkData
			}

			if time != networkDataMap[processName].Time[index] {
				networkData := networkDataMap[processName]
				networkData.Time = slices.Insert(networkData.Time, index, time)
				networkData.ReceivedMB = slices.Insert(networkData.ReceivedMB, index, "0.00000")
				networkData.SentMB = slices.Insert(networkData.SentMB, index, "0.00000")
				networkDataMap[processName] = networkData
			}
		}
	}

	return networkDataMap
}

// find the processName that has the longest time slice
func FindLongestTime(TopDesc []string, networkDataMap map[string]NetworkData) (processName string) {
	longest := 0

	if len(networkDataMap[TopDesc[1]].Time) > len(networkDataMap[TopDesc[0]].Time) {
		longest = 1
	}

	if len(networkDataMap[TopDesc[2]].Time) > len(networkDataMap[TopDesc[longest]].Time) {
		longest = 2
	}

	return TopDesc[longest]
}