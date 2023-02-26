package main

import (
	"fmt"

	"time"

	lib "github.com/lmpizarro/gorofex/pkg/lib"
	fin "github.com/alpeb/go-finance/fin"
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

type cashFlow struct {
	cash []float64
	dates []string
	name string
}

func Irrs(c *cashFlow) (float64, float64) {

	var dates []time.Time
	for _, date := range c.dates {
		t, _ := lib.ParserStringDate(date)
		dates = append(dates, t)
	}

	r1, _ := fin.InternalRateOfReturn(c.cash, .1)
	r2, _ := fin.ScheduledInternalRateOfReturn(c.cash, dates, 1)

	return r1, r2
}

func main() {

	// go scan()
	// fmt.Scanln()
	// panic("")

	begin_date := "2023-02-27"
	cfs := []cashFlow {{cash: []float64{-93.7, 100.0}, name: "S31M3",
	            dates: []string{begin_date, "2023-03-31"}},
					    {cash: []float64{-88.25, 100.0}, name: "S28A3",
	            dates: []string{begin_date, "2023-04-28"}},
					    {cash: []float64{-81.98, 100.0}, name: "S31Y3",
	            dates: []string{begin_date, "2023-05-31"}},
					    {cash: []float64{-76.98, 100.0}, name: "S30J3",
	            dates: []string{begin_date, "2023-06-30"}},
				// 31/08/2022 101 55.5275  78,6758
					    {cash: []float64{-139.5, 174.29}, name: "X16J3",
	            dates: []string{begin_date, "2023-06-16"}},

			}

	for _, cf := range cfs {
		r1, r2 := Irrs(&cf)
		fmt.Println(r1, r2)
	}


	panic("")

	token := lib.Login()
	allInstruments, _ := lib.GetAllInstruments(token)
	for _, inst := range allInstruments {
		fmt.Println(inst)
	}

	optionContracts := lib.AllOptionsContract("GIR.ROS", allInstruments)
	for _, contract := range optionContracts {
		op, _ := lib.LastPrice(contract.Position, token)
		fmt.Printf("%s %.2f %.2f\n", contract.Position, contract.TtmInDays/365.0, op)
	}
	panic("")
	map_messages, _ := lib.Build_messages("GGAL/FEB23", lib.Login())

	for k := range map_messages {
		fmt.Println(k, ".... ", map_messages[k])
	}
}
