package main

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"os"

	"encoding/csv"

	lib "github.com/lmpizarro/gorofex/pkg/lib"
	"gonum.org/v1/gonum/stat/distuv"

)

type credentials struct {
    User string
    Password string
    Account string
}

func read_credentials() credentials {

    csvFile, err := os.Open("./env.csv")
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

// https://golangdocs.com/golang-finance-go-package-stock-quote-options-chart

func bs(tipo string, S, K, T, r, sigma, div float64) float64 {

	/*
	   	Def
	          Calculador del precio de una opcion Europea con el modelo de Black Scholes
	      Inputs
	         - tipo : string - Tipo de contrato entre ["CALL","PUT"]
	         - S : float - Spot price del activo
	         - K : float - Strike price del contrato
	         - T : float - Tiempo hasta la expiracion (en años)
	         - r : float - Tasa 'libre de riesgo' (anualizada)
	         - sigma : float - Volatilidad implicita (anualizada)
	         - div : float - Tasa de dividendos continuos (anualizada)
	     Outputs
	         - precio_BS: float - Precio del contrato
	*/
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var price float64

	d1 := (math.Log(S/K) + (r-div+0.5*sigma*sigma)*T) / sigma / math.Sqrt(T)
	d2 := (math.Log(S/K) + (r-div-0.5*sigma*sigma)*T) / sigma / math.Sqrt(T)

	if tipo == "C" {
		price = math.Exp(-div*T)*S*dist.CDF(d1) - K*math.Exp(-r*T)*dist.CDF(d2)
	} else if tipo == "P" {
		price = K*math.Exp(-r*T)*dist.CDF(-d2) - S*math.Exp(-div*T)*dist.CDF(-d1)
	}

	return price
}

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
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

	cred := read_credentials()
	user := cred.User
	password := cred.Password

	token := token(user, password)

	// all_instruments, _ := lib.Get_All_Instruments(token)
	// s := []string{"SOJ.ROS", "TRI.ROS", "MAI.ROS", "MERV - XMEV", "DLR", "ORO", "WTI", ".MIN", "CNH", ".CME", "CONT", "I."}
	// contracts := lib.All_options_contract("MAI.ROS", all_instruments)

	e := "DLR/MAR23"
	fmt.Println(e)
	lib.Get_Market_Data(e, token)

}
