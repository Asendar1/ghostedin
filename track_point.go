package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func GetRows(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query("SELECT company, strftime('%Y-%m-%d', applied_date), status, role FROM applications ORDER BY created_at DESC;")
    if err != nil {
        http.Error(w, "Database query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        t := template.Must(template.ParseFiles("templates/application_row.html"))
        var company, applied_date, status, role string
        rows.Scan(&company, &applied_date, &status, &role)
        t.Execute(w, map[string]any{
            "Company":    company,
            "AppliedDate": applied_date,
            "Status":     status,
            "Role":       role,
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
            tech_stack TEXT,
            notes TEXT,
            job_url TEXT,
            resume_version TEXT,
            last_followup DATE,
            response_date DATE,
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
