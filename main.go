/*
   BMM - Battery Mileage Monitor
   ==============================
   Measures discharge rate by quering `/sys/class/power_supply/BAT0/capacity` file
   for every 5 mins. Record `capacity` and calc the average of consequent diffs (for now hc'ed 12).
   The average is then flushed to file.
*/
package main

import (
	"log"
	"os"
	"strings"
	"time"
)

var capacity int = 0
var status string = ""
var diffs []int
var logFilePath = "/tmp/batMileageStatus.log"
var bmmPath = "/data/sysConfigs/batDrainAvg.txt"
var freq = 1 * time.Minute
var logger *log.Logger

func initLogger() *os.File {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger = log.New(f, "[BAT Discharge Mileage] ", log.LstdFlags)
	return f
}

func main() {
	loggerFile := initLogger()
	defer loggerFile.Close()
	capacity, status = getBatteryStatus()
	for 1 > 0 {
		curCap, curStatus := getBatteryStatus()
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
		time.Sleep(freq)
	}
}
