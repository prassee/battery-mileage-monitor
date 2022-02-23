package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getBatteryStatus() (int, string) {
	basePath := "/sys/class/power_supply/BAT0/"
	capacity, capFileMissing := ioutil.ReadFile(basePath + "/capacity")
	if capFileMissing != nil {
		logger.Printf("BAT Capacity file missing %v \n", capFileMissing.Error())
	}
	capParsed, statParserError := strconv.ParseInt(strings.TrimSpace(string(capacity)), 10, 64)
	if statParserError != nil {
		logger.Printf("BAT status cannot be parsed %v \n", statParserError.Error())
	}
	cap := int(capParsed)
	status, _ := ioutil.ReadFile(basePath + "/status")
	statusStr := string(status)
	return cap, statusStr
}


// Function `calcAvg` - calculates the average discharge rate 
func calcAvg() (avg float64) {
	sum := 0.0
	for i := range diffs {
		sum += float64(diffs[i])
	}
	avg = sum / float64(len(diffs))
	return avg
}

func saveAvg(avg float64) {
	ioutil.WriteFile(bmmPath, []byte(fmt.Sprintf("%v", avg)), os.FileMode(0644))
}
