package main

import (
	"fmt"

	"time"

	lib "github.com/lmpizarro/gorofex/pkg/lib"
)

func scan() {
	token := lib.Login()
	for {
		fmt.Println("--- ", time.Now().Second())
		for _, ticker := range []string{"DLR/FEB23", "DLR/MAR23",
			"DLR/ABR23", "DLR/MAY23", "DLR/JUN23",
			"DLR/JUL23", "DLR/AGO23", "DLR/SEP23",
			"DLR/NOV23", "DLR/DIC23", "DLR/ENE24",
		} {
			cl, _ := lib.LastPrice(ticker, token)
			fmt.Println(ticker, " ", cl)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(30 * time.Second)
		fmt.Println("Press the Enter Key to stop anytime")
	}
}

func main() {

	// go scan()
	// fmt.Scanln()
	// panic("")

	token := lib.Login()
	allInstruments, _ := lib.GetAllInstruments(token)

	contracts := lib.AllOptionsContract("SOJ.ROS", allInstruments)
	for _, contract := range contracts {
		lprice, _ := lib.LastPrice(contract.Symbol, token)
		fmt.Println(contract.Symbol, lprice, contract.K, contract.MaturityDate)
	}

	map_messages, _ := lib.Build_messages("GGAL/FEB23", lib.Login())

	for k := range map_messages {
		fmt.Println(k, ".... ", map_messages[k])
	}
}
