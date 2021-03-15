package helper

import (
	"fmt"
	"time"
)

var (
	ThaiWeekDay = map[string]interface{}{
		"Sunday":    "อาทิตย์",
		"Monday":    "จันทร์",
		"Tuesday":   "อังคาร",
		"Wednesday": "พุธ",
		"Thursday":  "พฤหัสบดี",
		"Friday":    "ศุกร์",
		"Saturday":  "เสาร์",
	}
	ThaiShortWeekDay = map[string]interface{}{
		"Sunday":    "อา",
		"Monday":    "จ",
		"Tuesday":   "อ",
		"Wednesday": "พ",
		"Thursday":  "พฤ",
		"Friday":    "ศ",
		"Saturday":  "ส",
	}

	ThaiMonth = map[string]interface{}{
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

	ThaiShortMonth = map[string]interface{}{
		"January":   "ม.ค.",
		"February":  "ก.พ.",
		"March":     "มี.ค.",
		"April":     "เม.ย.",
		"May":       "พ.ค.",
		"June":      "มิ.ค.",
		"July":      "ก.ค.",
		"August":    "ส.ค.",
		"September": "ก.ย.",
		"October":   "ต.ค.",
		"November":  "พ.ย.",
		"December":  "ธ.ค.",
	}
)

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

func ConvertTimeToThaiShortDateFormat(t time.Time) string {
	return fmt.Sprintf("%s, %02d %s %d",
		ThaiShortWeekDay[t.Weekday().String()],
		t.Day(),
		ThaiShortMonth[t.Month().String()],
		t.Year()+543,
	)
}

func ConvertTimeToThaiFullTimeFormat(t time.Time) string {
	return fmt.Sprintf("เวลา %02d:%02dน.",
		t.Hour(),
		t.Minute(),
	)
}

func ConvertTimeToThaiShortTimeFormat(t time.Time) string {
	return fmt.Sprintf("%02d:%02dน.",
		t.Hour(),
		t.Minute(),
	)
}
