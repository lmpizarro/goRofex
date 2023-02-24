package lib

import (
	"fmt"

	"net/http"
	"os"

	"encoding/csv"
)

type credentials struct {
	User     string
	Password string
	Account  string
}

func Read_credentials(file string) credentials {

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
		User:     line[0],
		Password: line[1],
		Account:  line[2],
	}
	return emp
}

func Token(cred *credentials) string {
	r, err := http.NewRequest("POST", Url_Auth, nil)

	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Password", cred.Password)
	r.Header.Add("X-Username", cred.User)
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

func Login() string {

    cred := Read_credentials("./env.csv")
	return  Token(&cred)


}
