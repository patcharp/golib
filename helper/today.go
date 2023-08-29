package helper

import (
	"errors"
	"github.com/patcharp/golib/v2/util"
	"regexp"
	"strings"
	"time"
)

type Today struct {
	date time.Time
}

func NewToday() Today {
	return NewTodayWithTime(time.Now())
}

func NewTodayWithTime(t time.Time) Today {
	today := Today{}
	today.Set(t)
	return today
}

func (t *Today) ToDateTime() time.Time {
	return t.date
}

func (t *Today) Set(d time.Time) {
	t.date = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

func (t *Today) Convert(date string) error {
	// 02-12-2021
	if matched, _ := regexp.MatchString(`^\d{1,2}-\d{1,2}-\d{4}$`, date); matched {
		// Correct selected month
		splitText := strings.Split(date, "-")
		dateNumber := util.AtoI(splitText[0], -1)
		monthNumber := util.AtoI(splitText[1], -1)
		yearNumber := util.AtoI(splitText[2], -1)
		if dateNumber != -1 && monthNumber != -1 && yearNumber != -1 {
			t.date = time.Date(yearNumber, time.Month(monthNumber), dateNumber, 0, 0, 0, 0, time.Local)
			return nil
		}
	}
	return errors.New("invalid date format")
}

func (t *Today) NextDay() time.Time {
	return t.date.AddDate(0, 0, 1)
}

func (t *Today) PreviousDay() time.Time {
	return t.date.AddDate(0, 0, -1)
}
