package main

import (
	"fmt"

	"time"

	lib "github.com/lmpizarro/gorofex/pkg/lib"
)


func scan(token string) {
	for {
		fmt.Println("--- ", time.Now().Second())
		for _, ticker :=
		    range []string{"DLR/FEB23", "DLR/MAR23",
			"DLR/ABR23", "DLR/MAY23", "DLR/JUN23",
			"DLR/JUL23", "DLR/AGO23", "DLR/SEP23",
			"DLR/NOV23", "DLR/DIC23",
			} {
			cl, _ := lib.Last_Price(ticker, token)
			fmt.Println(ticker, " ", cl)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(1 * time.Second)
		fmt.Println("Press the Enter Key to stop anytime")
	}
}

func main() {

	cred := lib.Read_credentials("./env.csv")
	user := cred.User
	password := cred.Password

	token := lib.Token(user, password)

	go scan(token)
	fmt.Scanln()
}
