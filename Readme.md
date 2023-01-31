A very basic Golang interface to Rofex, draft, as is

License Apache 2.0

#### code

`

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

	// all_instruments, _ := lib.Get_All_Instruments(token)
	// s := []string{"SOJ.ROS", "TRI.ROS", "MAI.ROS", "MERV - XMEV", "DLR", "ORO", "WTI", ".MIN", "CNH", ".CME", "CONT", "I."}
	// contracts := lib.All_options_contract("MAI.ROS", all_instruments)

`

`

// https://golangdocs.com/golang-finance-go-package-stock-quote-options-chart

import (
	"math"
	"gonum.org/v1/gonum/stat/distuv"
)

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


`

## Learning Go
[empty-vs-nil-slices-golang](https://blog.boot.dev/golang/empty-vs-nil-slices-golang/)

[store password](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/09.5.html)

[contains](https://play.golang.org/p/Qg_uv_inCek)
