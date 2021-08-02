package rtime

import (
	"errors"
	"time"
)

// 计算2018年工作日间隔秒数
func SecondsAfterWorkDays(workDays int) (seconds int64) {
	now := time.Now()
	return int64(EndTimeInWorkDays(now, workDays).Sub(now).Seconds())
}

func IsWorkDay(day time.Time) bool {
	y, m, d := day.Date()
	t := time.Date(y, m, d, 0, 0, 0, 0, BeijingLocation)
	if holidays[t] {
		return false
	}
	if workdays[t] {
		return true
	}
	switch t.Weekday() {
	case time.Saturday:
		return false
	case time.Sunday:
		return false
	default:
		return true
	}
}

// work days between [start, end]
func WorkDaysBetween(start, end time.Time) (workDays int, err error) {
	start, end = start.In(BeijingLocation), end.In(BeijingLocation)

	if start.Before(effectStartDate) ||
		end.After(effectEndDate) ||
		start.After(end) {
		err = ErrTimeRangeNotSupport
		return
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, BeijingLocation)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, BeijingLocation)

	//TODO 有个优化方案，将每个月的工作日和假日数量落库，可以快速计算，将循环次数减少一个量级
	for start.Before(end) {
		if IsWorkDay(start) {
			workDays++
		}
		start = start.AddDate(0, 0, 1)
	}

	if IsWorkDay(end) {
		workDays++
	}
	return
}

func EndTimeInWorkDays(start time.Time, workDays int) (end time.Time) {
	if workDays <= 0 {
		return start
	}
	if workDays >= 60 {
		// 保险起见不允许计算超过60个工作日
		return start
	}
	//if start.Year() != 2018 {
	//	// 暂时只支持2018年的起始日期
	//	return start
	//}
	y, m, d := start.Date()
	t := time.Date(y, m, d, 0, 0, 0, 0, BeijingLocation)
	i := 0
	for {
		if i >= workDays {
			break
		}
		t = t.AddDate(0, 0, 1)
		if holidays[t] {
			goto HAPPY
		}
		if workdays[t] {
			goto WORK
		}
		switch t.Weekday() {
		case time.Saturday:
			goto HAPPY
		case time.Sunday:
			goto HAPPY
		default:
			goto WORK
		}
	WORK:
		i++
	HAPPY:
		continue
	}
	y, m, d = t.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, BeijingLocation)
}

var (
	effectStartDate time.Time = time.Date(2018, 1, 1, 0, 0, 0, 0, BeijingLocation)
	effectEndDate   time.Time = time.Date(2020, 12, 31, 23, 59, 59, 0, BeijingLocation)

	ErrTimeRangeNotSupport error = errors.New("time range not support")
)

var holidays = map[time.Time]bool{
	// 2018
	time.Date(2018, 10, 1, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 2, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 3, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 4, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 5, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 6, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 10, 7, 0, 0, 0, 0, BeijingLocation): true,

	// 2019
	time.Date(2019, 1, 1, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 2, 4, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 2, 5, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 2, 6, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 2, 7, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 2, 8, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 4, 5, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 5, 1, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 5, 2, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 5, 3, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 5, 4, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 6, 7, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 9, 13, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2019, 10, 1, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2019, 10, 2, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2019, 10, 3, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2019, 10, 4, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2019, 10, 7, 0, 0, 0, 0, BeijingLocation): true,

	// 2020
	time.Date(2020, 1, 1, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 1, 24, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 1, 27, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 1, 28, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 1, 29, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 1, 30, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 4, 6, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 5, 1, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 5, 4, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 5, 5, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 6, 25, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 6, 26, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 1, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 2, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 5, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 6, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 7, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2020, 10, 8, 0, 0, 0, 0, BeijingLocation): true,

	// 2021
	// 元旦
	time.Date(2021, 1, 1, 0, 0, 0, 0, BeijingLocation): true,
	// 春节
	time.Date(2021, 2, 11, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 2, 12, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 2, 15, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 2, 16, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 2, 17, 0, 0, 0, 0, BeijingLocation): true,
	// 清明节
	time.Date(2021, 4, 5, 0, 0, 0, 0, BeijingLocation): true,
	// 劳动节
	time.Date(2021, 5, 3, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 5, 4, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 5, 5, 0, 0, 0, 0, BeijingLocation): true,
	// 端午节
	time.Date(2021, 6, 14, 0, 0, 0, 0, BeijingLocation): true,
	// 中秋节
	time.Date(2021, 9, 20, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 9, 21, 0, 0, 0, 0, BeijingLocation): true,
	// 国庆节
	time.Date(2021, 10, 1, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 10, 4, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 10, 5, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 10, 6, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 10, 7, 0, 0, 0, 0, BeijingLocation): true,
}

var workdays = map[time.Time]bool{
	// 2018
	time.Date(2018, 9, 29, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2018, 9, 30, 0, 0, 0, 0, BeijingLocation): true,

	// 2019
	time.Date(2019, 2, 2, 0, 0, 0, 0, BeijingLocation):   true,
	time.Date(2019, 2, 3, 0, 0, 0, 0, BeijingLocation):   true,
	time.Date(2019, 4, 28, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 5, 5, 0, 0, 0, 0, BeijingLocation):   true,
	time.Date(2019, 9, 29, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2019, 10, 12, 0, 0, 0, 0, BeijingLocation): true,

	// 2020
	time.Date(2020, 1, 19, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 2, 1, 0, 0, 0, 0, BeijingLocation):   true,
	time.Date(2020, 4, 26, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 5, 9, 0, 0, 0, 0, BeijingLocation):   true,
	time.Date(2020, 6, 28, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 9, 27, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2020, 10, 10, 0, 0, 0, 0, BeijingLocation): true,

	// 2021
	time.Date(2021, 2, 7, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2021, 2, 20, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 4, 25, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 5, 8, 0, 0, 0, 0, BeijingLocation):  true,
	time.Date(2021, 9, 18, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 9, 26, 0, 0, 0, 0, BeijingLocation): true,
	time.Date(2021, 10, 9, 0, 0, 0, 0, BeijingLocation): true,
}
