package lib

import (
	"fmt"

)
/*
XADD temperatures:us-ny:10007 * temp_f 87.2 pressure 29.69 humidity 46
*/

func Build_messages(ticker, token string) (map[string] string, error){
	m := make(map[string]string)

	data, err := Get_Market_Data(ticker, token)
	if err != nil {
		return m, err
	}
	m["CL"] = fmt.Sprintf("RFX:TCKR:%v:CL PRICE %v DATE %v SIZE %v", ticker,
		data.MarketData.Cl.Price, data.MarketData.Cl.Date, data.MarketData.Cl.Size)

	m["LA"] = fmt.Sprintf("RFX:TCKR:%v:LA PRICE %v DATE %v SIZE %v", ticker,
		data.MarketData.La.Price, data.MarketData.La.Date, data.MarketData.La.Size)

	m["HI"] = fmt.Sprintf("RFX:TCKR:%v:HI PRICE %v ", ticker, data.MarketData.Hi)

	m["LO"] = fmt.Sprintf("RFX:TCKR:%v:LO PRICE %v ", ticker, data.MarketData.Lo)

	m["OP"] = fmt.Sprintf("RFX:TCKR:%v:OP PRICE %v ", ticker, data.MarketData.Op)

	m["OHLC"] = fmt.Sprintf("RFX:TCKR:%v:OHLC O %v H %v L %v C %v", ticker,
		data.MarketData.Op, data.MarketData.Hi, data.MarketData.Lo, data.MarketData.Cl.Price)

	if len(data.MarketData.Of) > 0 {
		m["OF"] = fmt.Sprintf("RFX:TCKR:%v:OF PRICE %v SIZE %v", ticker,
		data.MarketData.Of[0].Price, data.MarketData.Of[0].Size)
	} else {
		m["OF"] = fmt.Sprintf("RFX:TCKR:%v:OF PRICE %v SIZE %v", ticker, "--.--", "--.--")
	}

	if len(data.MarketData.Bi) > 0 {
		m["BI"] = fmt.Sprintf("RFX:TCKR:%v:BI PRICE %v SIZE %v", ticker,
		data.MarketData.Bi[0].Price, data.MarketData.Bi[0].Size)
	} else {
		m["BI"] = fmt.Sprintf("RFX:TCKR:%v:BI PRICE %v SIZE %v", ticker, "--.--", "--.--")
	}


	fmt.Println(data.MarketData.Of)
	return m, nil
}

func Message_bykey(m map[string]string, key string) (string, error) {
	message := m[key]
	return message, nil
}

func Message_CL(m map[string]string) (string, error) {
	return Message_bykey(m, "CL")
}

func Message_LA(m map[string]string) (string, error) {
	return Message_bykey(m, "LA")
}

func Message_HI(m map[string]string) (string, error) {
	return Message_bykey(m, "HI")
}

func Message_LO(m map[string]string) (string, error) {
	return Message_bykey(m, "LO")
}

func Message_OP(m map[string]string) (string, error) {
	return Message_bykey(m, "OP")
}

func Message_OF(m map[string]string) (string, error) {
	return Message_bykey(m, "OF")
}

func Message_BI(m map[string]string) (string, error) {
	return Message_bykey(m, "BI")
}

func Message_OHLC(m map[string]string) (string, error) {
	return Message_bykey(m, "OHLC")
}

