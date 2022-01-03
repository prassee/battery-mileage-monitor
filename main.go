package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var capacity int = 0
var status string = ""
var diffs []int

func main() {
	capacity, status = readBatStatus()
	for 1 > 0 {
		curCap, curStatus := readBatStatus()
		obsTime := time.Now().Format("2006-01-02 15:04:05")
		if strings.Contains(curStatus, "Discharging") {
			diff := capacity - curCap
			if len(diffs) == 12 {
				diffs = diffs[1:]
			}
			if diff > 0 {
				diffs = append(diffs, diff)
				avg := calcAvg()
				fmt.Printf("time %v currentCapacity %v discharged %v avg %v \n", obsTime, curCap, diff, avg)
				saveAvg(avg)
			} else {
				fmt.Printf("skipping calc \n")
			}
			capacity = curCap
			status = curStatus
		} else if strings.Contains(curStatus, "Charging") {
			fmt.Printf("status charging at %v \n", obsTime)
		}
		time.Sleep(5 * time.Minute)
	}
}

func readBatStatus() (int, string) {
	basePath := "/sys/class/power_supply/BAT0/"
	capacity, _ := ioutil.ReadFile(basePath + "/capacity")
	capParsed, _ := strconv.ParseInt(strings.TrimSpace(string(capacity)), 10, 64)
	cap := int(capParsed)
	status, _ := ioutil.ReadFile(basePath + "/status")
	statusStr := string(status)
	return cap, statusStr
}

func calcAvg() (avg float64) {
	sum := 0.0
	for i := range diffs {
		sum += float64(diffs[i])
	}
	return sum / float64(len(diffs))
}

func saveAvg(avg float64) {
	ioutil.WriteFile("/data/sysConfigs/batDrainAvg.txt", []byte(fmt.Sprintf("%v", avg)), os.FileMode(0644))
}
