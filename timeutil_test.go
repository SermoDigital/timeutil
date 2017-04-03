package timeutil

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

var days = func() (days [100]time.Weekday) {
	for i := range days {
		// Intn is [0, n)
		days[i] = time.Weekday(rand.Intn(int(time.Saturday) + 1))
	}
	return days
}()

func TestNext(t *testing.T) {
	const (
		start = 0000
		end   = 9999
	)

	// Forwards
	tm := time.Date(start, 1, 1, 0, 0, 0, 0, time.Local)
	for i := 0; tm.Year() <= end; i++ {
		tm = testNext(t, tm, i, Forward)
	}

	// Backwards
	tm = time.Date(end, 12, 31, 0, 0, 0, 0, time.Local)
	for i := 0; tm.Year() >= start; i++ {
		tm = testNext(t, tm, i, Backward)
	}
}

func testNext(t *testing.T, tm time.Time, i int, d Direction) time.Time {
	day := days[i%len(days)]
	n := i % (rand.Intn(10) + 1)
	next := NthWeekday(tm, day, n, d)

	var err error
	if next.Weekday() != day {
		err = errors.New("invalid weekday returned")
	}
	if n == 0 {
		if tm.Weekday() == day && !next.Equal(tm) {
			err = errors.New("n == 0, wkd == day, next != tm")
		}
	} else {
		switch d {
		case Forward:
			if !next.After(tm) {
				err = errors.New("moving forward but date isn't after tm")
			}
		case Backward:
			if !next.Before(tm) {
				err = errors.New("moving backward but date isn't prior to tm")
			}
		}
	}
	if err != nil {
		t.Fatalf("#%d: %s\n--> NthWeekday(%s [%s], %s, %d, %s): got %s [%s]",
			i, err,
			tfmt{tm}, tm.Weekday(), day, n, d,
			tfmt{next}, next.Weekday())
	}
	return next
}

func date(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestBounds(t *testing.T) {
	const (
		Jan = iota
		Feb
		Mar
		Apr
		May
		Jun
		Jul
		Aug
		Sep
		Oct
		Nov
		Dec
	)
	type year [12]struct{ start, end time.Time }

	for j, test := range [...]struct {
		year   int
		bounds year
	}{
		{
			year: 2017,
			bounds: year{
				Jan: {date("2017-01-02"), date("2017-01-29")},
				Feb: {date("2017-01-30"), date("2017-02-26")},
				Mar: {date("2017-02-27"), date("2017-04-02")},
				Apr: {date("2017-04-03"), date("2017-04-30")},
				May: {date("2017-05-01"), date("2017-05-28")},
				Jun: {date("2017-05-29"), date("2017-07-02")},
				Jul: {date("2017-07-03"), date("2017-07-30")},
				Aug: {date("2017-07-31"), date("2017-08-27")},
				Sep: {date("2017-08-28"), date("2017-10-01")},
				Oct: {date("2017-10-02"), date("2017-10-29")},
				Nov: {date("2017-10-30"), date("2017-11-26")},
				Dec: {date("2017-11-27"), date("2017-12-31")},
			},
			// TODO: 2016 and 2018-2024 since we have it in the .xlsx file.
		},
	} {
		for i, bs := range test.bounds {
			i++ // Months are 1-indexed, range is 0.

			// Every 4-4-5 month will be in [4, 26]
			r := rand.Intn((26-4)+1) + 4

			tm := time.Date(test.year, time.Month(i), r, 0, 0, 0, 0, time.UTC)
			start, end := Bounds(tm)
			t.Logf("%s: %s --> [%s, %s]\n",
				time.Month(i), tfmt{tm}, tfmt{start}, tfmt{end})
			if !start.Equal(bs.start) || !end.Equal(bs.end) {
				t.Fatalf("#%d: %s: wanted (%s, %s), got (%s, %s)",
					j,
					time.Month(i), tfmt{bs.start}, tfmt{bs.end},
					tfmt{start}, tfmt{end})
			}
		}
	}
}

type tfmt struct{ time.Time }

func (t tfmt) String() string {
	return t.Format("2006-01-02")
}
