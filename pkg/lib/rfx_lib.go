package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const Url = "https://api.remarkets.primary.com.ar/"
const auth = "auth/getToken"
const instruments = "rest/instruments/all"
const market_data = "rest/marketdata/get?marketId=ROFX&symbol=%v&entries=BI,OF,LA,OP,CL,HI,LO,SE,OI&depth=%v"

var MesesToString = map[string]string{"ENE": "01", "FEB": "02", "MAR": "03", "ABR": "04", "MAY": "05", "JUN": "06",
	"JUL": "07", "AGO": "08", "SEP": "9", "OCT": "10", "NOV": "11", "DIC": "12"}

var MesesToInt = map[string]int{"ENE": 1, "FEB": 2, "MAR": 3, "ABR": 4, "MAY": 5, "JUN": 6,
	"JUL": 7, "AGO": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DIC": 2}

var MesesToMonth = map[string]string{"ENE": "JAN", "FEB": "FEB", "MAR": "MAR", "ABR": "APR", "MAY": "05", "JUN": "06",
	"JUL": "07", "AGO": "AUG", "SEP": "SEP", "OCT": "OCT", "NOV": "NOV", "DIC": "DIC"}

/*
BI: BIDS Mejor oferta de compra en el Book
OF: OFFERS Mejor oferta de venta en el Book
LA: LAST Último precio operado en el mercado
OP: OPENING PRICE Precio de apertura
CL: CLOSING PRICE Precio de cierre
SE: SETTLEMENT PRICE Precio de ajuste (solo para futuros)
HI: TRADING SESSION HIGH PRICE Precio máximo de la rueda
LO: TRADING SESSION LOW PRICE Precio mínimo de la rueda
TV: TRADE VOLUME Volumen operado en contratos/nominales para ese security
OI: OPEN INTEREST Interés abierto (solo para futuros)
IV: INDEX VALUE Valor del índice (solo para índices)
EV: TRADE EFFECTIVE VOLUME Volumen efectivo de negociación para ese security
NV: NOMINAL VOLUME Volumen nominal de negociación para ese security
*/

func MarketDataUrl(symbol string, depth int) string {
	return fmt.Sprintf(Url+market_data, symbol, depth)
}

const Url_Auth = Url + auth
const Url_All_Instruments = Url + instruments

func rfx_get_req(url, token string) ([]byte, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Auth-Token", token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %v", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body) // response body is []byte
	if err != nil {
		return nil, fmt.Errorf("byte to string fail")
	}

	return body, nil
}

type RespAllInstruments struct {
	Status      string `json:"status"`
	Instruments []struct {
		InstrumentID struct {
			MarketID string `json:"marketId"`
			Symbol   string `json:"symbol"`
		} `json:"instrumentId"`
		Cficode string `json:"cficode"`
	} `json:"instruments"`
}

// https://mholt.github.io/json-to-go/

func UnmarshalAllInstruments(body []byte) ([]string, error) {
	var result RespAllInstruments
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		return nil, fmt.Errorf("can not unmarshal json")
	}

	var list_of_instruments []string

	if result.Status == "OK" {
		for _, inst := range result.Instruments {
			list_of_instruments = append(list_of_instruments, inst.InstrumentID.Symbol)
		}
	}

	return list_of_instruments, nil
}

