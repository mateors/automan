package main

import (
	"fmt"
	"log"
	"mateors/mastererp/models"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mateors/mcb"
	"github.com/mateors/mtool"
)

//constant for this project
const (
	LoginSessionTable = "login_session"
	DeviceLogTable    = "device_log"
	LoginTable        = "login"
	ActivityLogTable  = "activity_log"
)

//CountByQuery ...
func CountByQuery(sql string, db *mcb.DB) (count float64) {

	pResponse := db.Query(sql)
	for _, val := range pResponse.Result {
		//fmt.Printf("%v %T\n", val, val)
		//rMap:=val.(map)
		cMap := val.(map[string]interface{})
		if c, ok := cMap["cnt"]; ok {
			count = c.(float64)
		}
	}
	//fmt.Println("CountTable>", pResponse.Result)
	return
}

//CheckMultipleConditionTrue this func is used for checking multiple conditions valid or ERROR
func CheckMultipleConditionTrue(formData url.Values, funcsMap map[string]interface{}, db *mcb.DB) string {

	var response string
	// funcs := map[string]interface{}{
	// 	"pass": passAndConfirmPassCheck,
	// 	"dup":  duplicateUserName,
	// }

	resAray := make([]string, 0)

	for key := range funcsMap {

		//fmt.Printf("%v %v type: %T\n", key, v, v)
		result, err := mtool.Call(funcsMap, key, formData, db) //result is type of []reflect.Value
		if err != nil {
			log.Println(err)
		}
		res := result[0].Interface().(string) //Converting reflect.Value to string
		//fmt.Printf("**** Response: %v, error: %v, ResType: %T, %T ***\n\n", res, err, result[0], res)
		resAray = append(resAray, res)

	}

	//fmt.Printf("\n## REPONSE ARRAY: %v\n\n", resAray)
	i, errorExist := mtool.ErrorInSlice(resAray, "ERROR")
	if errorExist == true {
		response = fmt.Sprintf("%v", resAray[i]) //ERROR EXIST in CheckError =>>
		//fmt.Println(response)
	} else {
		response = "OKAY"
	}

	return response

}

//InsertIntoActivityLog ...
func InsertIntoActivityLog(db *mcb.DB, loginID, cid, activity, OwnerTable, pkf, pkv, log, ip string) bool {

	t1 := time.Now()
	createDate := t1.Format("2006-01-02 15:04:05")

	nextDocID := CountDoc(BucketName, ActivityLogTable, db) + 1
	docID := fmt.Sprintf("%s::%v", ActivityLogTable, nextDocID)
	//loginIDint64, _ := strconv.ParseInt(loginID, 10, 64)
	dataActivity := models.ActivityLog{
		ID:           docID,
		LogID:        int64(nextDocID),
		Type:         ActivityLogTable,
		CompanyID:    cid,
		ActivityType: activity,
		TableName:    OwnerTable,
		PkField:      pkf,
		PkValue:      pkv,
		LogDetails:   log,
		IPAddress:    ip,
		LoginID:      loginID,
		CreateDate:   createDate,
		Status:       1,
	}
	pRes := db.InsertIntoBucket(docID, BucketName, dataActivity)
	fmt.Println(pRes.Status, pRes.Metrics.ResultCount)

	if pRes.Status == "success" {
		return true
	}

	return false
}

//UpdateByMap ...
func UpdateByMap(db *mcb.DB, keyVal map[string]interface{}, where string) bool {

	sql := UpdateQueryBuilder2(keyVal, where)
	fmt.Println(sql)
	fmt.Println()
	pRes := db.Query(sql)
	rows := pRes.GetRows()

	//returning map
	retMap := make(map[string]interface{}, 0)

	for _, row := range rows {

		for key := range row {
			//val := fmt.Sprintf("%v", row[key])
			val := row[key]

			if _, isOk := retMap[key]; isOk == false {
				vType := mtool.ValueType(val)

				//bydefault from database int number are parsing as float64
				//so if we find any value of float64 we need to convert it to int
				if vType == "float64" {
					intval, _ := strconv.Atoi(fmt.Sprintf("%v", val))
					retMap[key] = intval

				} else {
					retMap[key] = val
				}

			}

		}
	}

	fmt.Println("KeyVal:", keyVal)
	fmt.Println("retMap:", retMap)
	fmt.Println()

	isEqual := reflect.DeepEqual(keyVal, retMap)
	fmt.Println("isEqual >> ", isEqual)

	return isEqual

}

//UpdateQueryBuilder2 ...
func UpdateQueryBuilder2(keyVal map[string]interface{}, where string) (sql string) {

	var vstr, keysOnly string
	for key, val := range keyVal {

		xtype := mtool.ValueType(val)
		if xtype == "string" {
			vstr += fmt.Sprintf(`%s="%v"`, key, val) + ","

		} else if xtype == "int" {
			vstr += fmt.Sprintf("%s=%v", key, val) + ","

		} else if xtype == "float64" {
			vstr += fmt.Sprintf("%s=%v", key, val) + ","

		} else if xtype == "slice" {
			vstr += fmt.Sprintf("%s='%v'", key, val) + ","

		} else {
			vstr += fmt.Sprintf("%s='%v'", key, val) + ","
		}

		keysOnly += key + ","

	}
	vstr = strings.TrimRight(vstr, ",")
	keysOnly = strings.TrimRight(keysOnly, ",")

	//where := `type="login" AND META().id="login::3"`
	qs := `UPDATE %s SET %s WHERE %s RETURNING %s`
	sql = fmt.Sprintf(qs, BucketName, vstr, where, keysOnly)
	//fmt.Println(sql)
	//fields = keysOnly

	return
}

//UpdateQueryBuilder ...
func UpdateQueryBuilder(keyVal map[string]string, where string) (sql string) {

	var vstr, keysOnly string
	for key, val := range keyVal {

		vstr += fmt.Sprintf("%s='%v'", key, val) + ","
		keysOnly += key + ","

	}
	vstr = strings.TrimRight(vstr, ",")
	keysOnly = strings.TrimRight(keysOnly, ",")

	//where := `type="login" AND META().id="login::3"`
	qs := `UPDATE %s SET %s WHERE %s RETURNING %s`
	sql = fmt.Sprintf(qs, BucketName, vstr, where, keysOnly)
	//fmt.Println(sql)
	//fields = keysOnly

	return
}

//FieldByValue get field value by fieldName
func FieldByValue(table, fieldName, where string, db *mcb.DB) (fieldValue string) {

	sql := fmt.Sprintf("SELECT %v FROM `%v` WHERE %v;", fieldName, table, where)
	pRes := db.Query(sql)

	rows := pRes.GetRows()
	for _, row := range rows {

		for key := range row {
			fieldValue = fmt.Sprintf("%v", row[key])
			return
		}
	}
	return
}

//CountDoc ..
func CountDoc(bucketName, tableName string, db *mcb.DB) (count float64) {

	sql := fmt.Sprintf(`SELECT COUNT(*) as cnt FROM %s WHERE type="%s"`, bucketName, tableName)
	pResponse := db.Query(sql)
	//fmt.Println(sql)
	for _, val := range pResponse.Result {
		//fmt.Printf("%v %T\n", val, val)
		//rMap:=val.(map)
		cMap := val.(map[string]interface{})
		if c, ok := cMap["cnt"]; ok {
			count = c.(float64)
		}
	}
	//fmt.Println("CountTable>", pResponse.Result)
	return
}
