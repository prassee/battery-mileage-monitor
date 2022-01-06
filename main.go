/*
   BMM - Battery Mileage Monitor
   ==============================
   Measures discharge rate by quering `/sys/class/power_supply/BAT0/capacity` file
   for every 5 mins. Record `capacity` and calc the average of consequent diffs (for now hc'ed 12).
   The average is then flushed to file.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var capacity int = 0
var status string = ""
var diffs []int
var logFilePath = "/tmp/batMileageStatus.log"
var logger *log.Logger

func main() {
	f, err := os.OpenFile(logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger = log.New(f, "[BAT Discharge Mileage] ", log.LstdFlags)
	capacity, status = readBatStatus()
	for 1 > 0 {
		curCap, curStatus := readBatStatus()
		obsTime := time.Now().Format("2006-01-02 15:04")
		if strings.Contains(curStatus, "Discharging") {
			diff := capacity - curCap
			if len(diffs) == 12 {
				diffs = diffs[1:]
			}
			if diff > 0 { // FIXME - convert this to a FSM for better accuracy
				diffs = append(diffs, diff)
				avg := calcAvg()
				logger.Printf("time %v currentCapacity %v discharged %v avg %v \n", obsTime, curCap, diff, avg)
				saveAvg(avg)
			} else {
				logger.Printf("skipping calc \n start status \n time %v charge %v \n", obsTime, curCap)
			}
			capacity = curCap
			status = curStatus
		} else if strings.Contains(curStatus, "Charging") {
			logger.Printf("status charging \n %v - %v started at %v \n", obsTime, curStatus, curCap)
			//set to until it starts discharging
			saveAvg(0.0)
		}
		time.Sleep(5 * time.Minute)
	}
}

func readBatStatus() (int, string) {
	basePath := "/sys/class/power_supply/BAT0/"
	capacity, capFileMissing := ioutil.ReadFile(basePath + "/capacity")
	if capFileMissing != nil {
		logger.Printf("BAT Capacity file missing %v \n", capFileMissing.Error())
	}
	capParsed, statFileMissing := strconv.ParseInt(strings.TrimSpace(string(capacity)), 10, 64)
	if statFileMissing != nil {
		logger.Printf("BAT status file missing %v \n", statFileMissing.Error())
	}
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
	avg = sum / float64(len(diffs))
	return avg
}

func saveAvg(avg float64) {
	ioutil.WriteFile("/data/sysConfigs/batDrainAvg.txt", []byte(fmt.Sprintf("%v", avg)), os.FileMode(0644))
}
