package qkd

import (
	"fmt"
	"strconv"
	"crypto/rand"
	"log"
	"strings"
)

type Segment struct {
    message string
    message_bit []int
    padded_message_bit  []int
}


func Generate_message_bit(message string) []int{
	var message_bit []int
	for _, as := range []byte(message) {
        s := fmt.Sprintf("%b", as)
        intSlice := StringToIntSlice(s)
        if len(intSlice) < 8 {
            intSlice = PadWithZeros(intSlice, 8)
        }
        message_bit = append(message_bit, intSlice...)
    }
	return message_bit
}

func StringToIntSlice(s string) []int {
	intSlice := make([]int, len(s))
	for i, c := range s {
		intVal, _ := strconv.Atoi(string(c))
		intSlice[i] = intVal
	}
	return intSlice
}

func PadWithZeros(bits []int, desiredBitCount int) []int {
	paddingCount := desiredBitCount - len(bits)
	paddedBits := make([]int, paddingCount)
	for i := range paddedBits {
		paddedBits[i] = 0
	}

	paddedBits = append(paddedBits, bits...)
	return paddedBits
}

func Generate_padded_message_bit(message_bit []int, key_length int) []int {
	var padded_bits []int
	if len(message_bit) < key_length {
        padded_bits =  PadWithRandomBit(message_bit, key_length)
    }
	return padded_bits
}

func PadWithRandomBit(bits []int, desiredBitCount int) []int {
    paddingCount := desiredBitCount - len(bits)
    paddedBits := make([]int, paddingCount)
    for i := range paddedBits {
        var err error
        paddedBits[i], err = RandomBit()
        if err != nil {
            log.Println("cannot generate random bit")
        }
    }
    paddedBits = append(paddedBits, bits...)
    return paddedBits
}

func RandomBit() (int, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return int(b[0]) & 1, nil
}

func Decryption_message_bit(decrypted_message_bit []int) string {
	var decrypted_message string
	for i := 0; i < len(decrypted_message_bit); i += 8 {
		end := i + 8
		demical := ConvertToDecimal(decrypted_message_bit[i:end])
		ascii := byte(demical)
		decrypted_message += string(ascii)
	}
	return decrypted_message
}

func ConvertToDecimal(bits []int) int {
	bitStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(bits)), ""), "[]")
	decimalValue, _ := strconv.ParseInt(bitStr, 2, 0)
	return int(decimalValue)
}
