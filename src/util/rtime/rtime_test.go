package rtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetBJNextDayStartTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  1520853880,
			Expect: 1520870400,
		},
		{
			Input:  1520802000,
			Expect: 1520870400,
		},
		{
			Input:  1520780400,
			Expect: 1520784000,
		},
	}

	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJNextDayStartTs(c.Input))
	}
}

func TestGetBJCurrentDayStartTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  1520853880,
			Expect: 1520784000,
		},
		{
			Input:  1520802000,
			Expect: 1520784000,
		},
		{
			Input:  1520780400,
			Expect: 1520697600,
		},
	}

	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJCurrentDayStartTs(c.Input))
	}
}

func TestGetBJNextMondayStartTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  1525881900,
			Expect: 1526227200,
		},
		{
			Input:  1525791600,
			Expect: 1526227200,
		},
		{
			Input:  1525816800,
			Expect: 1526227200,
		},
		{
			Input:  1526011200,
			Expect: 1526227200,
		},
		{
			Input:  1526054400,
			Expect: 1526227200,
		},
		{
			Input:  1526169600,
			Expect: 1526227200,
		},
		{
			Input:  1525622399,
			Expect: 1525622400,
		},
		{
			Input:  1525622401,
			Expect: 1526227200,
		},
	}
	for i, c := range cases {
		ast.EqualValues(c.Expect, GetBJNextMondayStartTs(c.Input), fmt.Sprint(i))
	}
}

func TestGetBJCurrentDayString(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect string
	}

	cases := []testCase{
		{
			Input:  1525881900,
			Expect: "2018-05-10",
		},
		{
			Input:  1546865825,
			Expect: "2019-01-07",
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJCurrentDayString(c.Input))
	}
}

func TestGetBJCurrentMondayStartTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  1525881900,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1525791600,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1525816800,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1526011200,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1526054400,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1526169600,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1525622399,
			Expect: 1525622400 - 7*86400,
		},
		{
			Input:  1525622401,
			Expect: 1526227200 - 7*86400,
		},
		{
			Input:  1483200000, // 2017-01-01
			Expect: 1482681600, // 2016-12-26
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJCurrentMondayStartTs(c.Input))
	}
}

func TestDayHourMinStr(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect string
	}

	cases := []testCase{
		{
			Input:  60,
			Expect: "1分",
		},
		{
			Input:  30,
			Expect: "1分",
		},
		{
			Input:  864100,
			Expect: "10天2分",
		},
		{
			Input:  86400*7 + 3600*2 + 60*3 + 20,
			Expect: "7天2小时4分",
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, DayHourMinStr(c.Input))
	}
}

func TestDayHourStr(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect string
	}

	cases := []testCase{
		{
			Input:  60,
			Expect: "1小时",
		},
		{
			Input:  864100,
			Expect: "10天 1小时",
		},
		{
			Input:  86400*7 + 3600*2 + 60*3,
			Expect: "7天 3小时",
		},
		{
			Input:  86400*7 + 3600*23 + 60*3,
			Expect: "8天",
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, DayHourStr(c.Input, " "))
	}
}

func TestGetBJCurrentMonthStartTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  time.Date(2019, 1, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 1, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 2, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 2, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 2, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 2, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 11, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 11, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 9, 15, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 9, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 8, 15, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 8, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 1, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 1, 1, 0, 0, 0, 0, BeijingLocation).Unix(),
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJCurrentMonthStartTs(c.Input))
	}
}

func TestGetBJCurrentMonthEndTs(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  time.Date(2019, 1, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 1, 31, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 2, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 2, 28, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 2, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 2, 29, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 11, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 11, 30, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 9, 15, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 9, 30, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2019, 8, 15, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2019, 8, 31, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
		{
			Input:  time.Date(2020, 1, 1, 1, 1, 1, 1, BeijingLocation).Unix(),
			Expect: time.Date(2020, 1, 31, 23, 59, 59, 0, BeijingLocation).Unix(),
		},
	}
	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBJCurrentMonthEndTs(c.Input))
	}
}

