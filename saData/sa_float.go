package saData

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

func RoundFloat64(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

func RoundFloat32(f float32, n int) float32 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return float32(inst)
}

// IntToFloat
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// digit是小数点后位数
func IntToFloat(value interface{}, digit int, roundType ...RoundType) float32 {
	if digit < 0 {
		return 0
	}

	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	div := decimal.New(1, int32(digit))

	var decimalValue decimal.Decimal
	switch v := value.(type) {
	case int:
		decimalValue = decimal.NewFromFloat(float64(v))
	case float64:
		decimalValue = decimal.NewFromFloat(v)
	case string:
		decimalValue, _ = decimal.NewFromString(v)
	case float32:
		decimalValue = decimal.NewFromFloat(float64(v))
	case int32:
		decimalValue = decimal.NewFromFloat(float64(v))
	case int64:
		decimalValue = decimal.NewFromFloat(float64(v))
	}

	//负数则返回0
	if decimalValue.Sign() <= 0 {
		return 0
	}

	f, _ := decimalValue.Float64()
	i := decimalValue.IntPart()

	//四舍五入
	if round == RoundTypeDefault {
		if f-float64(i) >= 0.5 {
			f, _ = decimal.NewFromInt(i + 1).Div(div).Float64()
			return float32(f)
		} else {
			f, _ = decimal.NewFromInt(i).Div(div).Float64()
			return float32(f)
		}
	} else
	//向上取整
	if round == RoundTypeUp {
		if f-float64(i) > 0 {
			f, _ = decimal.NewFromInt(i + 1).Div(div).Float64()
			return float32(f)
		} else {
			f, _ = decimal.NewFromInt(i).Div(div).Float64()
			return float32(f)
		}
	} else
	//向下取整
	if round == RoundTypeDown {
		f, _ = decimal.NewFromInt(i).Div(div).Float64()
		return float32(f)
	}
	return 0
}

// FloatToInt
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// digit是小数点后位数
func FloatToInt(value interface{}, digit int, roundType ...RoundType) int {
	if digit < 0 {
		return 0
	}

	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	mul := decimal.New(1, int32(digit))

	var decimalValue decimal.Decimal
	switch v := value.(type) {
	case int:
		decimalValue = decimal.NewFromFloat(float64(v)).Mul(mul)
	case int64:
		decimalValue = decimal.NewFromFloat(float64(v)).Mul(mul)
	case string:
		d, _ := decimal.NewFromString(v)
		f, _ := d.Float64()
		decimalValue = decimal.NewFromFloat(f).Mul(mul)
	case float32:
		decimalValue = decimal.NewFromFloat(float64(v)).Mul(mul)
	case float64:
		decimalValue = decimal.NewFromFloat(v).Mul(mul)
	case int32:
		decimalValue = decimal.NewFromFloat(float64(v)).Mul(mul)
	default:
		s := fmt.Sprint(value)
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			decimalValue = decimal.NewFromFloat(f).Mul(mul)
		}
	}

	f, _ := decimalValue.Float64()
	i := decimalValue.IntPart()

	//四舍五入
	if round == RoundTypeDefault {
		if f-float64(i) >= 0.5 {
			return int(i + 1)
		} else {
			return int(i)
		}
	} else
	//向上取整
	if round == RoundTypeUp {
		if f-float64(i) > 0 {
			return int(i + 1)
		} else {
			return int(i)
		}
	} else
	//向下取整
	if round == RoundTypeDown {
		return int(i)
	}
	return 0
}
