package main

import (
	"encoding/base64"
	"fmt"
)

func main() {

	//\x00smbiz.temp@gmail.com\x00Temp@2020
	resp := []byte("\x00" + "smbiz.temp@gmail.com" + "\x00" + "Temp@2020")
	fmt.Println(string(resp), resp)

	sEnc := base64.StdEncoding.EncodeToString([]byte(resp))
	fmt.Println(sEnc)

}
