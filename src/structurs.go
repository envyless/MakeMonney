package main

import "time"

type DailyData struct {
	date       time.Time
	lastPrice  int32
	deferPrice int32
	cap        int32
	highCap    int32
	lowCap     int32
	volume     int64
}
