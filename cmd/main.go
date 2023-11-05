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
		for _, ticker := range []string{
			"DLR/AGO23", "DLR/SEP23",
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


	lib.DatosBCRA(lib.SetDataSource("2022-02-21", "2023-02-21", "cer"))
	//panic("")
	// go scan()
	// fmt.Scanln()
	// panic("")

	lib.TestCashFlow()
	// panic("")
	token := lib.Login()
	allInstruments, _ := lib.GetAllInstruments(token)

	optionContracts := lib.AllOptionsContract("SOJ.ROS", allInstruments)
	for _, contract := range optionContracts {
		op, _ := lib.LastPrice(contract.Underlying, token)
		fmt.Printf("%s %.2f %.2f \n", contract.Symbol, contract.TtmInDays/365.0, op)
		cp, _ := lib.LastPrice(contract.Symbol, token)
		if cp != 0.0 {
			fmt.Println("cp ", cp)
		}
	}

	map_messages, _ := lib.Build_messages("GGAL/DIC23", lib.Login())

	for k := range map_messages {
		fmt.Println(k, ".... ", map_messages[k])
	}
}
