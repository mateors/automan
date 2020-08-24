package main

import (
	"fmt"
	"mateors/mastererp/models"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mateors/mcb"
	"github.com/mateors/mtool"
	uuid "github.com/satori/go.uuid"
)

//main encryption decryption password
const (
	PassWordEncryptionDecryption = "Mos$unAriWusVruR0mBil360"
	BucketName                   = "master_erp"
)

//ChangePassword ...
func ChangePassword(r *http.Request, db *mcb.DB) bool {

	r.ParseForm()

	currentPass := r.FormValue("current_pass")
	newPass := r.FormValue("new_pass")
	confirmNewpass := r.FormValue("confirm_newpass")
	loginID := r.FormValue("login_id")

	//cid := r.FormValue("cid")
	//accessName := r.FormValue("access_name")
	//username := r.FormValue("username")
	fmt.Println(currentPass, newPass, confirmNewpass)

	//SELECT passw FROM master_erp WHERE type="login" AND login_id=1
	whereCondition := fmt.Sprintf(`type="%v" AND login_id=%v`, "login", loginID)
	dbPass := FieldByValue(BucketName, "passw", whereCondition, db)

	//check if currnent password matches with database hash
	cPassValid := mtool.HashCompare(currentPass, dbPass)
	fmt.Println(cPassValid)
	if cPassValid == true {

		if strings.Compare(newPass, confirmNewpass) == 0 {

			newHashPass := mtool.HashBcrypt(newPass)
			fmt.Println("NEWPASS:", newHashPass)
			fmt.Println("NOW UPDATE TO DATABASE")

			docID := fmt.Sprintf("login::%v", loginID)
			//loginIDint64, _ := strconv.ParseInt(loginID, 10, 64)
			// data := models.Login{
			// 	ID:         docID,
			// 	Type:       LoginTable,
			// 	CompanyID:  cid,
			// 	LoginID:    loginIDint64,
			// 	AccessID:   "",
			// 	AccessName: accessName,
			// 	UserName:   username,
			// 	CreateDate: "",
			// 	Password:   newHashPass,
			// }
			// pRes := db.UpsertIntoBucket(docID, BucketName, data)
			// fmt.Println(pRes.Status)

			keyVal := make(map[string]interface{}, 0)
			keyVal["passw"] = newHashPass
			//keyVal["status"] = 1

			qs := `type="login" AND META().id="%v"`
			where := fmt.Sprintf(qs, docID)
			isUpdated := UpdateByMap(db, keyVal, where)

			return isUpdated

		}

	}

	//fmt.Println(r.Form)
	return false

}

func checkSessionValid(sessionCode string) bool {

	//time.Sleep(200 * time.Millisecond)
	//fmt.Println("wait 200 milli second")
	qs := `SELECT count(*) as row_count FROM master_erp as login_session
	WHERE login_session.type="login_session"
	AND login_session.session_code="%s"
	AND login_session.status=0;`

	sql := fmt.Sprintf(qs, sessionCode)
	//fmt.Println("log_session >", sql)
	pRes := db.Query(sql)
	rows := pRes.GetRows()

	//fmt.Println(rows, len(rows), pRes.Status, "RES_COUNT:>", pRes.Metrics.ResultCount)
	for _, row := range rows {

		//fmt.Println(">", row)
		if val, isOk := row["row_count"]; isOk {

			rCount := val.(float64)
			//fmt.Println("row_count >>", ValueType(val))
			if rCount == 1 {
				return true
			}
		}

	}

	fmt.Println("PROBLEM IN checkSessionValid false")
	return false
}

//CheckIfAlreadyLoggedIn check database using session id from cookie, count==1 meaning user loggedin properly
func CheckIfAlreadyLoggedIn(r *http.Request, db *mcb.DB) (isLoggedIn bool, loginType string, loginData map[string]string) {

	loginData = make(map[string]string, 0)
	cookie, err := r.Cookie("login_session")
	if err == nil {

		sessionHexcode := cookie.Value //hexCode value
		plainText := mtool.DecodeStr(sessionHexcode, PassWordEncryptionDecryption)

		loginData = mtool.StringToMap(plainText)
		//fmt.Println("CheckIfAlreadyLoggedIn>", loginData, len(loginData))

		var sessionCode string
		if val, isExist := loginData["sid"]; isExist {
			sessionCode = val
		}

		if val, isExist := loginData["access_name"]; isExist {
			loginType = strings.ToUpper(val)
		}

		//fmt.Println("sessionCode:", sessionCode)
		isLoggedIn = checkSessionValid(sessionCode)
	}

	return
}

