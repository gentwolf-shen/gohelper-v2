package convert

import (
	"fmt"
	"strconv"
)

func ToInt8(str string) int8 {
	i, _ := strconv.Atoi(str)
	return int8(i)
}

func ToUint8(str string) uint8 {
	i, _ := strconv.Atoi(str)
	return uint8(i)
}

func ToInt16(str string) int16 {
	i, _ := strconv.Atoi(str)
	return int16(i)
}

func ToUint16(str string) uint16 {
	i, _ := strconv.Atoi(str)
	return uint16(i)
}

func ToInt32(str string) int32 {
	i, _ := strconv.Atoi(str)
	return int32(i)
}

func ToUint32(str string) uint32 {
	i, _ := strconv.Atoi(str)
	return uint32(i)
}

func ToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func ToUint64(str string) uint64 {
	i, _ := strconv.ParseUint(str, 10, 64)
	return i
}

func ToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func ToUint(str string) uint {
	return uint(ToUint64(str))
}

func ToFloat32(str string) float32 {
	i, _ := strconv.ParseFloat(str, 32)
	return float32(i)
}

func ToFloat64(str string) float64 {
	i, _ := strconv.ParseFloat(str, 64)
	return i
}

func ToStr(n interface{}, args ...int) string {
	if len(args) == 0 {
		return fmt.Sprintf("%d", n)
	}

	return fmt.Sprintf("%."+fmt.Sprintf("%d", args[0])+"f", n)
}
