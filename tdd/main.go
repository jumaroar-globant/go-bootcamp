package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	errNotANumber = errors.New("the input is not an integer")
)

func isMultipleOfThree(number int) bool {
	module := number % 3

	return module == 0
}

func isMultipleOfFive(number int) bool {
	module := number % 5

	return module == 0
}

func FizzBuzz(number int) string {
	res := ""
	if isMultipleOfThree(number) {
		res += "Fizz"
	}

	if isMultipleOfFive(number) {
		res += "Buzz"
	}

	if res != "" {
		return res
	}

	return strconv.Itoa(number)
}

func main() {
	insertedNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", errNotANumber)
		return
	}

	fmt.Println(FizzBuzz(insertedNumber))
}
