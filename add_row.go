package main

import (
	"html/template"
	"net/http"
)

func RenderQuickAddRow(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/quick_add_row.html"))
	t.Execute(w, nil)
}
