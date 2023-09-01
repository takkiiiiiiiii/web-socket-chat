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

func Qkd(message string) {
	fmt.Println("Generating a 96-bit key by simulating BB84...")
	key := simulateBB84(2048)
	hex_key := generateHex(key)
	fmt.Println("Got key :                             "+ hex_key + "\n")

	var padded_message_bit []int
	
	message_bit := generate_message_bit(message)
	if len(message_bit) < len(key) {
		padded_message_bit = generate_padded_message_bit(message_bit, len(key))	
	}
	seg := Segment{
        message : message,
        message_bit : message_bit,
		padded_message_bit : padded_message_bit,
    }

	hex_message := generateHex(seg.padded_message_bit)
	fmt.Println("Using key to send secret message:     "+ hex_message + "\n")


	encrypted_message := applyOneTimePad(seg.padded_message_bit, key)
	hex_encrypted_message := generateHex(encrypted_message)
	fmt.Println("Encrypted message :                   "+ hex_encrypted_message + "\n")

	decrypted_message_bit := applyOneTimePad(encrypted_message, key)
	padded_len := len(seg.padded_message_bit) - len(seg.message_bit)
	
	decrypted_message := decryption_message_bit(decrypted_message_bit[padded_len:])
	fmt.Println("received message :      "+ decrypted_message + "\n")
}


func generate_message_bit(message string) []int{
	var message_bit []int
	for _, as := range []byte(message) {
        s := fmt.Sprintf("%b", as)
        intSlice := stringToIntSlice(s)
        if len(intSlice) < 8 {
            intSlice = padWithZeros(intSlice, 8)
        }
        message_bit = append(message_bit, intSlice...)
    }
	return message_bit
}

func stringToIntSlice(s string) []int {
	intSlice := make([]int, len(s))
	for i, c := range s {
		intVal, _ := strconv.Atoi(string(c))
		intSlice[i] = intVal
	}
	return intSlice
}

func padWithZeros(bits []int, desiredBitCount int) []int {
	paddingCount := desiredBitCount - len(bits)
	paddedBits := make([]int, paddingCount)
	for i := range paddedBits {
		paddedBits[i] = 0
	}

	paddedBits = append(paddedBits, bits...)
	return paddedBits
}

func generate_padded_message_bit(message_bit []int, key_length int) []int {
	var padded_bits []int
	if len(message_bit) < key_length {
        padded_bits =  padWithRandomBit(message_bit, key_length)
    }
	return padded_bits
}

func padWithRandomBit(bits []int, desiredBitCount int) []int {
    paddingCount := desiredBitCount - len(bits)
    paddedBits := make([]int, paddingCount)
    for i := range paddedBits {
        var err error
        paddedBits[i], err = randomBit()
        if err != nil {
            log.Println("cannot generate random bit")
        }
    }
    paddedBits = append(paddedBits, bits...)
    return paddedBits
}

func randomBit() (int, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return int(b[0]) & 1, nil
}

func decryption_message_bit(decrypted_message_bit []int) string {
	var decrypted_message string
	for i := 0; i < len(decrypted_message_bit); i += 8 {
		end := i + 8
		demical := convertToDecimal(decrypted_message_bit[i:end])
		ascii := byte(demical)
		decrypted_message += string(ascii)
	}
	return decrypted_message
}

func convertToDecimal(bits []int) int {
	bitStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(bits)), ""), "[]")
	decimalValue, _ := strconv.ParseInt(bitStr, 2, 0)
	return int(decimalValue)
}
