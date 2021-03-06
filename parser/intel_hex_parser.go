package main

import (
	"fmt"
	"strconv"
)

func NumberOfBytes(line string) uint8 {
	slice := line[1:3]

	value, err := strconv.ParseInt(slice, 16, 64)
	if err != nil {
		fmt.Printf("Conversion failed: %s\n", err)
		panic(".")
	}

	return uint8(value)

}

func StartingAddress(line string) uint16 {
	slice := line[3:7]

	value, err := strconv.ParseInt(slice, 16, 64)
	if err != nil {
		fmt.Printf("Conversion failed: %s\n", err)
		panic(".")
	}

	return uint16(value)
}

func Record(line string) uint8 {
	slice := line[7:9]

	value, err := strconv.ParseInt(slice, 16, 64)
	if err != nil {
		fmt.Printf("Conversion failed: %s\n", err)
		panic(".")
	}

	return uint8(value)
}

func Payload(line string) [32]uint8 {
	var result [32]uint8
	n := 2 * NumberOfBytes(line)
	slice := line[9 : 9+n]

	for i := 0; i < len(slice); i += 2 {
		hex := slice[i : i+2]
		value, err := strconv.ParseInt(hex, 16, 64)
		result[i/2] = uint8(value)

		if err != nil {
			fmt.Printf("Conversion failed: %s\n", err)
			panic(".")
		}
	}

	return result
}

func CRC(line string) uint8 {
	n := 2 * NumberOfBytes(line)
	slice := line[9+n : 9+n+2]

	value, err := strconv.ParseInt(slice, 16, 64)
	if err != nil {
		fmt.Printf("Conversion failed: %s\n", err)
		panic(".")
	}

	return uint8(int8(value))
}

func IsCRCValid(line string) bool {
	n := NumberOfBytes(line)

	if n == 0 {
		return true
	}

	a := StartingAddress(line)
	r := Record(line)
	p := Payload(line)
	expectedCRC := CRC(line)

	computedCRC := n + uint8(a>>8) + uint8(a) + r
	for _, e := range p {
		computedCRC = computedCRC + e
	}

	return TwosComplement(computedCRC) == expectedCRC
}

func IsFileValid(file []string) bool {
	for _, line := range file {

		actual := IsCRCValid(line)
		if actual == false {
			return false
		}

	}
	return true
}

func TotalNumberOfBytes(file []string) int {
	result := 0
	for _, line := range file {
		result = result + int(NumberOfBytes(line))
	}
	return result
}

func TwosComplement(in uint8) uint8 {
	return ^in + 1
}
