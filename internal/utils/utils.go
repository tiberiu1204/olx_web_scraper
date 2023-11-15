package utils

import (
	"errors"
	"math"
)

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func GetNumberFromString(str string) (int32, error) {
	var number int32 = 0
	var err error
	ZERO, NINE := '0', '9'
	for index := range str {
		char := str[index]
		if ZERO <= rune(char) && rune(char) <= NINE {
			number = number*10 + (rune(char) - ZERO)
		}
	}
	if number == 0 {
		err = errors.New("")
	}
	return number, err
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
