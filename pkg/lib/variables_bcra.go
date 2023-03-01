package lib

import "github.com/PuerkitoBio/goquery"
import "strings"
import "net/http"
import "fmt"

const myurl = "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp"

func hacerFecha(fecha string) string {
	splitD:= strings.Split(fecha, "-")
	return fmt.Sprintf("%v%v%v", splitD[0], splitD[1], splitD[2])
}

func Cer(fechaDesde, fechaHasta string){
	r, err := http.NewRequest("POST", myurl, nil)
	if err != nil {
		panic(err)
	}

	form := r.URL.Query()
	form.Add("primeravez", "1")
	form.Add("fecha_desde", hacerFecha(fechaDesde))
	form.Add("fecha_hasta", hacerFecha(fechaHasta))
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic("-....-...._")
	}
 	// busca el par fecha cer
  	doc.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
		tr.Find("td").Each(func(ix int, td *goquery.Selection) {
			switch ix {
				case 0:
					fmt.Printf("%v ",  td.Text())
				case 1:
					fmt.Printf("%v \n", strings.TrimSpace(td.Text()))
			}
		})
	})
}
