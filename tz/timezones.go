// Package tz provides frequently-used timezones and utility functions.
package tz

import "time"

// MustParse calls time.LoadLocation, panicking if it returns a non-nil error.
func MustParse(zone string) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		panic(err)
	}
	return loc
}

// American timezones from §71.1 of
// https://www.gpo.gov/fdsys/pkg/CFR-2010-title49-vol1/pdf/CFR-2010-title49-vol1-part71.pdf
var (
	Atlantic       = MustParse("America/Puerto_Rico")
	Eastern        = MustParse("America/New_York")
	Central        = MustParse("America/Chicago")
	Mountain       = MustParse("America/Denver")
	Pacific        = MustParse("America/Los_Angeles")
	Alaska         = MustParse("America/Anchorage")
	HawaiiAleutian = MustParse("America/Adak")
	Samoa          = MustParse("Pacific/Pago_Pago")
	Chamorro       = MustParse("Pacific/Guam")
)
