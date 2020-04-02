package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// RubyDateLayout represents
	RubyDateLayout = "Mon Jan 02 15:04:05 MST 2006"
	// DateLayout represents
	DateLayout = "2006-01-02"
	// DateTimeLayout represents
	DateTimeLayout = "2006-01-02 15:04:05"
	// DateTimeLocationLayout represents
	DateTimeLocationLayout = "2006-01-02 15:04:05 -0700 MST"
	// DateTimeIDLayout represents
	DateTimeIDLayout = "20060102150405"
)

var (
	Local *time.Location
)

func init() {
	local, err := time.LoadLocation("Local") //服务器设置的时区
	if err != nil {
		fmt.Println(err)
	}
	Local = local
}

// GetWeekday returns current/specified weekday, 1-7
func GetWeekday(day ...time.Time) int {
	t := time.Now()
	if len(day) > 0 {
		t = day[0]
	}
	wd := int(t.Weekday())
	if wd == 0 {
		return 7
	}
	return wd
}

// ParseWeek returns weeks of input duration
func ParseWeek(from, to time.Time) (yws [][2]int) {
	var year, week int
	for {
		y, w := from.ISOWeek()
		if y != year || w != week {
			year = y
			week = w
			// fmt.Println(year, week)
			yws = append(yws, [2]int{year, week})
		}
		from = from.AddDate(0, 0, 1)
		if to.Sub(from) < 0 {
			break
		}
	}
	// fmt.Println(yws)
	return
}

// ContainWeek returns weeks of input duration
func ContainWeek(from, to time.Time, year, week int) bool {
	for {
		y, w := from.ISOWeek()
		if y == year || w == week {
			return true
		}
		from = from.AddDate(0, 0, 1)
		if to.Sub(from) < 0 {
			break
		}
	}
	return false
}

// Between tells whether t is between start and end
func Between(t, start, end time.Time) bool { return start.Before(t) && t.Before(end) }

// NowBetween tells whether t is between start and end
func NowBetween(start, end time.Time) bool { return Between(time.Now(), start, end) }

// NowBetweens tells whether t is between start and end
func NowBetweens(start, end string) bool {
	s, e1 := ParseDateTime(start)
	e, e2 := ParseDateTime(end)
	if e1 != nil || e2 != nil {
		return false
	}
	return NowBetween(s, e)
}

// ParseIDTimeInLocation returns
func ParseIDTimeInLocation(t string) (time.Time, error) {
	return time.ParseInLocation(DateTimeIDLayout, t, Local)
}

// ParseDateTimeInCST8 parse time in format 2006-01-02 15:04:05 -0700 MST
func ParseDateTimeInCST8(t string) (time.Time, error) {
	return time.Parse(DateTimeLocationLayout, t+" +0800 CST")
}

// ParseDateTime parse time in format 2006-01-02 15:04:05
func ParseDateTime(t string) (time.Time, error) {
	return time.Parse(DateTimeLayout, t)
}

// ParseDateTimeLocal parse local time in format 2006-01-02 15:04:05
func ParseDateTimeLocal(t string) (time.Time, error) {
	return time.ParseInLocation(DateTimeLayout, t, Local)
}

// ParseDateLocal parse local time in format 2006-01-02 15:04:05
func ParseDateLocal(t string) (time.Time, error) {
	return time.ParseInLocation(DateLayout, t, Local)
}

// ParseDate parse time in format 2006-01-02
func ParseDate(t string) (time.Time, error) {
	return time.Parse(DateLayout, t)
}

// ParseRubyTime parse time in format Mon Jan 02 15:04:05 MST 2006
func ParseRubyTime(t string) (time.Time, error) {
	return time.Parse(RubyDateLayout, t)
}

// Since returns duration from start in string format
func Since(start time.Time) string { return FormatDuration(time.Since(start)) }

// FormatDuration formats duration to string format
func FormatDuration(d time.Duration) string {
	v := d.String()
	arr := strings.Split(v, "h")
	if len(arr) == 2 {
		strh := arr[0]
		hours, err := strconv.Atoi(strh)
		if err != nil || hours < 24 {
			return v
		}
		d := hours / 24
		if d >= 30 {
			return fmt.Sprintf("%dM%dd%dh%s", d/30, d%30, hours%24, arr[1])
		}
		return fmt.Sprintf("%dd%dh%s", hours/24, hours%24, arr[1])
	}
	return v
}

// Now returns now time in string format 2006-01-02 15:04:05
func Now() string { return time.Now().Format(DateTimeLayout) }

// GetCost return the duration of time
func GetCost(start time.Time) time.Duration { return time.Now().Sub(start) }

// GetCostStr return the duration of time in string
func GetCostStr(start time.Time) string { return GetCost(start).String() }

// GetTimeStr return the current time in format 2006-01-02 15:04:05
func GetTimeStr(t time.Time) string { return t.Format(DateTimeLayout) }

// GetNow return the current time in int64
func GetNow() int64 { return time.Now().UnixNano() }

// GetTimestamp return the current/specified time in format "20060102150405"
func GetTimestamp(str ...string) string {
	t := time.Now().Unix()
	if len(str) > 0 {
		return time.Unix(t, 0).Format(str[0] + "_20060102150405")
	}
	return time.Unix(t, 0).Format("20060102150405")
}

// DaysOfMonth returns days of the specified year and month
func DaysOfMonth(year int, month int) int {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			return 30
		}
		return 31
	}

	if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
		return 29
	}
	return 28
}

// TodayStr returns today in format 2006-01-02
func TodayStr() string { return time.Now().Format(DateLayout) }

// Today returns today in format 2006-01-02
func Today() time.Time {
	today, _ := time.Parse(DateLayout, time.Now().Format(DateLayout))
	return today
}

// IsToday returns whether t is today
func IsToday(t time.Time) bool { return IsSameDay(time.Now(), t) }

// IsSameDay returns whether t1, t2 are the same day
func IsSameDay(t1, t2 time.Time) bool {
	return (t1.Year() == t2.Year()) && (t1.YearDay() == t2.YearDay())
}

// IsSameWeek returns whether t1, t2 are in the same week
func IsSameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return (y1 == y2) && (w1 == w2)
}

// IsCurrentWeek returns whether t is in this week
func IsCurrentWeek(t time.Time) bool {
	return IsSameWeek(t, time.Now())
}

// IsSameMonth returns whether t1, t2 are in the same month
func IsSameMonth(t1, t2 time.Time) bool {
	return (t1.Year() == t2.Year()) && (t1.Month() == t2.Month())
}

// SubDays returns days from after to before
func SubDays(before, after time.Time) int {
	for i := 0; i < 1000*365; i++ {
		if int(before.AddDate(0, 0, i).Sub(after)) == 0 {
			if int(after.Sub(before)) >= 0 {
				return i
			}
			return -i
		}
	}
	return 0
}

func subDays(before, after string) int {
	start, _ := time.Parse(DateLayout, before)
	end, _ := time.Parse(DateLayout, after)
	return SubDays(start, end)
}
