package lib

import "github.com/PuerkitoBio/goquery"
import "encoding/json"
import "strings"
import "net/http"
import "fmt"

func Cer(fechaDesde, fechaHasta string){
	splitDesde := strings.Split(fechaDesde, "-")
	splitHasta := strings.Split(fechaHasta, "-")

	fechaDesde = fmt.Sprintf("%v%v%v", splitDesde[0], splitDesde[1], splitDesde[2])
	fechaHasta = fmt.Sprintf("%v%v%v", splitHasta[0], splitHasta[1], splitHasta[2])

    myurl := "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp"
	r, err := http.NewRequest("POST", myurl, nil)

	if err != nil {
		panic(err)
	}
	form := r.URL.Query()
	form.Add("primeravez", "1")
	form.Add("fecha_desde", fechaDesde)
	form.Add("fecha_hasta", fechaHasta)
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
