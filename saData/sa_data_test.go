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

func TestId2Code(t *testing.T) {
	ids := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 233, 4578, 1989, 67273, 8856370, 17666544, 1111111, 1000, 10000, 100000, 1000000, 100000000, 1000000000}
	//for _, id := range ids {
	//	encoded := I64ToCode(id, 8)
	//	decoded := Code2Id(encoded, 8)
	//	fmt.Println(id, " => ", encoded, " => ", decoded)
	//}
	//
	//fmt.Println("========================")

	for _, id := range ids {
		encoded := IdToChar(id)
		decoded := CharToId(encoded)
		fmt.Println(id, " => ", encoded, " => ", decoded)
	}
}
