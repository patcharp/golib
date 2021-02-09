package helper

import (
	"fmt"
	"time"
)

var thaiWeekDay = map[string]interface{}{
	"Sunday":    "อาทิตย์",
	"Monday":    "จันทร์",
	"Tuesday":   "อังคาร",
	"Wednesday": "พุธ",
	"Thursday":  "พฤหัสบดี",
	"Friday":    "ศุกร์",
	"Saturday":  "เสาร์",
}

var thaiMonth = map[string]interface{}{
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

func ConvertTimeToThaiFormat(t time.Time) string {
	return fmt.Sprintf("วัน%sที่ %02d %s พ.ศ. %d เวลา %02d:%02dน.",
		thaiWeekDay[t.Weekday().String()],
		t.Day(),
		thaiMonth[t.Month().String()],
		t.Year()+543,
		t.Hour(),
		t.Minute(),
	)
}
