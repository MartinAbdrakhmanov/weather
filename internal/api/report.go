package api

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday/v2"
)

func GetReport(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join("data", "Report.md")

	report, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
	}

	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/report.html"))
	reportb := string(blackfriday.Run(report))
	data := struct {
		Title      string
		ReportHTML string
	}{
		Title:      "Report",
		ReportHTML: reportb,
	}
	err = tmpl.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		log.Fatal(err)
	}

}