type marketData struct {
	Status     string `json:"status"`
	MarketData struct {
		Oi struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"OI"`
		Cl struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"CL"`
		Hi float64 `json:"Hi"`
		Lo float64 `json:"LO"`
		Of []struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
		} `json:"OF"`
		Se struct {
			Price float64 `json:"price"`
			Size  int64   `json:"size"`
			Date  int64   `json:"date"`
		} `json:"SE"`
		La struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"LA"`
		Bi []struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
		} `json:"BI"`
		Op float64 `json:"OP"`
	} `json:"marketData"`
	Depth      int  `json:"depth"`
	Aggregated bool `json:"aggregated"`
}

func GetMarketData(contract, token string) (marketData, error) {

	url := MarketDataUrl(contract, 2)
	res, err := rfx_get_req(url, token)
	var unmarshaled_data marketData
	if err != nil {
		return unmarshaled_data, fmt.Errorf("error %v", err)
	}
	err = json.Unmarshal(res, &unmarshaled_data)
	if unmarshaled_data.Status != "OK" {

		return unmarshaled_data, fmt.Errorf("error unmarshall")

	}
	return unmarshaled_data, err
}

func LastPrice(ticker, token string) (float64, error) {
	data, err := GetMarketData(ticker, token)
	if err != nil {
		return 0, err
	}

	return data.MarketData.La.Price, err
}

type optionContract struct {
	Underlying   string
	Symbol       string
	K            float64
	Type         string
	MaturityDate string
	TtmInDays    float64
}

type OptionParameters struct {
	Tipo                 string
	S, K, T, R, Sigma, Q float64
}

func parseOptionContract(e string) optionContract {
	split1 := strings.Split(e, "/")
	split2 := strings.Split(split1[1], " ")
	fecha := split2[0]
	month := fecha[0:3]
	year := "20" + fecha[3:5]
	K, _ := strconv.ParseFloat(split2[1], 64)
	tipo := split2[2]
	month_index := MesesToString[month]
	maturityDate := fmt.Sprintf("%v-%v-%v", year, month_index, "28")
	_, secondsTtm := ParserStringDate(maturityDate)
	ttmInDays := TtmInDays(secondsTtm)

	return optionContract{Underlying: strings.Split(e, " ")[0],
		Symbol: e, K: K,
		Type: tipo, MaturityDate: maturityDate,
		TtmInDays: ttmInDays}
}

func mapOptions(key string) []string {
	mapForOptions := make(map[string][]string)

	mapForOptions["SOJ.ROS"] = append(mapForOptions["SOJ.ROS"], `^SOJ.ROS.*P$`)
	mapForOptions["SOJ.ROS"] = append(mapForOptions["SOJ.ROS"], `^SOJ.ROS.*C$`)

	mapForOptions["MAI.ROS"] = append(mapForOptions["MAI.ROS"], `^MAI.ROS.*P$`)
	mapForOptions["MAI.ROS"] = append(mapForOptions["MAI.ROS"], `^MAI.ROS.*C$`)

	mapForOptions["TRI.ROS"] = append(mapForOptions["TRI.ROS"], `^TRI.ROS.*P$`)
	mapForOptions["TRI.ROS"] = append(mapForOptions["TRI.ROS"], `^TRI.ROS.*C$`)

	mapForOptions["DLR"] = append(mapForOptions["DLR"], `^DLR.*P$`)
	mapForOptions["DLR"] = append(mapForOptions["DLR"], `^DLR.*C$`)

	mapForOptions["GIR.ROS"] = append(mapForOptions["GIR.ROS"], `^GIR.ROS.*P$`)
	mapForOptions["GIR.ROS"] = append(mapForOptions["GIR.ROS"], `^GIR.ROS.*C$`)

	return mapForOptions[key]
}

func compileRegExForUnderlying(underlying string) (*regexp.Regexp, *regexp.Regexp) {
	const regExForCall = 1
	const regExForPut = 0
	regexCall := mapOptions(underlying)[regExForCall]
	regexPut := mapOptions(underlying)[regExForPut]
	reCall, _ := regexp.Compile(regexCall)
	rePut, _ := regexp.Compile(regexPut)

	return reCall, rePut
}

func AllOptionsContract(underlying string, all_instruments []string) []optionContract {
	reCall, rePut := compileRegExForUnderlying(underlying)

	var contracts []optionContract
	for _, instrument := range all_instruments {
		matched := reCall.MatchString(instrument) || rePut.MatchString(instrument)
		if matched {
			optContract := parseOptionContract(instrument)
			contracts = append(contracts, optContract)
		}
	}
	return contracts
}

func GetAllInstruments(token string) ([]string, error) {
	response, err := rfx_get_req(Url_All_Instruments, token)
	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return UnmarshalAllInstruments(response)
}
