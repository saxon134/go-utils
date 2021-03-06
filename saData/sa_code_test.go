package saData

import (
	"fmt"
	"testing"
)

func TestId2Code(t *testing.T) {
	ids := []int64{798, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 233, 4578, 1989, 67273, 8856370, 17666544, 1111111, 1000, 10000, 100000, 1000000, 100000000, 1000000000}
	for _, id := range ids {
		encoded := Id2Code(id, 6)
		decoded := Code2Id(encoded, 8)
		fmt.Println(id, " => ", encoded, " => ", decoded)
	}
}
