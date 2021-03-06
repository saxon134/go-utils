package saData

import (
	"errors"
	"strings"
)

var _source = []string{"2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

const salt = 367892
const prime1 = 7
const prime2 = 13

func Id2Code(id int64, length int) (code string) {
	return Id2CodeWithSource(id, length, _source)
}

func Code2Id(ivcode string, length int) (id int64) {
	return Code2IdWithSource(ivcode, length, _source)
}

func Id2CodeWithSource(id int64, length int, source []string) (code string) {
	if len(source) == 0 {
		source = _source
	}

	if length < 4 {
		length = 4
	}

	id = id * prime1

	// add salt
	id = id + salt
	numberResult := make([]int64, length)

	// transform to 32 base
	numberResult[0] = id
	for i := 0; i < length-1; i++ {
		numberResult[i+1] = numberResult[i] / 32
		numberResult[i] = (numberResult[i] + numberResult[0]*int64(i)) % 32
	}

	if numberResult[length-1] >= 32 {
		return
	}

	// p-box
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = source[numberResult[i*prime2%length]]
	}

	code = strings.Join(result, "")
	return
}

func Code2IdWithSource(ivcode string, length int, source []string) (id int64) {
	if len(source) == 0 {
		source = _source
	}

	if length < 4 {
		length = 4
	}

	if len(ivcode) != length {
		return
	}

	// return to int64
	result := strings.Split(ivcode, "")
	numberResult := make([]int64, length)
	numbers := make([]int64, length)
	for i := 0; i < length; i++ {
		index, err := findIndexOf(source, result[i])
		if err != nil {
			return
		}
		numberResult[i] = int64(index)
	}

	for i := 0; i < length; i++ {
		numbers[i] = numberResult[i*prime2%length]
	}

	b := make([]int64, length)
	for i := length - 2; i >= 0; i-- {
		b[i] = (numbers[i] - numbers[0]*int64(i) + 32*int64(i)) % 32
	}

	for i := length - 2; i > 0; i-- {
		id += b[i]
		id *= 32
	}
	id = (id + b[0] - salt) / prime1
	return
}

func findIndexOf(source []string, v string) (int, error) {
	for index, value := range source {
		if value == v {
			return index, nil
		}
	}
	return -1, errors.New("value not found")
}
