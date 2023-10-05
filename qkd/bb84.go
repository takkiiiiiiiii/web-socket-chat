package qkd

import (
	// "log"
	"strconv"
	"strings"
	// "crypto/rand"
	// "math/big"
)

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

func ApplyOneTimePad(message []int, key []int, env int) []int {
	var message_index int
	if env == 1 {
		encrypted_message := make([]int, len(message))
		message_index = 0
		if len(message) > len(key) {
			count := len(message) / len(key)
			over := len(message) % len(key)
			for i := 0; i < count+1; i++ {
				if i != count {
					for j := 0; j < len(key); j++ {
						encrypted_message[message_index] = message[message_index] ^ key[j]
						message_index++
					}
				} else {
					if over != 0 {
						for j := 0; j < over; j++ {
							encrypted_message[message_index] = message[message_index] ^ key[j]
							message_index++
						}
						break
					}
				} 
			}
		} else {
			for i := 0; i < len(message); i++ {
				encrypted_message[i] = message[i] ^ key[i]
			}
		}
		return encrypted_message
	} else {
		decrypted_message := make([]int, len(message)) 
		message_index = 0
		if len(message) > len(key) {
			count := len(message) / len(key)
			over := len(message) % len(key)
			for i := 0; i < count+1; i++ {
				if i != count {
					for j := 0; j < len(key); j++ {
						decrypted_message[message_index] = message[message_index] ^ key[j]
						message_index++
					}
				} else {
					if over != 0 {
						for j := 0; j < over; j++ {
							decrypted_message[message_index] = message[message_index] ^ key[j]
							message_index++
						}
					}
				}
			}
		} else {
			for i := 0; i < len(message); i++ {
				decrypted_message[i] = message[i] ^ key[i]
			}
		}
		return decrypted_message
	}
}