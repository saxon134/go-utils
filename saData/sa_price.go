package saData

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

type RoundType int8

const (
	RoundTypeDefault RoundType = iota //四舍五入
	RoundTypeUp                       //向上取整 取2位小数 如：1.234 => 1.24  1.238 => 1.24
	RoundTypeDown                     //向下取整 取2位小数 如：1.234 => 1.23  1.238 => 1.23
)

// RoundYuan
// 格式化，解决小数点问题
// 提供四舍五入、上下取整方式
// 返回的元是2位小数
func RoundYuan(yuan interface{}, roundType ...RoundType) float64 {
	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	var yuanDecimal decimal.Decimal
	switch v := yuan.(type) {
	case int:
		yuanDecimal = decimal.NewFromFloat(float64(v))
	case float64:
		yuanDecimal = decimal.NewFromFloat(v)
	case string:
		yuanDecimal, _ = decimal.NewFromString(v)
	case float32:
		yuanDecimal = decimal.NewFromFloat(float64(v))
	case int32:
		yuanDecimal = decimal.NewFromFloat(float64(v))
	case int64:
		yuanDecimal = decimal.NewFromFloat(float64(v))
	}

	//负数则返回0
	if yuanDecimal.Sign() <= 0 {
		return 0
	}

	f, _ := yuanDecimal.Float64()
	i := yuanDecimal.IntPart()

	//四舍五入
	if round == RoundTypeDefault {
		if f-float64(i) >= 0.5 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(100)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
			return f
		}
	} else
	//向上取整
	if round == RoundTypeUp {
		if f-float64(i) > 0 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(100)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
			return f
		}
	} else
	//向下取整
	if round == RoundTypeDown {
		f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
		return f
	}
	return 0
}

// Fen2Yuan
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// 不可为负数，负数则返回0
// 返回的元是2位小数
func Fen2Yuan(fen interface{}, roundType ...RoundType) (yuan float64) {
	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	var fenDecimal decimal.Decimal
	switch v := fen.(type) {
	case int:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case float64:
		fenDecimal = decimal.NewFromFloat(v)
	case string:
		fenDecimal, _ = decimal.NewFromString(v)
	case float32:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case int32:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case int64:
		fenDecimal = decimal.NewFromFloat(float64(v))
	}

	//负数则返回0
	if fenDecimal.Sign() <= 0 {
		return 0
	}

	f, _ := fenDecimal.Float64()
	i := fenDecimal.IntPart()

	//四舍五入
	if round == RoundTypeDefault {
		if f-float64(i) >= 0.5 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(100)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
			return f
		}
	} else
	//向上取整
	if round == RoundTypeUp {
		if f-float64(i) > 0 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(100)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
			return f
		}
	} else
	//向下取整
	if round == RoundTypeDown {
		f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(100)).Float64()
		return f
	}
	return
}

// Yuan2Fen
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// 不可为负数，负数则返回0
func Yuan2Fen(yuan interface{}, roundType ...RoundType) (fen int) {
	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	var fenDecimal decimal.Decimal
	switch v := yuan.(type) {
	case int:
		fenDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(100))
	case int64:
		fenDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(100))
	case string:
		d, _ := decimal.NewFromString(v)
		f, _ := d.Float64()
		fenDecimal = decimal.NewFromFloat(f).Mul(decimal.NewFromInt(100))
	case float32:
		fenDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(100))
	case float64:
		fenDecimal = decimal.NewFromFloat(v).Mul(decimal.NewFromInt(100))
	case int32:
		fenDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(100))
	default:
		s := fmt.Sprint(yuan)
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			fenDecimal = decimal.NewFromFloat(f).Mul(decimal.NewFromInt(100))
		}
	}

	//负数则返回0
	if fenDecimal.Sign() <= 0 {
		return 0
	}

	f, _ := fenDecimal.Float64()
	i := fenDecimal.IntPart()

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
	return
}

// Li2Yuan
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// 不可为负数，负数则返回0
// 返回的元是3位小数
func Li2Yuan(li interface{}, roundType ...RoundType) (yuan float64) {
	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	var fenDecimal decimal.Decimal
	switch v := li.(type) {
	case int:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case float64:
		fenDecimal = decimal.NewFromFloat(v)
	case string:
		fenDecimal, _ = decimal.NewFromString(v)
	case float32:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case int32:
		fenDecimal = decimal.NewFromFloat(float64(v))
	case int64:
		fenDecimal = decimal.NewFromFloat(float64(v))
	}

	//负数则返回0
	if fenDecimal.Sign() <= 0 {
		return 0
	}

	f, _ := fenDecimal.Float64()
	i := fenDecimal.IntPart()

	//四舍五入
	if round == RoundTypeDefault {
		if f-float64(i) >= 0.5 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(1000)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(1000)).Float64()
			return f
		}
	} else
	//向上取整
	if round == RoundTypeUp {
		if f-float64(i) > 0 {
			f, _ = decimal.NewFromInt(i + 1).Div(decimal.NewFromInt(1000)).Float64()
			return f
		} else {
			f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(1000)).Float64()
			return f
		}
	} else
	//向下取整
	if round == RoundTypeDown {
		f, _ = decimal.NewFromInt(i).Div(decimal.NewFromInt(1000)).Float64()
		return f
	}
	return
}

// Yuan2Li
// 解决精度丢失问题 例如：19.9 * 100 = 1889
// 提供四舍五入、上下取整方式
// 不可为负数，负数则返回0
func Yuan2Li(yuan interface{}, roundType ...RoundType) (li int) {
	round := RoundTypeDefault
	if roundType != nil && len(roundType) > 0 {
		round = roundType[0]
	}

	var liDecimal decimal.Decimal
	switch v := yuan.(type) {
	case int:
		liDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(1000))
	case int64:
		liDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(1000))
	case string:
		d, _ := decimal.NewFromString(v)
		f, _ := d.Float64()
		liDecimal = decimal.NewFromFloat(f).Mul(decimal.NewFromInt(1000))
	case float32:
		liDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(1000))
	case float64:
		liDecimal = decimal.NewFromFloat(v).Mul(decimal.NewFromInt(1000))
	case int32:
		liDecimal = decimal.NewFromFloat(float64(v)).Mul(decimal.NewFromInt(1000))
	default:
		s := fmt.Sprint(yuan)
		if f, err := strconv.ParseFloat(s, 32); err == nil {
			liDecimal = decimal.NewFromFloat(f).Mul(decimal.NewFromInt(1000))
		}
	}

	//负数则返回0
	if liDecimal.Sign() <= 0 {
		return 0
	}

	f, _ := liDecimal.Float64()
	i := liDecimal.IntPart()

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
	return
}
