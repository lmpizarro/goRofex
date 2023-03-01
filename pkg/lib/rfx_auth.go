package lib

import (
	"fmt"
	"io/ioutil"

	"net/http"
	"os"

	"encoding/csv"

	"encoding/json"

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

func Cer(){
	r, err := http.NewRequest("POST", "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp", nil)

	if err != nil {
		panic(err)
	}
	form := r.URL.Query()
	form.Add("primeravez", "1")
	form.Add("fecha_desde", "20230206")
	form.Add("fecha_hasta", "20230227")
	form.Add("serie", "3540")
	form.Add("serie1", "0")
	form.Add("serie2", "0")
	form.Add("serie3", "0")
	form.Add("serie4", "0")
	form.Add("detalle", "CER (Base 2.2.2002=1)")
	r.URL.RawQuery = form.Encode()

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.StatusCode)
	}

	var resp map[string]interface{}

    json.NewDecoder(res.Body).Decode(&resp)

    b, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
}
