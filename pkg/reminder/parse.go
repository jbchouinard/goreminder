package reminder

import (
	"errors"
	"regexp"
	"time"
)

var regexTime = regexp.MustCompile(`(\d\d:\d\d) (.*)`)
var regexTomorrow = regexp.MustCompile(`tomorrow (\d\d:\d\d) (.*)`)
var regexDate = regexp.MustCompile(`(\d\d/\d\d \d\d:\d\d) (.*)`)
var regexYearDate = regexp.MustCompile(`(\d\d\d\d-\d\d-\d\d \d\d:\d\d) (.*)`)

var timeSpecs = []struct {
	regex *regexp.Regexp
	f     func(string, *time.Location) (time.Time, error)
}{
	{
		regexTomorrow,
		func(s string, loc *time.Location) (time.Time, error) {
			t, err := parseLocalTime(s, loc)
			tplus1 := t.AddDate(0, 0, 1)
			if err != nil {
				return time.Time{}, err
			}
			return tplus1, nil
		},
	},
	{
		regexDate,
		parseLocalDateTime,
	},
	{
		regexYearDate,
		parseLocalYearDateTime,
	},
	{
		regexTime,
		parseLocalTime,
	},
}

func parseLocalTime(s string, loc *time.Location) (time.Time, error) {
	localTime, err := time.Parse("15:04", s)
	if err != nil {
		return localTime, err
	}
	now := time.Now().In(loc)
	dateTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		localTime.Hour(), localTime.Minute(),
		0, 0, loc,
	)
	return dateTime, nil
}

func parseLocalDateTime(s string, loc *time.Location) (time.Time, error) {
	localTime, err := time.Parse("01/02 15:04", s)
	now := time.Now().In(loc)
	if err != nil {
		return localTime, err
	}
	dateTime := time.Date(
		now.Year(), localTime.Month(), localTime.Day(),
		localTime.Hour(), localTime.Minute(),
		0, 0, loc,
	)
	return dateTime, nil
}

func parseLocalYearDateTime(s string, loc *time.Location) (time.Time, error) {
	localTime, err := time.Parse("2006-01-02 15:04", s)
	if err != nil {
		return localTime, err
	}
	dateTime := time.Date(
		localTime.Year(), localTime.Month(), localTime.Day(),
		localTime.Hour(), localTime.Minute(),
		0, 0, loc,
	)
	return dateTime, nil
}

func parseSpec(s string, loc *time.Location) (time.Time, string, error) {
	for _, spec := range timeSpecs {
		if m := spec.regex.FindStringSubmatch(s); m != nil {
			dueTime, err := spec.f(m[1], loc)
			if err != nil {
				return time.Time{}, "", err
			}
			return dueTime, m[2], nil
		}
	}
	return time.Time{}, "", errors.New("not a valid reminder spec")
}
