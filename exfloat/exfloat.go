package exfloat

import (
	"fmt"
	"strconv"
)

// 8个小数点的精度
func Float8Format(fNum float64) float64 {
	retNum,err := strconv.ParseFloat(fmt.Sprintf("%.8f", fNum), 64)
	if err != nil {
		return fNum
	} else {
		return retNum
	}
}


// 6个小数点的精度
func Float6Format(fNum float64) float64 {
	retNum,err := strconv.ParseFloat(fmt.Sprintf("%.6f", fNum), 64)
	if err != nil {
		return fNum
	} else {
		return retNum
	}
}
