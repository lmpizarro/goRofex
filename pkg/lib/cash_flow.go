package lib

import (
	"time"
	fin "github.com/alpeb/go-finance/fin"
	"fmt"
	"math"
)

type CashFlow struct {
	cash []float64
	dates []string
	name string
}

func Irrs(c *CashFlow) (float64, float64) {

	var dates []time.Time
	for _, date := range c.dates {
		t, _ := ParserStringDate(date)
		dates = append(dates, t)
	}

	r1, _ := fin.InternalRateOfReturn(c.cash, .1)
	r2, _ := fin.ScheduledInternalRateOfReturn(c.cash, dates, 1)

	return r1, r2
}

func TestCashFlow() {
	begin_date := "2023-02-27"
	cfs := []CashFlow {{cash: []float64{-93.7, 100.0}, name: "S31M3",
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
}

type RateFactor struct {
	Rate float64
	Term float64
}

func AccumulatedInflation(rf []RateFactor) float64{
	accumulatedInflation := 1.0
	for _, e := range rf{
		f := math.Exp(e.Rate * e.Term)
		accumulatedInflation *= f
	}
	return accumulatedInflation
}

func TestInflation(){
	inflat := []RateFactor{
		{Rate: 0.06, Term: 1.0},
		{Rate: 0.058, Term: 1.0},
		{Rate: 0.057, Term: 1.0},
		{Rate: 0.056, Term: .5},
	}

	fmt.Println(AccumulatedInflation(inflat))

}