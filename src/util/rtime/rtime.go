package rtime

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	MysqlInsertFormat = "2006-01-02 15:04:05"
)

// 获得传入ts的北京时间当天0:0:0时间戳
func GetBJCurrentDayStartTs(input int64) int64 {
	return (input+28800)/86400*86400 - 28800
}

// 获得传入ts的北京时间当天23:59:59时间戳
func GetBJCurrentDayEndTs(input int64) int64 {
	return GetBJNextDayStartTs(input) - 1
}

func GetBJCurrentDayString(input int64) string {
	return time.Unix(input, 0).In(BeijingLocation).Format("2006-01-02")
}

func GetBJCurrentTimeString(input int64) string {
	return time.Unix(input, 0).In(BeijingLocation).Format(MysqlInsertFormat)
}

// 获得传入ts的北京时间第二天0:0:0时间戳
func GetBJNextDayStartTs(input int64) int64 {
	return GetBJCurrentDayStartTs(input) + 86400
}

// 获得传入时间的北京时间当月第一天的0:0:0的时间
func GetBJCurrentMonthStart(input time.Time) time.Time {
	year, month, _ := input.Date()

	t, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%v-%02d-01", year, int(month)), BeijingLocation)
	if err != nil {
		panic(err)
	}

	return t
}

// 获得传入时间的北京时间下月第一天的0:0:0的时间
func GetBJNextMonthStart(input time.Time) time.Time {
	year, month, _ := input.Date()
	month++

	if month == time.December {
		month = 1
		year++
	}

	t, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%v-%02d-01", year, int(month)), BeijingLocation)
	if err != nil {
		panic(err)
	}

	return t
}

// 获得传入ts的北京时间下一周周一的0:0:0时间戳
func GetBJNextMondayStartTs(input int64) int64 {
	t := time.Unix(input, 0).In(BeijingLocation)
	weekDay := t.Weekday()
	deltaDay := int64(0)
	if weekDay == time.Sunday {
		deltaDay = 1
	} else {
		deltaDay = 8 - int64(weekDay)
	}
	return GetBJCurrentDayStartTs(input) + deltaDay*86400
}

// 获得传入ts的北京时间本周周一的0:0:0时间戳
func GetBJCurrentMondayStartTs(input int64) int64 {
	t := time.Unix(input, 0).In(BeijingLocation)
	weekDay := t.Weekday()
	deltaDay := int64(0)
	if weekDay == time.Sunday {
		deltaDay = 7
	} else {
		deltaDay = int64(weekDay)
	}
	deltaDay--
	return GetBJCurrentDayStartTs(input) - deltaDay*86400
}

// 获得传入ts的北京时间当天某个整点(0~23)时间戳
func GetBJCurrentDayHourTs(input int64, hour int64) int64 {
	return GetBJCurrentDayStartTs(input) + hour*3600
}

// 获得传入ts的北京时间第二天某个整点(0~23)时间戳
func GetBJNextDayHourTs(input int64, hour int64) int64 {
	return GetBJNextDayStartTs(input) + hour*3600
}

var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)

const (
	timeTdFormat = "20060102"
)

// ads 会把一个整形20190304转成日期
func TransformToTd(now time.Time) (int64, error) {
	ret := now.Format(timeTdFormat)
	return strconv.ParseInt(ret, 10, 64)
}

func TDInMysql(createAt time.Time) (td uint64) {
	createAt = createAt.Add(time.Nanosecond * 500) // golang官方Mysql库在存储时基于微秒做了四舍五入
	return TD(createAt)
}

func TD(t time.Time) (td uint64) {
	dateFormat := t.In(BeijingLocation).Format("20060102")
	td, _ = strconv.ParseUint(dateFormat, 10, 64)
	return
}

// 将时间中的毫秒以下给截断掉
// '2019-07-04 23:59:59.512222'  => '2019-07-04 23:59:59'
func TrimMicroSeconds(t time.Time) time.Time {
	return t.Truncate(time.Second)
}

// 获得传入ts的北京时间当天0:0:0时间戳,测试环境180秒算一天，方便测试
func GetBJCurrentDayStartTsWithTest(input int64) int64 {
	section := int64(86400)
	d := section / 3
	return (input+d)/section*section - d
}

// 解析2006-01-02 15:04:05格式的北京时间
func ParseBJYMDHMSTime(input string) (t time.Time, err error) {
	if input == "" {
		err = errors.New("empty input")
		return
	}

	return time.ParseInLocation("2006-01-02 15:04:05", input, BeijingLocation)
}

