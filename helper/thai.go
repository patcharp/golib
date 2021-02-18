package helper

import (
	"fmt"
	"time"
)

var ThaiWeekDay = map[string]interface{}{
	"Sunday":    "อาทิตย์",
	"Monday":    "จันทร์",
	"Tuesday":   "อังคาร",
	"Wednesday": "พุธ",
	"Thursday":  "พฤหัสบดี",
	"Friday":    "ศุกร์",
	"Saturday":  "เสาร์",
}

var ThaiMonth = map[string]interface{}{
	"January":   "มกราคม",
	"February":  "กุมภาพันธ์",
	"March":     "มีนาคม",
	"April":     "เมษายน",
	"May":       "พฤษภาคม",
	"June":      "มิถุนายน",
	"July":      "กรกฎาคม",
	"August":    "สิงหาคม",
	"September": "กันยายน",
	"October":   "ตุลาคม",
	"November":  "พฤศจิกายน",
	"December":  "ธันวาคม",
}

func ConvertTimeToThaiFullDateTimeFormat(t time.Time) string {
	return fmt.Sprintf("%s %s",
		ConvertTimeToThaiFullDateFormat(t),
		ConvertTimeToThaiFullTimeFormat(t),
	)
}

func ConvertTimeToThaiFullDateFormat(t time.Time) string {
	return fmt.Sprintf("วัน%sที่ %02d %s พ.ศ. %d",
		ThaiWeekDay[t.Weekday().String()],
		t.Day(),
		ThaiMonth[t.Month().String()],
		t.Year()+543,
	)
}

func ConvertTimeToThaiFullTimeFormat(t time.Time) string {
	return fmt.Sprintf("เวลา %02d:%02dน.",
		t.Hour(),
		t.Minute(),
	)
}
