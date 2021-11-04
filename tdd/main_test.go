package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFizzBuzz(t *testing.T) {
	c := require.New(t)

	output := FizzBuzz(3)
	c.Equal("Fizz", output)

	output = FizzBuzz(4)
	c.Equal("4", output)

	output = FizzBuzz(5)
	c.Equal("Buzz", output)

	output = FizzBuzz(15)
	c.Equal("FizzBuzz", output)
}
