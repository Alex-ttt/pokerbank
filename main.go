package main

import (
	"fmt"
	"html/template"
	"net/http"
	"pokerscore/models"
	. "strconv"
	"time"
)

type Welcome struct {
	Name string
	Time string
}

func main() {
	models.InitDb("./poker.sqlite")

	rows, err := models.Db.Query("SELECT Id, Value FROM items")
	checkErr(err)

	var id, value int
	var result string
	for rows.Next() {
		err = rows.Scan(&id, &value)
		checkErr(err)
	}
	result = Itoa(id) + " " + Itoa(value)
	fmt.Printf("%s", result)
	welcome := Welcome{"Anonymous" + " " + result, time.Now().Format(time.Stamp)}

	templates := template.Must(template.ParseFiles("templates/welcome-template.html"))

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.
	http.HandleFunc("/" , func(w http.ResponseWriter, r *http.Request) {

		//Takes the name from the URL query e.g ?name=Martin, will set welcome.Name = Martin.
		if name := r.FormValue("name"); name != "" {
			welcome.Name = name;
		}
		//If errors show an internal server error message
		//I also pass the welcome struct to the welcome-templates.html file.
		if err := templates.ExecuteTemplate(w, "welcome-template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening");
	fmt.Println(http.ListenAndServe(":8080", nil));
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
