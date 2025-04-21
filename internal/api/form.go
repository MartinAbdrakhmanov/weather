package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"weather/internal/model"
)

func GetForm(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.New("layout.html").ParseFiles("templates/layout.html", "templates/form.html"))

	err := tmpl.Execute(w, "layout.html")
	if err != nil {
		log.Fatal(err)
	}

}

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	data := model.FormData{
		Surname:     r.FormValue("surname"),
		Name:        r.FormValue("name"),
		Patronymic:  r.FormValue("patronymic"),
		Approval:    r.FormValue("approval"),
		Suggestions: r.FormValue("suggestions"),
		Email:       r.FormValue("email"),
	}

	fmt.Fprintf(w, "Полученные данные:\n%+v", data)

}