//Login ...
func Login(r *http.Request, w http.ResponseWriter, db *mcb.DB) {

	//check username and password are valid
	isValid, userData := LoginAuth(r, db)
	//fmt.Println(isValid, userData, len(userData))

	userType := userData["access_name"]

	battery := r.FormValue("battery")
	geolocation := r.FormValue("geolocation")
	IPAddress := r.FormValue("ip")
	username := r.FormValue("username")
	screenSize := r.FormValue("screen_size")

	t1 := time.Now()
	createDate := t1.Format("2006-01-02 15:04:05")

	cookie, err := r.Cookie("login_session")
	if err == http.ErrNoCookie {
		fmt.Println("NO SESSION FOUND SEEMS OKAY...")
	}
	if err != http.ErrNoCookie {
		//session found so destroying them
		cookie.Name = "login_session"
		cookie.MaxAge = -1
		cookie.Value = ""
		cookie.Path = "/"
		cookie.HttpOnly = true
		http.SetCookie(w, cookie)
		fmt.Println("login_session SESSION FOUND SO DESTROYING...")
	}

	lcookie, err := r.Cookie("login_error")
	if err != http.ErrNoCookie {
		lcookie := &http.Cookie{Name: "login_error", MaxAge: -1, Value: "", Path: "/login", HttpOnly: true}
		http.SetCookie(w, lcookie)
		fmt.Println("login_error SESSION FOUND SO DESTROYING...")
	}

	//set cookie
	if isValid {

		//create session
		sID, _ := uuid.NewV4() //session code
		sessionCode := sID.String()
		userData["sid"] = sessionCode
		delete(userData, "passw") //removing password from map
		//fmt.Println()
		//fmt.Println(userData)

		sessionStr := mtool.MapToString(userData)
		//fmt.Println()
		//fmt.Println("sessionValue:", sessionStr)
		sessionValue := mtool.EncodeStr(sessionStr, PassWordEncryptionDecryption)

		c := &http.Cookie{
			Name:   "login_session",
			Value:  sessionValue,
			MaxAge: 86400,
		}

		http.SetCookie(w, c)
		//insert into login_session table for future reference

		bInfo := mtool.BrowserInfo(r.UserAgent(), battery)

		loginID := fmt.Sprintf("login::%v", userData["login_id"])
		cid := fmt.Sprintf("company::%v", userData["cid"]) // "cid": "company::1",
		ip := mtool.ReadUserIP(r)

		//cid := userData["cid"].(float64)
		//loginID := userData["login_id"].(float64) //??

		nextDeviceLog := CountDoc(BucketName, DeviceLogTable, db) + 1
		deviceLog := fmt.Sprintf("device_log::%v", nextDeviceLog)

		deviceData := models.DeviceLog{
			ID:          deviceLog,
			Type:        DeviceLogTable,
			LoginID:     loginID,
			DeviceType:  bInfo["device"],
			CompanyID:   cid,
			Browser:     bInfo["browser_version"],
			Os:          bInfo["os_version"],
			Platform:    bInfo["platform"],
			ScreenSize:  screenSize,
			GeoLocation: geolocation,
			CreateDate:  createDate,
			Status:      1,
		}
		db.InsertIntoBucket(deviceLog, BucketName, deviceData)

		//fmt.Println(">>", cid, loginID, deviceLog)
		//loginIDint64, _ := strconv.ParseInt(loginID, 10, 64)

		//update old login_session status 0 to 1
		//UPDATE master_erp SET status=1 WHERE type="login_session" AND status=0 AND login_id="login::2" RETURNING *;
		sql := fmt.Sprintf(`UPDATE master_erp SET status=1 WHERE type="login_session" AND status=0 AND login_id="%s";`, loginID)
		pRes := db.Query(sql)
		fmt.Println("OldSessionUpdate:", pRes.Status, pRes.Metrics.ResultCount)

		nextDocID := CountDoc(BucketName, LoginSessionTable, db) + 1
		docID := fmt.Sprintf("login_session::%v", nextDocID)
		data := models.LoginSession{
			ID:          docID,
			Type:        LoginSessionTable,
			CompanyID:   cid,
			DeviceInfo:  deviceLog,
			SessionCode: sessionCode,
			LoginID:     loginID,
			IPAddress:   ip,
			LoginTime:   createDate,
			CreateDate:  createDate,
			Status:      0,
		}

		pRes = db.InsertIntoBucket(docID, BucketName, data)
		fmt.Println(pRes.Status, pRes.Metrics.ResultCount)

		//insert into activity log
		// nextDocID = CountDoc(BucketName, ActivityLogTable, db) + 1
		// docID = fmt.Sprintf("%s::%v", ActivityLogTable, nextDocID)
		// loginIDint64, _ := strconv.ParseInt(userData["login_id"], 10, 64)
		// dataActivity := models.ActivityLog{
		// 	ID:           docID,
		// 	LogID:        int64(nextDocID),
		// 	Type:         ActivityLogTable,
		// 	CompanyID:    cid,
		// 	ActivityType: "Login",
		// 	TableName:    LoginTable,
		// 	PkField:      "aid",
		// 	PkValue:      loginID,
		// 	LogDetails:   "",
		// 	IPAddress:    IPAddress,
		// 	LoginID:      loginIDint64,
		// 	CreateDate:   createDate,
		// 	Status:       1,
		// }
		// db.InsertIntoBucket(docID, BucketName, dataActivity)
		var activity, OwnerTable, pkf, pkv string = "Login", "login", "login_id", loginID
		var log string = fmt.Sprintf("%v %s logged in using %s %s %s", userType, username, bInfo["device"], screenSize, bInfo["browser_version"])
		InsertIntoActivityLog(db, loginID, cid, activity, OwnerTable, pkf, pkv, log, IPAddress)

		time.Sleep(200 * time.Millisecond)
		//fmt.Println("wait 200 milli second")

		//redirect to dashboard
		//fmt.Println("redirecting to /dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return

	}

	//login failed
	lcookie = &http.Cookie{Name: "login_error", Value: "invalid login", Path: "/login", HttpOnly: true}
	http.SetCookie(w, lcookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return

}

func duplicateUserName(req url.Values, db *mcb.DB) string {

	response := fmt.Sprintf("ERROR user %v already exist", req.Get("email"))

	length := len(req)
	//fmt.Printf("duplicateUserName FORM_DATA: %v, LENGTH: %v\n", req, length)

	if length > 2 {

		email := req.Get("email")
		bucketName := req.Get("bucket")
		//whereClause := fmt.Sprintf("username='%v'", email)
		sql := fmt.Sprintf("SELECT count(*)as cnt FROM `%s` WHERE type='%s' AND username='%s';", bucketName, "login", email)
		//fmt.Println(sql)
		count := CountByQuery(sql, db)
		//count := dbf.CheckCount("login", whereClause, db)
		//fmt.Printf("USER_COUNT: %v\n", count)

		if count == 0 && email != "" {
			response = "valid"
		}
	}

	return response
}

func passAndConfirmPassCheck(req url.Values, db *mcb.DB) string {

	response := "ERROR password and confirm password does not match."
	//length := len(req)
	//fmt.Printf("passAndConfirmPassCheck FORM_DATA: %v, LENGTH: %v\n", req, length)

	password := req.Get("password")
	confirmPassword := req.Get("confirm_password")

	if password == confirmPassword {
		response = "valid"
	}
	return response
}

//CreateAccount ...
func CreateAccount(r *http.Request, db *mcb.DB) (bool, string) {

	r.ParseForm()
	form := r.Form

	funcsMap := map[string]interface{}{
		//"pass": passAndConfirmPassCheck,
		"dup": duplicateUserName,
	}

	response := CheckMultipleConditionTrue(form, funcsMap, db)
	//fmt.Println("RESPONSE >", response)

	//if strings.Compare(response, "OKAY")
	if response == "OKAY" {

		//form url.Values
		//map[cid:[] email:[bill.rassel@gmail.com] full_name:[MOSTAIN BILLAH] mobile:[01672710028] password:[1234]]
		//fmt.Println("CreateAccount>>>>>", form)
		bucketName := form.Get("bucket")
		mobile := form.Get("mobile")
		email := form.Get("email")
		password := form.Get("password")
		cid := form.Get("cid")
		company := fmt.Sprintf("company::%s", cid)

		hashPass := mtool.HashBcrypt(password)

		//db.Insert(form)
		//sql := fmt.Sprintf(`SELECT COUNT(*) as cnt FROM %s WHERE type="%s"`, "master_erp", "login")
		count := CountDoc(bucketName, "login", db)
		nextLoginID := count + 1
		loginAid := fmt.Sprintf("login::%v", nextLoginID)
		loginID := fmt.Sprintf("%v", nextLoginID)

		count = CountDoc(bucketName, "account", db) + 1
		accountID := fmt.Sprintf("%v", count)

		t1 := time.Now()
		createDate := t1.Format("2006-01-02 15:04:05")

		var acc models.Account
		var login models.Login

		//login table
		form.Set("cid", company)
		form.Set("aid", loginAid)
		form.Set("type", "login")
		form.Set("username", email)
		form.Set("passwrd", hashPass)
		form.Set("access_name", "student")
		form.Set("login_id", loginID)
		form.Set("create_date", createDate)
		form.Set("status", "1")
		//db.ProcessData()

		// var contact []models.Contact
		// contact = append(contact, models.Contact{ID: "contact::1", Type: "contact"})
		// acc.ContactInfo = contact

		//json.Unmarshal()

		//login table created
		db.Insert(form, &login)
		//resMessage := db.Insert(form, &login)
		//fmt.Println("Login::", resMessage)
		//fmt.Println()

		//account table
		accountAid := fmt.Sprintf("account::%v", accountID)
		form.Set("aid", accountAid)
		form.Set("login_id", loginAid)
		form.Set("account_id", accountID)
		form.Set("account_type", "student")
		form.Set("account_name", form.Get("full_name"))
		form.Set("type", "account")
		form.Set("mobile", mobile)
		form.Set("email", email)

		db.Insert(form, &acc)
		//resMessage = db.Insert(form, &acc)
		//fmt.Println(form)
		//fmt.Println("Account::", resMessage)

		return true, "User registration successful."

	}

	return false, response
	// for key := range form {
	// 	form.Del(key)
	// }
	// fmt.Println("FORM:", form)
}

//LoginAuth check against couchbd
func LoginAuth(r *http.Request, db *mcb.DB) (bool, map[string]string) {

	r.ParseForm()
	//form := r.Form
	username := r.FormValue("username")
	password := r.FormValue("password")

	//sql := fmt.Sprintf("SELECT * FROM master_erp WHERE username='%v' AND type='login'", username)
	qs := `SELECT company.cid,
	company.company_name,
	company.website,
	account.account_id,
	account.account_name,
	account.mobile,
	account.email,
	login.login_id,
	login.access_name,
	login.passw FROM master_erp AS login
	LEFT JOIN master_erp as company ON login.cid=META(company).id
	LEFT JOIN master_erp as account ON account.login_id=META(login).id AND account.type='account'
	WHERE login.username='%v' AND login.type='login';`

	sql := fmt.Sprintf(qs, username)
	//fmt.Println(sql)
	pres := db.Query(sql)

	rows := pres.GetRows()
	var dbHashPassword string
	userData := make(map[string]string, 0)

	for _, row := range rows {
		dbHashPassword = row["passw"].(string)
		//userData = row
		for k := range row {
			if val, isOk := row[k]; isOk {
				userData[k] = fmt.Sprintf("%v", val)
			}
		}
	}

	isTrue := mtool.HashCompare(password, dbHashPassword)

	return isTrue, userData
}
