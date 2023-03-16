package saData

import (
	"errors"
	"strings"
)

////////////////////////////////////////////////////////////////
// int64和字符串互转，转成的字符串固定长度
// 数字范围较小，如8位字符串，最大能表示的数字为999999999
// 尽量使用saIdChar，只有在需要固定字符长度的情形使用saIdCode
////////////////////////////////////////////////////////////////

var _source = []string{
	"2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

const salt = 192
const prime1 = 7
const prime2 = 13

func IdToCode(id int64, length int) (code string) {
	return IdToCodeWithSource(id, length, _source)
}

func CodeToId(code string, length int) (id int64) {
	return CodeToIdWithSource(code, length, _source)
}

func IdToCodeWithSource(id int64, length int, source []string) (code string) {
	if len(source) == 0 {
		source = _source
	}

	if length < 3 {
		length = 3
	}

	id = id * prime1

	// add salt
	id = id + salt
	numberResult := make([]int64, length)

	// transform to 32 base
	SourceLen := int64(len(source))
	numberResult[0] = id
	for i := 0; i < length-1; i++ {
		numberResult[i+1] = numberResult[i] / SourceLen
		numberResult[i] = (numberResult[i] + numberResult[0]*int64(i)) % SourceLen
	}

	if numberResult[length-1] >= SourceLen {
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

func CodeToIdWithSource(code string, length int, source []string) (id int64) {
	if len(source) == 0 {
		source = _source
	}

	if length < 3 {
		length = 3
	}

	if len(code) != length {
		return
	}

	// return to int64
	result := strings.Split(code, "")
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

	SourceLen := int64(len(source))
	b := make([]int64, length)
	for i := length - 2; i >= 0; i-- {
		b[i] = (numbers[i] - numbers[0]*int64(i) + SourceLen*int64(i)) % SourceLen
	}

	for i := length - 2; i > 0; i-- {
		id += b[i]
		id *= SourceLen
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
