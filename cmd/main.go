package main

import (
	"fmt"
	"net/http"
	"strings"
	"os"

	"encoding/csv"

	lib "github.com/lmpizarro/gorofex/pkg/lib"

)

type credentials struct {
    User string
    Password string
    Account string
}

func read_credentials(file string) credentials {

    csvFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

    csvReader := csv.NewReader(csvFile)
	csvLines, err := csvReader.ReadAll()
    if err != nil {
        panic(err)
    }
    line := csvLines[0]
    emp := credentials{
        User: line[0],
        Password: line[1],
		Account: line[2],
    }
	return emp
}



func token(user, password string) string {
	r, err := http.NewRequest("POST", lib.Url_Auth, nil)

	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Password", password)
	r.Header.Add("X-Username", user)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}

	token := res.Header["X-Auth-Token"]
	return token[0]
}

func main() {

	cred := read_credentials("./env.csv")
	user := cred.User
	password := cred.Password

	token := token(user, password)

	e := "DLR/MAR23"
	fmt.Println(e)
	lib.Get_Market_Data(e, token)

}
