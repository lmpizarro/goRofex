package lib

import (
	"fmt"
	"math"
	"time"

	fin "github.com/alpeb/go-finance/fin"
)

type CashFlow struct {
	cash  []float64
	dates []string
	name  string
	Irr float64
	Rate float64
}

func Irrs(c *CashFlow) {

	var dates []time.Time
	for _, date := range c.dates {
		t, _ := ParserStringDate(date)
		dates = append(dates, t)
	}

	r1, _ := fin.InternalRateOfReturn(c.cash, .1)
	r2, _ := fin.ScheduledInternalRateOfReturn(c.cash, dates, 1)

	c.Irr = r2
	c.Rate = r1

}

func TestCashFlow() {
	begin_date := "2023-03-07"
	s31m3 := CashFlow{cash: []float64{-95.1, 100.0}, name: "S31M3",
		dates: []string{begin_date, "2023-03-31"}}

	s28a3 := CashFlow{cash: []float64{-89.17, 100.0}, name: "S28A3",
		dates: []string{begin_date, "2023-04-28"}}

	s31y3 := CashFlow{cash: []float64{-82.97, 100.0}, name: "S31Y3",
		dates: []string{begin_date, "2023-05-31"}}

	s30j3 := CashFlow{cash: []float64{-77.75, 100.0}, name: "S30J3",
		dates: []string{begin_date, "2023-06-30"}}
		
	x16j3 := CashFlow{cash: []float64{-142.5, 144.5}, name: "X16J3",
		dates: []string{begin_date, "2023-06-16"}}

	cfs := []CashFlow{s31m3, s28a3, s31y3, s30j3, x16j3}

	for _, cf := range cfs {
		Irrs(&cf)
		fmt.Println(cf.Irr, cf.Rate, cf.name)
	}
}

// TimeFactor in years
type RateTerm struct {
	Rate       float64
	TimeFactor float64
}

func CompoundRates(rf []RateTerm) float64 {
	accumulatedRate := 1.0
	for _, e := range rf {
		f := math.Exp(e.Rate * e.TimeFactor)
		accumulatedRate *= f
	}
	return accumulatedRate
}

func modeloCer(cer0, cer1 float64, term RateTerm) float64 {
	return math.Exp(term.TimeFactor*term.Rate) * cer1 / cer0
}

func DiscountFactor(rt * RateTerm, c string) float64 {
	var df float64
	if c == "simp" {
		df =  1.0 / (1.0 + rt.Rate*rt.TimeFactor)
	} else if c == "comp" {
		df =  math.Pow(1.0 + rt.Rate, rt.TimeFactor)
	} else if c == "cont" {
		df =  math.Exp(-rt.Rate*rt.TimeFactor)
	}
	return df
}

func ForwardRate(term1, term2 RateTerm, c string) (float64, error){
	if term1.TimeFactor > term2.TimeFactor {
		return 0.0, fmt.Errorf("error")
	}
	var fr float64
	df1 := DiscountFactor(&term1, c)
	df2 := DiscountFactor(&term2, c)
	time_diff := term2.TimeFactor - term1.TimeFactor
	if c == "simp" {
		fr =  ((df1/df2) - 1) / time_diff
	} else if c == "comp" {
		fr = math.Pow(df2/df1, 1 / time_diff) - 1
	} else if c == "cont" {
		fr = math.Log(df1/df2) / time_diff
	}

	return fr, nil
}

func TestForwardRate() {

	for _, t := range []string{"comp", "simp", "cont"} {
		fr, _ := ForwardRate(RateTerm{Rate: 0.01, TimeFactor: 0.09},
			RateTerm{Rate: 0.05, TimeFactor: 0.5}, t)
		fmt.Println(fr)
	}

	for _, t := range []string{"comp", "simp", "cont"} {
		fr, _ := ForwardRate(RateTerm{Rate: 0.07, TimeFactor: 0.25},
			RateTerm{Rate: 0.05, TimeFactor: 0.5}, t)
		fmt.Println(fr)
	}
}


func modeloInflacion() float64{

	inflat := []RateTerm{
		{Rate: 0.06, TimeFactor: 1.0},
		{Rate: 0.058, TimeFactor: 1.0},
		{Rate: 0.057, TimeFactor: 1.0},
		{Rate: 0.056, TimeFactor: .5},
	}
	inflationAccum := CompoundRates(inflat)

	return inflationAccum
}

func TestCer() {
	cerBegin := 43.9619 // 55.5275
	cerOperation := 78.6758 // 2023-02-10
	// cerEnd := cerOperation * inflationAccum
	precio := 1.7775
	_, dateBegin := ParserStringDate("2022-04-29")  //  2022-08-31
	dateTime2, dateOperation := ParserStringDate("2023-02-27")
	dateTime3, _ := ParserStringDate("2023-05-19")  // 2023-05-19 2023-06-16

	diffOperBeg := fin.DaysDifference(dateBegin, dateOperation, 1)
	// diffEndBegin := fin.DaysDifference(dateBegin, dateEnd, 1)
	// fmt.Println(diffOperBeg, diffEndBegin, inflationAccum)

	disc1 := RateTerm{Rate: 0.000, TimeFactor: float64(diffOperBeg) / 365.0}
	// disc2 := RateTerm{Rate: 0.0051, TimeFactor: float64(diffEndBegin) / 365.0}
	estimate1 := modeloCer(cerBegin, cerOperation, disc1)
	// fmt.Println(estimate1, CompoundRates([]RateTerm{disc1}))
	// fmt.Println(modeloCer(cerBegin, cerEnd, disc2), CompoundRates([]RateTerm{disc2}))
	tir,_ :=fin.ScheduledInternalRateOfReturn([]float64{-precio, estimate1}, []time.Time{dateTime2, dateTime3}, 0.01)
	fmt.Println(tir, estimate1)

}
