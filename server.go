package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

//json for phones
type product struct {
	ID         int    `json:"id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	Os         string `json:"os"`
	Image      string `json:"image"`
	Screensize int    `json:"screensize"`
}

//array store rows
type products []product

//get phones
func get(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM phones")
	var phones products
	for rows.Next() {
		var phone product
		rows.Scan(&phone.ID, &phone.Brand, &phone.Model, &phone.Os, &phone.Image, &phone.Screensize)
		phones = append(phones, phone)
	}
	data, err := json.Marshal(phones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//post form
func post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	os := r.FormValue("os")
	image := r.FormValue("image")
	screensize := r.FormValue("screensize")

	insert, err := database.Prepare(`INSERT INTO phones (brand, model, os, image, screensize) VALUES (?, ?, ?, ?, ?)`)
	insert.Exec(brand, model, os, image, screensize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "brand = %s\n", brand)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
}

//update phone
func put(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	os := r.FormValue("os")
	image := r.FormValue("image")
	screensize := r.FormValue("screensize")
	id := r.FormValue("id")

	update, err := database.Prepare(`UPDATE phones SET brand = ?,model=?, os=?, image=?, screensize=? WHERE id=?`)
	update.Exec(brand, model, os, image, screensize, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "information updated! Please go back and refresh!id = %s\n", id)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
}

//reset database
func reset(w http.ResponseWriter, r *http.Request) {
	database.Exec("DELETE FROM phones")
	w.WriteHeader(200)
}

//delete infomation
func delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	stmt, err := database.Prepare("DELETE FROM phones WHERE id = ?")
	stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "information deleted!Please go back and refresh!id = %s\n", id)
	w.WriteHeader(200)
}

func main() {
	database, _ = sql.Open("sqlite3", "./phones.db")
	database.Exec("CREATE TABLE IF NOT EXISTS phones (id 	INTEGER PRIMARY KEY, brand	CHAR(100) NOT NULL, model 	CHAR(100) NOT NULL, os 	CHAR(10) NOT NULL, image 	CHAR(254) NOT NULL, screensize INTEGER NOT NULL)")
	database.Exec(`INSERT INTO phones (brand, model, os, image, screensize) VALUES ('apple','8p','nowhere','fly',9)`)

	http.Handle("/", http.FileServer(http.Dir("./file")))
	http.HandleFunc("/phones", get)
	http.HandleFunc("/post", post)
	http.HandleFunc("/update", put)
	http.HandleFunc("/reset", reset)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":8080", nil)
}
