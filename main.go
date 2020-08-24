package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mateors/mcb"
	"github.com/mateors/mtool"
)

var db *mcb.DB
var workingDirectory string
var websiteTemplateAbsPath string
var adminTemplateAbsPath string

func init() {

	workingDirectory, _ = os.Getwd()
	//fmt.Println("workingDirectory:", workingDirectory)
	adminTemplateAbsPath = filepath.Join(workingDirectory, "templates", "admin", "*.gohtml")
	//fmt.Println("adminTemplateAbsPath:", adminTemplateAbsPath)

	websiteTemplateAbsPath = filepath.Join(workingDirectory, "templates", "website", "*.gohtml")
	//fmt.Println("websiteTemplateAbsPath:", websiteTemplateAbsPath)

	db = mcb.Connect("128.199.136.190", "mostain", "Mosta!n2020$")
	res, err := db.Ping()
	if err != nil {
		fmt.Println(res)
		os.Exit(1)
	}
	fmt.Println(res, err)

}

func main() {

	//smux.PathPrefix("/vdata/").Handler(http.StripPrefix("/vdata/", http.FileServer(http.Dir(data_path))))

	dataPath := filepath.Join(workingDirectory, "data")
	http.Handle("/vdata/", http.StripPrefix("/vdata/", http.FileServer(http.Dir(dataPath))))

	assetPath := filepath.Join(workingDirectory, "assets")
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(assetPath))))

	http.HandleFunc("/", home)
	http.HandleFunc("/why", why)
	http.HandleFunc("/login", login)
	http.HandleFunc("/forgot_password", forgetPassword)
	http.HandleFunc("/changepass", changepass)

	http.HandleFunc("/dashboard", dashboard)

	// retMap := make(map[string]interface{}, 0)

	// retMap["status"] = 10
	// retMap["name"] = "Mostain"
	// retMap["marks"] = 50.89
	// retMap["ages"] = []string{"10", "20"}

	// for k, v := range retMap {

	// 	xType := ValueType(v)
	// 	fmt.Println(k, "=", xType)
	// }

	// fmt.Println(retMap)

	http.ListenAndServe(":8080", nil)

	/*
		fmt.Println()
		d := "access_name:student,cid:1,company_name:AUTOMAN LTD.,login_id:2,sid:7eafb6b5-d116-4362-82b8-ff8b45252135,account_id:2,account_name:MOSTAIN BILLAH,email:bill.rassel@gmail.com,mobile:01672710028"
		var key string = "Mos$unAriWusVruR0mBil360"

		hcode := EncodeStr(d, key)
		fmt.Println(hcode, len(hcode))

		fmt.Println()
		plaintext := DecodeStr(hcode, key)
		fmt.Println("Decrypted:", plaintext, len(plaintext))

		fmt.Println()
		sMap := stringToMap(d)
		fmt.Println(sMap, len(sMap))
		fmt.Println(sMap["sid"])
	*/

}

