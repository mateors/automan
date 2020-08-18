package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/mateors/mcb"
)

var db *mcb.DB
var bucketName string

func init() {

	db = mcb.Connect("128.199.136.190", "mostain", "Mosta!n2020$")
	res, err := db.Ping()
	if err != nil {
		fmt.Println(res)
		os.Exit(1)
	}
	fmt.Println(res, err)

}

func main() {

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", home)
	http.HandleFunc("/why", why)

	http.ListenAndServe(":8080", nil)

}

func home(w http.ResponseWriter, r *http.Request) {

	ptmp, err := template.ParseGlob("template/*.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = ptmp.ParseFiles("page/home.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	data := struct {
		Title string
	}{
		Title: "Homepage",
	}

	ptmp.Execute(w, data)

}

func why(w http.ResponseWriter, r *http.Request) {

	//New("abase.gohtml")
	ptmp, err := template.ParseGlob("template/*.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = ptmp.ParseFiles("page/why.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	data := struct {
		Title string
	}{
		Title: "Why Automan",
	}

	//smtp.SendMail()
	ptmp.Execute(w, data)

}
