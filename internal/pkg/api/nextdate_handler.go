package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type monthDaysRule struct {
	day     [32]bool
	last    bool
	preLast bool
}

func matchDays(day time.Time, rule monthDaysRule) bool {
	lastDayOfMonth := time.Date(day.Year(), day.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	if rule.day[day.Day()] {
		return true
	}

	if rule.last && day.Day() == lastDayOfMonth {
		return true
	}

	if rule.preLast && day.Day() == lastDayOfMonth-1 {
		return true
	}

	return false
}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeatElements := strings.Fields(repeat)
	if len(repeatElements) == 0 {
		return "", errIncorrectRepeat
	}

	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errParseStartDate
	}

	fixedStartDate := startDate

	switch repeatElements[0] {
	case "d":
		if len(repeatElements) != 2 {
			return "", errIncorrectRepeat
		}

		days, err := strconv.Atoi(repeatElements[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errIncorrectRepeat
		}

		for !startDate.After(now) || !startDate.After(fixedStartDate) {
			startDate = startDate.AddDate(0, 0, days)
		}

	case "y":
		if len(repeatElements) != 1 {
			return "", errIncorrectRepeat
		}

		for !startDate.After(now) || !startDate.After(fixedStartDate) {
			startDate = startDate.AddDate(1, 0, 0)
		}

	case "w":
		if len(repeatElements) != 2 {
			return "", errIncorrectRepeat
		}

		daysOfWeekStr := strings.Split(repeatElements[1], ",")

		daysOfWeek := make([]int, 0)
		for _, dayStr := range daysOfWeekStr {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day > 7 || day < 1 {
				return "", errIncorrectRepeat
			}

			daysOfWeek = append(daysOfWeek, day)
		}

		var weekDay [8]bool
		for _, day := range daysOfWeek {
			weekDay[day] = true
		}

		for !startDate.After(now) || !startDate.After(fixedStartDate) {
			startDate = startDate.AddDate(0, 0, 1)
		}

		var dateFound bool
		for i := 0; i < 5000; i++ {
			wd := int(startDate.Weekday())

			if wd == 0 {
				wd = 7
			}

			if weekDay[wd] {
				dateFound = true
				break
			}

			startDate = startDate.AddDate(0, 0, 1)
		}

		if !dateFound {
			return "", errNextDateNotFound
		}

	case "m":
		if len(repeatElements) != 2 && len(repeatElements) != 3 {
			return "", errIncorrectRepeat
		}

		days := make([]int, 0)
		daysStr := strings.Split(repeatElements[1], ",")
		for _, dayStr := range daysStr {
			dayInt, err := strconv.Atoi(dayStr)
			if err != nil || dayInt > 31 || dayInt < -2 || dayInt == 0 {
				return "", errIncorrectRepeat
			}

			days = append(days, dayInt)
		}

		months := make([]int, 0)
		if len(repeatElements) == 3 {
			monthsStr := strings.Split(repeatElements[2], ",")
			for _, monthStr := range monthsStr {
				monthInt, err := strconv.Atoi(monthStr)
				if err != nil || monthInt > 12 || monthInt < 1 || monthInt == 0 {
					return "", errIncorrectRepeat
				}

				months = append(months, monthInt)
			}

		}

		var daysRule monthDaysRule
		for _, d := range days {
			switch d {
			case -1:
				daysRule.last = true
			case -2:
				daysRule.preLast = true
			default:
				daysRule.day[d] = true
			}
		}

		var monthBool [13]bool
		if len(months) == 0 {
			for i := 1; i <= 12; i++ {
				monthBool[i] = true
			}
		} else {
			for _, m := range months {
				monthBool[m] = true
			}
		}

		for !startDate.After(now) || !startDate.After(fixedStartDate) {
			startDate = startDate.AddDate(0, 0, 1)
		}

		var dateFound bool
		for i := 0; i < 5000; i++ {
			if monthBool[int(startDate.Month())] && matchDays(startDate, daysRule) {
				dateFound = true
				break
			}

			startDate = startDate.AddDate(0, 0, 1)
		}

		if !dateFound {
			return "", errNextDateNotFound
		}

	default:
		return "", errIncorrectRepeat
	}

	return startDate.Format(dateFormat), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, errMissingQueryParameters.Error(), http.StatusBadRequest)
		return
	}

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, errIncorrectNowFormat.Error(), http.StatusBadRequest)
		return
	}

	next, err := nextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
