package utils

import (
	"fmt"
	"time"
)

var (
	days = map[string]string{
		"Sunday":    "Minggu",
		"Monday":    "Senin",
		"Tuesday":   "Selasa",
		"Wednesday": "Rabu",
		"Thursday":  "Kamis",
		"Friday":    "Jumat",
		"Saturday":  "Sabtu",
	}

	months = map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}
	// Loc is the pre-loaded location for Asia/Jakarta
	Loc, _ = time.LoadLocation("Asia/Jakarta")
)

// FormatIndonesian formats a time.Time to "Day, 02 Month Year - 15:04 WIB" in Bahasa Indonesia
func FormatIndonesian(t time.Time) string {
	if t.IsZero() {
		return "-"
	}

	if Loc != nil {
		t = t.In(Loc)
	}

	day := days[t.Weekday().String()]
	month := months[t.Month().String()]
	date := t.Format("02")
	year := t.Format("2006")
	timeStr := t.Format("15:04")

	return fmt.Sprintf("%s, %s %s %s - %s WIB", day, date, month, year, timeStr)
}

// FormatIndonesianPtr formats a *time.Time to string, returning "-" if nil
func FormatIndonesianPtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return FormatIndonesian(*t)
}
