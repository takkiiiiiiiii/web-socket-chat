package qkd

import (
	"log"
	"strconv"
	"strings"
	"crypto/rand"
	"math/big"
)

var QuantumChannel []Qubit
var ClassicalChannel []int


func SampleRamdomBit(device QuantumDevice) int {
	q := device.using_qubit()
	q.Hadamard(q.state)
	result := q.Measure()
	q.Reset()
	return result
}

func PrepareMessageQubit(message int, basis int, q Qubit) {
	if message == 1 {
		q.Hadamard(q.state)
	}
	if basis == 1 {
		q.X(q.state)
	}
}

func MeasureMessageQubit(basis int, q Qubit) int {
	if basis == 1 {
		q.Hadamard(q.state)
	}
	result := q.Measure()
	q.Reset()
	return result
}

func ConvertToHex(bits []int) string {
	binStr := ""
	for _, bit := range bits {
		if bit == 1 {
			binStr += "1"
		} else {
			binStr += "0"
		}
	}

	binInt, _ := strconv.ParseInt(binStr, 2, 64)
	hexStr := strconv.FormatInt(binInt, 16)
	return hexStr
}

func GenerateHex(bits []int) string {
	var hexStr string
	var hexChunk []string
	for i := 0; i < len(bits); i += 4 {
		end := i + 4
		fourBit := bits[i:end]
		hexStr = ConvertToHex(fourBit)
		hexChunk = append(hexChunk, hexStr)
	}
	finalHex := "0x"
	finalHex += strings.Join(hexChunk, "")
	return finalHex
}

// BB84 protocol for sending a classical bit
func SendSingleBitWithBB84(alice_device QuantumDevice, bob_device QuantumDevice) [4]int {
	var info [4]int
	alice_bit := SampleRamdomBit(alice_device)
	alice_basis := SampleRamdomBit(alice_device)
	info[0] = alice_bit
	info[1] = alice_basis

	q := alice_device.using_qubit()
	PrepareMessageQubit(alice_bit, alice_basis, q)

	bob_basis := SampleRamdomBit(bob_device)
	info[2] = bob_basis


	bob_result := MeasureMessageQubit(bob_basis, q)
	info[3] = bob_result

	return info
}

func ApplyOneTimePad(message []int, key []int, index int64) ([]int, int64) {
	var random_index int64
	var diff int
	message_index := 0
	if index == -1 {
		// Encryption
		encrypted_message := make([]int, len(message))
		// message_index = 0
		if len(message) > len(key) {
			count := len(message) / len(key)
			var linked_key []int
			// over := len(message) % len(key)
			for i := 0; i < count+1; i++ {
				linked_key = append(linked_key, key...)
			}
			diff := len(linked_key) - len(message)
			n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
			if err != nil {
				log.Println("err : ", err)
			}
			random_index = n.Int64()
			for j := random_index; j < random_index + int64(len(message)); j++ {
				encrypted_message[message_index] = message[message_index] ^ linked_key[j]
				message_index++
			}
		} else if len(message) == len(key) {
			random_index = 0
			for i := 0; i < len(message); i++ {
				encrypted_message[i] = message[i] ^ key[i]
			}
		} else {
			diff = len(key) - len(message)
			n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
			if err != nil {
				log.Println("err : ", err)
			}
			random_index = n.Int64()
			for i := random_index; i < random_index + int64(len(message)); i++ {
				encrypted_message[message_index] = message[message_index] ^ key[i]
				message_index++
			}
		}
		return encrypted_message, random_index
	} else {
		// Decryption 
		decrypted_message := make([]int, len(message)) 
		// message_index = 0
		if len(message) > len(key) {
			count := len(message) / len(key)
			var linked_key []int

			for i := 0; i < count+1; i++ {
				linked_key = append(linked_key, key...)
			}
			for j := index; j < index + int64(len(message)); j++ {
				decrypted_message[message_index] = message[message_index] ^ linked_key[j]
				message_index++
			}
		} else {
			for i := index+1; i < index + int64(len(message)); i++ {
				decrypted_message[message_index] = message[message_index] ^ key[i]
				message_index++
			}
		}
		return decrypted_message, index
	}
}