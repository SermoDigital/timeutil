// Package timeutil provides useful routimes for manipulating time.Time values.
// In particular, it includes routines for 4-4-5 (accounting) years and ISO
// 8601 dates.
package timeutil

import "time"

// TODO: provide different 'schemes', like 4-5-4 and 5-4-4.

// Direction tells NthWeekday whether to move forwards or backwards to find
// dates.
type Direction uint8

const (
	// Forward moves NthWeekday forward toward the end of the year.
	Forward Direction = iota

	// Backward moves NthWeekday backward toward the start of the year.
	Backward

	// Either moves NthWeekday toward the closest instance of the day.
	Either
)

//go:generate stringer -type=Direction

func absmin(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if a < b {
		return a
	}
	return b
}

// NthWeekday returns the nth instance of the given weekday, starting at t. If
// n == 0 and t.Weekday() == day, t will not be changed.
func NthWeekday(t time.Time, day time.Weekday, n int, d Direction) time.Time {
	if n < 0 {
		return t
	}

	wkd := t.Weekday()
	if wkd == day {
		if d == Backward {
			return t.AddDate(0, 0, n*-7)
		}
		// d == Either has no sane representation, so just move forward.
		return t.AddDate(0, 0, n*7)
	}

	incr := int(day - wkd)
	switch d {
	case Forward:
		incr += 7
		if n > 0 {
			n = (n - 1) * 7
		}
	case Backward:
		incr -= 7
		if n > 0 {
			n = (n - 1) * -7
		}
	case Either:
		incr = absmin(incr+7, incr-7)
		if incr > 0 {
			n = (n - 1) * 7
		} else {
			n = (n - 1) * -7
		}
	}

	incr %= 7
	incr += n
	return t.AddDate(0, 0, incr)
}

// Next returns the next instance of day, advancing even if t.Weekday() == day.
func Next(t time.Time, day time.Weekday) time.Time {
	return NthWeekday(t, day, 1, Forward)
}

// Previous returns the last instance of day, receding even if
// t.Weekday() == day.
func Previous(t time.Time, day time.Weekday) time.Time {
	return NthWeekday(t, day, 1, Either)
}

// Closest returns the closest instance of day, advancing even if
// t.Weekday() == day.
func Closest(t time.Time, day time.Weekday) time.Time {
	return NthWeekday(t, day, 1, Either)
}

var months = [54]uint8{
	0,

	// Q1
	1, 1, 1, 1, // January
	2, 2, 2, 2, // February
	3, 3, 3, 3, 3, // March

	// Q2
	4, 4, 4, 4, // April
	5, 5, 5, 5, // May
	6, 6, 6, 6, 6, // June

	// Q3
	7, 7, 7, 7, // July
	8, 8, 8, 8, // August
	9, 9, 9, 9, 9, // September

	// Q4
	10, 10, 10, 10, // October
	11, 11, 11, 11, // November
	12, 12, 12, 12, 12, // December

	// Leap year. Use a fake December.
	13,
}

// Month returns the ISO month number corresponding with the given date using
// the 4-4-5 calendar. For example, in 2017 2 January (week 01) is in month 1,
// while 27 February (week 09) is in month 3.
func Month(t time.Time) time.Month {
	_, week := t.ISOWeek()
	return time.Month(months[week])
}

// Quarter returns the quarter for the given date using the 4-4-5 calendar.
func Quarter(t time.Time) int {
	_, week := t.ISOWeek()
	// ISO years are divided into 4 13-week quarters.
	return (week / 13) + 1
}

// Bounds returns the start and end date of the given month using the 4-4-5
// calendar.
func Bounds(t time.Time) (start, end time.Time) {
	// The week containing 4 January is the first week of the ISO 8061 year.
	// (Formally, it's the week with the year's first Thursday, but this is
	// quicker.)
	jan4 := time.Date(t.Year(), 1, 4, 0, 0, 0, 0, t.Location())

	// The first date of an ISO 8601 year is the first Monday after or on
	// 4 January.
	start = NthWeekday(jan4, time.Monday, 0, Either)

	month := Month(t)
	start = start.AddDate(0, 0, daysBefore[month])

	// Subtract 1 since a month's bounds are [start, end).
	switch {
	// 9 of 12 months are 4 weeks, or 28 days.
	default:
		end = start.AddDate(0, 0, 28-1)
	// In 4-4-5 scheme every third month is 5 weeks, or 35 days.
	case month%3 == 0:
		end = start.AddDate(0, 0, 35-1)
	// On leap years December is 6 weeks, or 42 days.
	case month == 13:
		end = start.AddDate(0, 0, 42-1)
	}
	return start, end
}

// daysBefore[m] counts the number of days prior to the month, assuming a 4-4-5
// calendar year.
var daysBefore = [...]int{
	0,
	0,
	28,
	28 + 28,
	28 + 28 + 35,
	28 + 28 + 35 + 28,
	28 + 28 + 35 + 28 + 28,
	28 + 28 + 35 + 28 + 28 + 35,
	28 + 28 + 35 + 28 + 28 + 35 + 28,
	28 + 28 + 35 + 28 + 28 + 35 + 28 + 28,
	28 + 28 + 35 + 28 + 28 + 35 + 28 + 28 + 35,
	28 + 28 + 35 + 28 + 28 + 35 + 28 + 28 + 35 + 28,
	28 + 28 + 35 + 28 + 28 + 35 + 28 + 28 + 35 + 28 + 28,
	28 + 28 + 35 + 28 + 28 + 35 + 28 + 28 + 35 + 28 + 28,
}
