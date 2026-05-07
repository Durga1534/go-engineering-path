package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"unit-converter/internal/converter"
)

func main() {
	http.HandleFunc("/", handleHome)
	fmt.Println("Server starting at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("ui/index.html"))

	data := make(map[string]interface{})

	if r.Method == http.MethodPost {
		valStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")
		unitType := r.FormValue("type")

		val, _ := strconv.ParseFloat(valStr, 64)
		var result float64

		if unitType == "length" {
			result = converter.ConvertLength(val, from, to)
		} else if unitType == "temperature" {
			result = converter.ConvertTemp(val, from, to)
		} else {
			result = 0
		}

		data["Result"] = result
		data["Submitted"] = true
	}
	tmpl.Execute(w, data)
}