func home(w http.ResponseWriter, r *http.Request) {

	ptmp, err := template.ParseGlob(websiteTemplateAbsPath)
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
	ptmp, err := template.ParseGlob(websiteTemplateAbsPath)
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

func login(w http.ResponseWriter, r *http.Request) {

	ptmp, err := template.ParseFiles("page/login.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	r.ParseForm()
	fmt.Println("Method:", r.Method)

	if r.Method == "POST" {

		//fmt.Println("form received", r.Form)
		//method 1 to catch html form value
		username := r.FormValue("username")
		password := r.FormValue("password")
		geolocation := r.FormValue("geolocation")
		ip := r.FormValue("ip")
		platform := r.FormValue("platform")
		screenSize := r.FormValue("screen_size")
		battery := r.FormValue("battery")

		fmt.Println("username:", username)
		fmt.Println("password:", password)
		fmt.Println("geolocation:", geolocation)
		fmt.Println("ip:", ip)
		fmt.Println("platform:", platform)
		fmt.Println("screen_size:", screenSize)
		fmt.Println("Battery:", battery)

		//fmt.Println()
		fmt.Println()

		Login(r, w, db)

		//Mthod 2 /Advance way to catch html form value map[string][]string
		// for key, valA := range r.Form {
		// 	fmt.Println(key, valA[0])
		// }

	}

	data := struct {
		Title      string
		LoginError string
	}{
		Title:      "Login",
		LoginError: "",
	}

	//smtp.SendMail()
	ptmp.Execute(w, data)

}

func forgetPassword(w http.ResponseWriter, r *http.Request) {

	ptmp, err := template.ParseFiles("page/forget_password.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}

	data := struct {
		Title      string
		LoginError string
	}{
		Title:      "Login",
		LoginError: "",
	}

	//smtp.SendMail()
	ptmp.Execute(w, data)

}

func dashboard(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("dashboard panicking with value >>", r)
		}
	}()

	//check if user logged in or not
	isLoggedIn, loginType, logData := CheckIfAlreadyLoggedIn(r, db)
	fmt.Println("/dashboard ....")
	fmt.Println(isLoggedIn, loginType)

	if isLoggedIn == false {
		fmt.Println("No session found, redirecting /dashboard to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var cid, accountID string

	cid = mtool.GetMapValue(logData, "cid")
	accountID = mtool.GetMapValue(logData, "account_id")
	fmt.Println(cid, accountID)

	ptmp, err := template.ParseGlob(adminTemplateAbsPath)
	if err != nil {
		fmt.Println(err.Error())
	}

	pageName := filepath.Join(workingDirectory, "page", "dashboard.gohtml")
	ptmp, err = ptmp.ParseFiles(pageName)
	if err != nil {
		fmt.Println(err.Error())
	}

	data := struct {
		Title           string
		AccountName     string
		CompanyName     string
		CompanyWebsite  string
		LoginType       string
		Page            string
		WebuserCount    int
		SystemuserCount int
	}{
		Title:           "Dashboard",
		AccountName:     mtool.GetMapValue(logData, "account_name"),
		CompanyName:     mtool.GetMapValue(logData, "company_name"),
		CompanyWebsite:  mtool.GetMapValue(logData, "website"),
		LoginType:       "SUPERADMIN",
		Page:            r.RequestURI,
		WebuserCount:    0,
		SystemuserCount: 5,
	}

	if err := ptmp.Execute(w, data); err != nil {
		log.Fatal(err)
	}

}

func changepass(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("changepass panicking with value >>", r)
		}
	}()

	isLoggedIn, loginType, logData := CheckIfAlreadyLoggedIn(r, db)
	if isLoggedIn == false {
		fmt.Println("No session found, redirecting /dashboard to /login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	//handling form submission
	if r.Method == http.MethodPost {

		data := make(map[string]string, 0)
		r.ParseForm()

		loginID := mtool.GetMapValue(logData, "login_id")
		fmt.Println(loginID)
		r.Form.Add("login_id", loginID)

		isOk := ChangePassword(r, db)
		var message string
		if isOk == true {
			message = "OK"
			//fmt.Println("Password")
		} else {
			message = "ERROR"
		}

		data["message"] = message
		//response return to client
		js, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Something wrong @ json.Marshal", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(js))

	}

	if r.Method == http.MethodGet {

		ptmp, err := template.ParseGlob(adminTemplateAbsPath)
		if err != nil {
			fmt.Println(err.Error())
		}

		pageName := filepath.Join(workingDirectory, "page", "changepass.gohtml")
		ptmp, err = ptmp.ParseFiles(pageName)
		if err != nil {
			fmt.Println(err.Error())
		}

		data := struct {
			Title          string
			LoginType      string
			Page           string
			CompanyName    string
			CompanyWebsite string
		}{
			Title:          "ChangePassword",
			LoginType:      loginType,
			Page:           r.RequestURI,
			CompanyName:    mtool.GetMapValue(logData, "company_name"),
			CompanyWebsite: mtool.GetMapValue(logData, "website"),
		}

		if err := ptmp.Execute(w, data); err != nil {
			log.Fatal(err)
		}
	}

}
