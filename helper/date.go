package helper

import (
	"errors"
	"github.com/patcharp/golib/v2/util"
	"regexp"
	"strings"
	"time"
)

type Date struct {
	date time.Time
}

func (d *Date) ToDateTime() time.Time {
	return d.date
}

func (d *Date) Set(t time.Time) {
	d.date = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func (d *Date) Convert(date string) error {
	// 02-12-2021
	if matched, _ := regexp.MatchString(`^\d{1,2}-\d{1,2}-\d{4}$`, date); matched {
		// Correct selected month
		splitText := strings.Split(date, "-")
		dateNumber := util.AtoI(splitText[0], -1)
		monthNumber := util.AtoI(splitText[1], -1)
		yearNumber := util.AtoI(splitText[2], -1)
		if dateNumber != -1 && monthNumber != -1 && yearNumber != -1 {
			d.date = time.Date(yearNumber, time.Month(monthNumber), dateNumber, 0, 0, 0, 0, time.Local)
			return nil
		}
	}
	return errors.New("invalid date format")
}

func (d *Date) NextDay() time.Time {
	return d.date.AddDate(0, 0, 1)
}

func (d *Date) PreviousDay() time.Time {
	return d.date.AddDate(0, 0, -1)
}