func TestParseBJYMDHMSTime(t *testing.T) {
	ast := assert.New(t)
	type testCase struct {
		Input  string
		Expect int64
	}
	cases := []testCase{
		{
			Input:  "0000-00-00 00:00:00",
			Expect: -62135596800,
		},
		{
			Input:  "2018-07-06 14:13:16",
			Expect: 1530857596,
		},
	}
	for _, c := range cases {
		ret, _ := ParseBJYMDHMSTime(c.Input)
		ast.EqualValues(c.Expect, ret.Unix())
	}
}

func TestIsMinuteBetween(t *testing.T) {
	ast := assert.New(t)
	type testCase struct {
		Start  string
		End    string
		Now    string
		Expect bool
	}
	cases := []testCase{
		{
			Start:  "21:55",
			End:    "22:05",
			Now:    "22:06",
			Expect: false,
		},
		{
			Start:  "21:55",
			End:    "21:55",
			Now:    "21:55",
			Expect: true,
		},
		{
			Start:  "11:45",
			End:    "11:55",
			Now:    "11:50",
			Expect: true,
		},
		{
			Start:  "22:45",
			End:    "00:55",
			Now:    "22:50",
			Expect: true,
		},
		{
			Start:  "22:45",
			End:    "00:55",
			Now:    "00:50",
			Expect: true,
		},
		{
			Start:  "22:45",
			End:    "22:35",
			Now:    "22:50",
			Expect: true,
		},
		{
			Start:  "22:45",
			End:    "22:35",
			Now:    "22:36",
			Expect: false,
		},
	}
	nowDate := time.Now().Format("2006-01-02")
	for _, c := range cases {
		now, _ := time.ParseInLocation("2006-01-02 15:04", nowDate+" "+c.Now, BeijingLocation)
		ret := IsMinuteBetween(now, c.Start, c.End)
		ast.EqualValues(c.Expect, ret)
	}
}

func TestTimeMilliSecond(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  time.Duration
		Expect float64
	}
	cases := []testCase{
		{
			Input:  300 * time.Microsecond,
			Expect: 0.3,
		},
		{
			Input:  2345 * time.Microsecond,
			Expect: 2.345,
		},
		{
			Input:  time.Microsecond,
			Expect: 0.001,
		},
		{
			Input:  time.Second,
			Expect: 1000,
		},
	}
	for _, c := range cases {
		ast.InDelta(c.Expect, TransformMilliSecond(c.Input), 0.001)
	}
}

func TestGetBackOffInterval(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  uint32
		Expect time.Duration
	}

	step := time.Minute
	cases := []testCase{
		{
			Input:  0,
			Expect: 2 * step,
		},
		{
			Input:  1,
			Expect: 2 * step,
		},
		{
			Input:  2,
			Expect: 4 * step,
		},
	}

	for _, c := range cases {
		ast.EqualValues(c.Expect, GetBackOffInterval(c.Input, step, false))
	}
}

func TestMonthsBetween(t *testing.T) {
	ast := assert.New(t)
	startTime := time.Date(2021, 1, 1, 0, 0, 0, 0, BeijingLocation)
	endTime := time.Date(2021, 2, 1, 0, 0, 0, 0, BeijingLocation)

	ast.Equal(int32(1), MonthsBetween(startTime, endTime))
	ast.Equal(int32(-1), MonthsBetween(endTime, startTime))

	startTime = time.Date(2020, 12, 15, 0, 0, 0, 0, BeijingLocation)
	endTime = time.Date(2021, 2, 1, 0, 0, 0, 0, BeijingLocation)
	ast.Equal(int32(2), MonthsBetween(startTime, endTime))

	startTime = time.Date(2020, 12, 31, 23, 59, 0, 0, BeijingLocation)
	endTime = time.Date(2021, 1, 1, 0, 0, 0, 0, BeijingLocation)
	ast.Equal(int32(1), MonthsBetween(startTime, endTime))
}

func TestMilliStampToSecond(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  int64
		Expect int64
	}

	cases := []testCase{
		{
			Input:  1623396600,
			Expect: 1623396600,
		},
		{
			Input:  1623396600000,
			Expect: 1623396600,
		},
	}

	for _, c := range cases {
		ret := MilliStampToSecond(c.Input)
		t.Log(ret)
		ast.EqualValues(c.Expect, ret)
	}
}
