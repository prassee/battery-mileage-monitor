# Battery Mileage Monitor (BMM)

## What is BMM 

BMM is a tiny program running in background to measure discharge in a unit of time (hc'ded now for 5 mins). On ideal conditions, 1-2 % discharge (a.k.a Battery Mileage !!) for every 5 mins is a healthy one.
This measure is independent of hardware (CPU, RAM etc) and OEM. 

**This works only on Linux machines for now**

BMM queries for the following files for every 5 mins record the capacity and then calcualtes the average of diff for the last 1 hr (sliding). 

    "/sys/class/power_supply/BAT0/status"
    "/sys/class/power_supply/BAT0/capacity"

## Why I created this?

I'm using i3wm and wanted to plot this as part of my i3status. BMM writes the status to a file so its easy to create a i3status segment by exrtracting `cat` and `cut` commands !

**Start BMM Daemon**

on system start run the binary __(build one from this repo)__

**i3blocks - segment**

    [power-drain]
    label=
    command=echo "$(cat </data/sysConfigs/batDrainAvg.txt> | cut -c 1-3)%"
    interval=10
    separator=true

#### TODO
- Implement TEST
- restart process around `sleep` & `charging` mode
- detect if charge  suddenly drops
- TOML driven config
- release binaries via github actions
- option to install as **systemd** service on linux
