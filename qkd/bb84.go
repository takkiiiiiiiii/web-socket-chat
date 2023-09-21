package qkd

import (
	// "log"
	"strconv"
	"strings"
	"fmt"
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

func ApplyOneTimePad(message []int, key []int, index int, env int) ([]int, int64) {
	if env == 1 {
		encrypted_message := make([]int, len(message), len(message))
		fmt.Println(len(message))
		if len(message) > len(key) {
			message_index := 0
			var n int64
			rest_message := len(message)
			key_index := 0
			count := len(message) / len(key)
			for i := 0; i < count+1; i++ {

				key_index = 0
				if i == count {
					// rest_key := int64(len(key) - rest_message + 1)
					// n, err := rand.Int(rand.Reader, big.NewInt(rest_key))
					// if err != nil {
					// 	log.Println(err)
					// }
					// key_index := n.Int64()
					// for j := message_index; j < len(message); j++{
					// 	encrypted_message[j] = message[j] ^ key[key_index]
					// 	key_index++
					// }
					key_index = 0
					for j := message_index; j < len(message); j++ {
						encrypted_message[j] = message[j] ^ key[key_index]
						key_index++
					}
				} else {
					k := message_index
					for j := message_index; j < k + len(key); j++ {
						encrypted_message[j] = message[j] ^ key[key_index]
						key_index++
						message_index++
					}
					rest_message -= len(key)
				}
			}
			return encrypted_message, n
		} else {
			fmt.Println("aaaaaa")
			for i := 0; i < len(message); i++ {
				encrypted_message[i] = message[i] ^ key[i]
			}
			return encrypted_message, 0
		}
	} else {
		decrypted_message := make([]int, len(message), len(message))
		if len(message) > len(key) {
			message_index := 0
			count := len(message) / len(key)
			for i := 0; i < count+1; i++ {
				key_index := 0
				if i == count {
					index = 0
					for j := message_index; j < len(message); j++ {
						decrypted_message[j] = message[j] ^ key[index]
						index++
					}
				} else {
					k := message_index
					for j := message_index; j < k + len(key); j++ {
						decrypted_message[j] = message[j] ^ key[key_index]
						key_index++
						message_index++
					}
				}
			}
		} else {
			fmt.Println("dddddd")
			for j := 0; j < len(message); j++ {
				decrypted_message[j] = message[j] ^ key[j]
			}
		}
		return decrypted_message, 0
	}
}
