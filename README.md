### Battery Mileage Monitor (BMM)

#### What is Battery Mileage 

I created this to measure how much % of discharge in a unit of time (hc'ded now for 5 mins). According to my measure ideally , 1% discharge rate (a.k.a Battery Mileage !!) for every 5 mins is a healthy one (on no specific conditions).

This works only on linux machines for now. 

BMM queries for the following files for every 5 mins record the capacity and then calcualtes the average of diff for the last 1 hr (sliding). 

```
- "/sys/class/power_supply/BAT0/status"
- "/sys/class/power_supply/BAT0/capacity"
```

#### Why I created this?

I'm using i3wm and wanted to plot this as part of my i3status. BMM writes the status to a file so its easy to create a i3status segment by exrtracting `cat` and `cut` commands !

**i3status**

```

```

#### TODO
- externalize the files 
- remove hc'ded frequency 5 mins
