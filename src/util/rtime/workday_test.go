package rtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEndTimeInWorkDays(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Input  time.Time
		Expect time.Time
	}

	cases := []testCase{
		{
			Input:  time.Date(2018, 9, 28, 12, 0, 0, 0, BeijingLocation),
			Expect: time.Date(2018, 10, 8, 23, 59, 59, 0, BeijingLocation),
		},
		{
			Input:  time.Date(2018, 9, 26, 12, 0, 0, 0, BeijingLocation),
			Expect: time.Date(2018, 9, 29, 23, 59, 59, 0, BeijingLocation),
		},
		{
			Input:  time.Date(2018, 10, 9, 12, 0, 0, 0, BeijingLocation),
			Expect: time.Date(2018, 10, 12, 23, 59, 59, 0, BeijingLocation),
		},
	}

	for i, c := range cases {
		ast.EqualValues(c.Expect.Unix(), EndTimeInWorkDays(c.Input, 3).Unix(), fmt.Sprint(i))
	}

	fmt.Println(SecondsAfterWorkDays(3))
}

func TestWorkDaysBetween(t *testing.T) {
	ast := assert.New(t)

	type testCase struct {
		Start    time.Time
		End      time.Time
		Error    error
		WorkDays int
	}

	cases := []testCase{
		{
			Start:    time.Date(2020, 9, 28, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 9, 8, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 0,
			Error:    ErrTimeRangeNotSupport,
		},
		{
			Start:    time.Date(2020, 9, 28, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 9, 28, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 1,
			Error:    nil,
		},
		{
			Start:    time.Date(2020, 9, 28, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 9, 28, 12, 0, 0, 0, BeijingLocation),
			WorkDays: 1,
			Error:    nil,
		},
		{
			Start:    time.Date(2020, 9, 28, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 10, 8, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 3,
		},
		{
			Start:    time.Date(2020, 9, 26, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 9, 29, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 3,
		},
		{
			Start:    time.Date(2020, 9, 11, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 9, 18, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 6,
		},
		{
			Start:    time.Date(2020, 10, 9, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 10, 12, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 3,
		},
		{
			Start:    time.Date(2020, 9, 9, 12, 0, 0, 0, BeijingLocation),
			End:      time.Date(2020, 11, 12, 23, 59, 59, 0, BeijingLocation),
			WorkDays: 43,
		},
	}

	for _, c := range cases {
		workDays, err := WorkDaysBetween(c.Start, c.End)
		ast.Equal(c.Error, err)
		ast.Equal(c.WorkDays, workDays)
	}
}
