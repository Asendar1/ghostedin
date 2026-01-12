package main

import (
	"database/sql"
	"time"
)

func GetTotalApplications(DB *sql.DB) int {
    var total int
    err := DB.QueryRow("SELECT COUNT(*) FROM applications;").Scan(&total)
	if err != nil {
		return -1
	}
	return total
}

func GetDailyApplication(DB *sql.DB) int {
	var daily int
	today := time.Now().Format("2006-01-02")
	err := DB.QueryRow("SELECT COUNT(*) FROM applications WHERE applied_date = ?;", today).Scan(&daily)
	if err != nil {
		return -1
	}
	return daily
}

func GetInterviewPercentage(DB *sql.DB) float64 {
	var interviewsRate float64
	err := DB.QueryRow("SELECT COUNT(*) FROM applications WHERE status = 'Interviewing';").Scan(&interviewsRate)
	if err != nil {
		return -1
	}
	interviewsRate /= float64(GetTotalApplications(DB)) * 100
	return interviewsRate
}
