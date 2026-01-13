package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RenderQuickAddRow(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/quick_add_row.html"))
	t.Execute(w, nil)
}

func AddRow(w http.ResponseWriter, r *http.Request) {
	var company, role, status, notes, jobURL string
	company = r.FormValue("company")
	role = r.FormValue("role")
	status = r.FormValue("status")
	notes = r.FormValue("notes")
	jobURL = r.FormValue("job_url")

	if company == "" || role == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `<tr id="quick-add-row-error">
        <td colspan="6" style="color:#f87171;padding:8px 12px;">
            Please fill in all required fields.
        </td>
    </tr>`)
		return
	}

	_, err := DB.Exec(`
		INSERT INTO applications (company, role, status, notes, job_url, created_at)
		VALUES (?, ?, ?, ?, ?, datetime('now'));
	`, company, role, status, notes, jobURL)
	if err != nil {
		http.Error(w, "Failed to add row to database", http.StatusInternalServerError)
		return
	}

	var id int64
	err = DB.QueryRow("SELECT last_insert_rowid();").Scan(&id)
	if err != nil {
		http.Error(w, "Failed to retrieve new row ID", http.StatusInternalServerError)
		return
	}

	t := template.Must(template.ParseFiles("templates/application_row.html"))
	t.Execute(w, map[string]any{
		"ID":        id,
		"Company":   company,
		"Role":      role,
		"Status":    status,
		"Notes":     notes,
		"JobUrl":    jobURL,
		"CreatedAt": "Just now",
	})
}

func UpdateRow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var company, role, status, notes, jobURL string
	company = r.FormValue("company")
	role = r.FormValue("role")
	status = r.FormValue("status")
	notes = r.FormValue("notes")
	jobURL = r.FormValue("job_url")
	date := r.FormValue("date")

	if company == "" || role == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `<tr id="edit-row-error">
		<td colspan="6" style="color:#f87171;padding:8px 12px;">
			Please fill in all required fields.
		</td>
	</tr>`)
		return
	}

	_, err := DB.Exec(`
		UPDATE applications
		SET company = ?, role = ?, status = ?, notes = ?, job_url = ?
		WHERE id = ?;
	`, company, role, status, notes, jobURL, id)
	if err != nil {
		http.Error(w, "Failed to update row in database", http.StatusInternalServerError)
		return
	}

	if date != "" {
		_, err = DB.Exec(`
			UPDATE applications
			SET created_at = ?
			WHERE id = ?;
		`, date, id)
		if err != nil {
			http.Error(w, "Failed to update date in database", http.StatusInternalServerError)
			return
		}
	} else {
		err = DB.QueryRow(`
		SELECT strftime('%Y-%m-%d', created_at)
		FROM applications
		WHERE id = ?;
	`, id).Scan(&date)
		if err != nil {
			http.Error(w, "Failed to retrieve updated row from database", http.StatusInternalServerError)
			return
		}
	}

	t := template.Must(template.ParseFiles("templates/application_row.html"))
	t.Execute(w, map[string]any{
		"ID":        id,
		"Company":   company,
		"Role":      role,
		"Status":    status,
		"Notes":     notes,
		"JobUrl":    jobURL,
		"CreatedAt": date,
	})
}
