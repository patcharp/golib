package crontab

const (
	// Minute
	SpecEveryMinute = "* * * * *"
	SpecEvery5Min   = "*/5 * * * *"
	SpecEvery10Min  = "*/10 * * * *"
	SpecEvery15Min  = "*/15 * * * *"
	SpecEvery30Min  = "*/30 * * * *"
	SpecEvery45Min  = "*/45 * * * *"

	// Hour
	SpecEveryHour   = "0 * * * *"
	SpecEvery3Hour  = "0 */3 * * *"
	SpecEvery6Hour  = "0 */6 * * *"
	SpecEvery12Hour = "0 */12 * * *"

	// Day
	SpecEveryDay       = "0 0 * * *"
	SpecEverySunday    = "0 0 * * 0"
	SpecEveryMonday    = "0 0 * * 1"
	SpecEveryTuesday   = "0 0 * * 2"
	SpecEveryWednesday = "0 0 * * 3"
	SpecEveryThursday  = "0 0 * * 4"
	SpecEveryFriday    = "0 0 * * 5"
	SpecEverySaturday  = "0 0 * * 6"
	SpecEveryWeekDay   = "0 0 * * 1-5"
	SpecEveryWeekEnd   = "0 0 * * 0,6"

	// Time
	SpecEveryMidnight = SpecEveryDay
	SpecEveryMidday   = "0 12 * * *"
)
