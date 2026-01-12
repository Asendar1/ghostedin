package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func GetStatistics(w http.ResponseWriter, r *http.Request) {
    var totalApplications, dailyApplications int
    var interviewsRate float64

    totalApplications = GetTotalApplications(DB)
    dailyApplications = GetDailyApplication(DB)
    if totalApplications > 0 {
        interviewsRate = GetInterviewPercentage(DB)
    } else {
        interviewsRate = 0
    }

    // replace with template later
    if totalApplications == -1 || dailyApplications == -1 || interviewsRate == -1 {
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }

    t := template.Must(template.ParseFiles("templates/statistics_cards.html"))
    t.Execute(w, map[string]any{
        "Total": totalApplications,
        "AppliedToday": dailyApplications,
        "InterviewRate": interviewsRate,
    })
}

func GetRows(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query(`
        SELECT
            company,
            strftime('%Y-%m-%d', created_at),
            status,
            role,
            COALESCE(notes, ''),
            COALESCE(job_url, '')
        FROM applications
        ORDER BY created_at DESC;
    `)
    if err != nil {
        log.Println("Database query error:", err)
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        t := template.Must(template.ParseFiles("templates/application_row.html"))
        var company, status, role, notes, job_url, created_at string
        err := rows.Scan(
            &company,
            &status,
            &created_at,
            &role,
            &notes,
            &job_url,
        )
        if err != nil {
             log.Println("Error scanning row:", err)
             continue
        }

        t.Execute(w, map[string]any{
            "Company":       company,
            "Status":        status,
            "Role":          role,
            "Notes":         notes,
            "JobUrl":        job_url,
            "CreatedAt":     created_at,
        })
    }
}

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "jobs.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS applications (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            company TEXT NOT NULL,
            role TEXT NOT NULL,
            applied_date DATE NOT NULL DEFAULT (date()),
            status TEXT NOT NULL DEFAULT 'Applied',
            notes TEXT,
            job_url TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadSheet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/static/track.html")
}

func CreateApplication(w http.ResponseWriter, r *http.Request) {

}
