package api

import (
	"errors"
	"net/http"
	"time"
)

const (
	dateFormat = "20060102"
)

var (
	errMethodNotAllowed       = errors.New("method not allowed")
	errMissingQueryParameters = errors.New("missing query parameters")
	errInvalidNowFormat       = errors.New("invalid now format")
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, errMissingQueryParameters.Error(), http.StatusBadRequest)
		return
	}

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, errInvalidNowFormat.Error(), http.StatusBadRequest)
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
