package saData

import (
	"fmt"
	"testing"
)

func TestDataConvert(t *testing.T) {
	fmt.Println(Yuan2Fen(19.9*100, RoundTypeDefault))
	fmt.Println(Yuan2Fen(1.230, RoundTypeDefault))
	fmt.Println(Yuan2Fen(1.235, RoundTypeDefault))
	fmt.Println(Yuan2Fen(1.239, RoundTypeDefault))
	fmt.Println("1=====")

	fmt.Println(Yuan2Fen(19.9*100, RoundTypeUp))
	fmt.Println(Yuan2Fen(1.230, RoundTypeUp))
	fmt.Println(Yuan2Fen(1.231, RoundTypeUp))
	fmt.Println(Yuan2Fen(1.235, RoundTypeUp))
	fmt.Println(Yuan2Fen(1.239, RoundTypeUp))
	fmt.Println("2=====")

	fmt.Println(Yuan2Fen(19.9*100, RoundTypeDown))
	fmt.Println(Yuan2Fen(1.230, RoundTypeDown))
	fmt.Println(Yuan2Fen(1.235, RoundTypeDown))
	fmt.Println(Yuan2Fen(1.239, RoundTypeDown))
	fmt.Println("3=====")

	fmt.Println(Fen2Yuan(1990, RoundTypeDefault))
	fmt.Println(Fen2Yuan(123.1, RoundTypeDefault))
	fmt.Println(Fen2Yuan(123.5, RoundTypeDefault))
	fmt.Println(Fen2Yuan(123.9, RoundTypeDefault))
	fmt.Println("4=====")

	fmt.Println(Fen2Yuan(1990, RoundTypeUp))
	fmt.Println(Fen2Yuan(123.0, RoundTypeUp))
	fmt.Println(Fen2Yuan(123.1, RoundTypeUp))
	fmt.Println(Fen2Yuan(123.5, RoundTypeUp))
	fmt.Println(Fen2Yuan(123.9, RoundTypeUp))
	fmt.Println("5=====")

	fmt.Println(Fen2Yuan(19.9*100, RoundTypeDown))
	fmt.Println(Fen2Yuan(123.0, RoundTypeDown))
	fmt.Println(Fen2Yuan(123.1, RoundTypeDown))
	fmt.Println(Fen2Yuan(123.5, RoundTypeDown))
	fmt.Println(Fen2Yuan(123.9, RoundTypeDown))
	fmt.Println("6=====")
}
