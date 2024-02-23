package main

import (
	"log"
	"time"
	_ "time/tzdata"

	"github.com/zsefvlol/timezonemapper"
)

func nowTime() time.Time {
	if len(lockTimeStr) > 0 {
		return lockTime
	} else {
		return time.Now()
	}
}

func latLonToTimezone(lat float64, lon float64) (string, *time.Location) {
	// Get timezone string from lat/long
	var timezone string = timezonemapper.LatLngToTimezoneString(lat, lon)
	// Should print "Timezone: Asia/Shanghai"
	// Load location from timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println("时区加载失败:", timezone, err)
	}
	// Parse time string with location
	// t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2010-01-01 00:00:00", loc)
	// Should print
	// 2010-01-01 00:00:00 +0800 CST
	// 2009-12-31 16:00:00 +0000 UTC
	// fmt.Println(t)
	// fmt.Println(t.UTC())
	return timezone, loc
}

func nameToTimezone(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Println("时区加载失败:", name, err)
	}
	return loc
}

func getUTCOffset(loc *time.Location) int {
	now := time.Now().In(loc)
	_, offset := now.Zone()
	offsetHours := offset / 3600
	return offsetHours // UTC +?
}
