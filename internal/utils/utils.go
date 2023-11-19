package utils

import (
	"errors"
	"math"
)

// This function takes in a num var of type float64 and a precision var of type int and returns
// a float64 var representing num with the precision 'precision'

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// This function takes in a str var of type string and returns an int32 type representing
// the first number contained in the string, or an error type if there is none

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

// This function converts the num var of type uin32 representing m^2
// to a float64 value representing the hectars

func MetersToHa(num uint32) float64 {
	return float64(float64(num) / 10000)
}

// This function takes in a price var of type uin32 representing a monetary value
// and an area var of type uint32 representing area in m^2 and returns a
// float64 value representing the price per ha computed from the input parameters

func PricePerHa(price uint32, area uint32) float64 {
	return ToFixed(float64(price)/float64(MetersToHa(area)), 2)
}

// This function takes in a num var of type float64 and return a int type
// represented num rounded to the closest integer value

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
