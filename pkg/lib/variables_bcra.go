package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const myUrl = "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp"

type FuenteDatos struct {
	Serie   string
	Detalle string
	FechaHasta string
	FechaDesde string
}

var DatosCer = FuenteDatos{Serie: "3540", Detalle: "CER (Base 2.2.2002=1)"}

var DatosBadlarEA = FuenteDatos{Serie: "7935", Detalle: "BADLAR en pesos de bancos privados (en  e.a.)"}

var DatosA3500 = FuenteDatos{Serie: "272", Detalle: "Tipo de Cambio Mayorista ($ por US$) Comunicación A 3500 - Referencia"}

func SetDataSource(fechaDesde, fechaHasta, fuente string) FuenteDatos {

	var s FuenteDatos
	if fuente == "cer" {
		s = DatosCer
	} else if fuente == "badlarEA"{
		s = DatosBadlarEA
	} else if fuente == "a3500" {
		s = DatosA3500
	} else {
		panic("not yet ready")
	}
	s.FechaDesde = fechaDesde
	s.FechaHasta = fechaHasta

	return s
}

func hacerFecha(fecha string) string {
	splitD := strings.Split(fecha, "-")
	return fmt.Sprintf("%v%v%v", splitD[0], splitD[1], splitD[2])
}

func DatosBCRA(service FuenteDatos) {
	r, err := http.NewRequest("POST", myUrl, nil)
	if err != nil {
		panic(err)
	}

	form := r.URL.Query()
	form.Add("primeravez", "1")
	form.Add("fecha_desde", hacerFecha(service.FechaDesde))
	form.Add("fecha_hasta", hacerFecha(service.FechaHasta))
	form.Add("serie", service.Serie)
	form.Add("serie1", "0")
	form.Add("serie2", "0")
	form.Add("serie3", "0")
	form.Add("serie4", "0")
	form.Add("detalle", service.Detalle)
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic("-....-...._")
	}
	// busca el par fecha cer
	doc.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
			case 0:
				fmt.Printf("%v ", td.Text())
			case 1:
				fmt.Printf("%v \n", strings.TrimSpace(td.Text()))
			}
		})
	})
}