// 时间间隔解析成"XX天XX小时XX分"文案
func DayHourMinStr(seconds int64) string {
	if seconds <= 0 {
		return "0分"
	}
	day := seconds / 86400
	hour := (seconds % 86400) / 3600
	minute := ((seconds % 86400) % 3600) / 60
	if seconds%86400%3600%60 > 0 {
		minute++
	}
	ret := ""
	if day > 0 {
		ret += fmt.Sprintf("%d天", day)
	}
	if hour > 0 {
		ret += fmt.Sprintf("%d小时", hour)
	}
	if minute > 0 {
		ret += fmt.Sprintf("%d分", minute)
	}

	return ret
}

// 时间间隔解析成"XX天XX小时"文案，上取整
func DayHourStr(seconds int64, sep string) string {
	if seconds <= 0 {
		return "0小时"
	}
	day := seconds / 86400
	hour := (seconds % 86400) / 3600
	if ((seconds % 86400) % 3600) > 0 {
		hour++
	}
	if hour%24 >= 0 {
		day += hour / 24
		hour = hour % 24
	}
	strS := make([]string, 0, 2)
	if day > 0 {
		strS = append(strS, fmt.Sprintf("%d天", day))
	}
	if hour > 0 {
		strS = append(strS, fmt.Sprintf("%d小时", hour))
	}

	return strings.Join(strS, sep)
}

// 获得传入ts的北京时间的本月第一天0:0:0的时间戳
func GetBJCurrentMonthStartTs(input int64) int64 {
	t := time.Unix(input, 0).In(BeijingLocation)
	currentYear, currentMonth, _ := t.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, BeijingLocation)
	return firstOfMonth.Unix()
}

// 获得传入ts的北京时间的本月的最后一天的23:59:59时间戳
func GetBJCurrentMonthEndTs(input int64) int64 {
	t := time.Unix(input, 0).In(BeijingLocation)
	currentYear, currentMonth, _ := t.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, BeijingLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0)
	return lastOfMonth.Unix() - 1
}

//判断时间是否在某个分钟级区间, 参数格式21:55,22:10
func IsMinuteBetween(now time.Time, begin, end string) bool {
	begins := strings.Split(begin, ":")
	ends := strings.Split(end, ":")
	if len(begins) != 2 || len(ends) != 2 {
		return false
	}
	nowDate := now.Format("2006-01-02")
	startTime, err := time.ParseInLocation("2006-01-02 15:04", nowDate+" "+begin, BeijingLocation)
	if err != nil {
		return false
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04", nowDate+" "+end, BeijingLocation)
	if err != nil {
		return false
	}
	//跨天了,end加一天
	if endTime.Before(startTime) {
		endTime = endTime.AddDate(0, 0, 1)
		if now.Before(startTime) {
			now = now.AddDate(0, 0, 1)
		}
	}
	if now.Unix() >= startTime.Unix() && now.Unix() < endTime.Unix()+60 {
		return true
	}
	return false
}

// 纳秒转换成毫秒
func TransformMilliSecond(d time.Duration) float64 {
	ms := d / time.Millisecond
	nsec := d % time.Millisecond
	return float64(ms) + float64(nsec)/1e6
}

var factor float64 = 2
var max float64 = 128

// 15 次
// 0-14
// 2/2/4/4/8/8/16/16/32/32/64/64/128/128/128 - 4/4/8/8/16/16/32/32/64/64/128/128/128/128/128
// 根据重试次数获取 下次时间间隔
func GetBackOffInterval(retries uint32, step time.Duration, fullJitter bool) time.Duration {
	retries /= 2

	r := 1.0
	if fullJitter {
		r = 1.0 + rand.Float64() // random number in [1..2]
	}
	m := r * factor * math.Pow(2, float64(retries))
	if m > max {
		m = max
	}
	d := time.Duration(int64(m)) * step
	return d
}

// 获得当天整点时间的时间戳，hour为[0,24]
func GetBJTodayHourTs(input int64, hour int8) (int64, error) {
	if hour < 0 || hour > 24 {
		return 0, errors.New("hour must in [0,24]")
	}
	todayStart := (input+28800)/86400*86400 - 28800
	hourGap := int64(hour) * 3600
	return todayStart + hourGap, nil
}

func MonthsBetween(startDate, endDate time.Time) int32 {
	return int32((endDate.Year()-startDate.Year())*12) + int32(endDate.Month()-startDate.Month())
}

func MilliStampToSecond(milliStamp int64) int64 {

	if milliStamp > 1000000000 && milliStamp < 3000000000 {
		return milliStamp
	}
	milliStamp = milliStamp / 1000
	needTime := time.Unix(milliStamp, 0)

	return needTime.Unix()
}
