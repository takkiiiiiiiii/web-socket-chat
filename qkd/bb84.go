package qkd

import (
	"log"
	"fmt"
	"strconv"
	"strings"
)

func SampleRamdomBit(device QuantumDevice) int {
	q, err := device.Using_qubit()
	if err != nil {
		log.Println(err)
	}
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
	alice_message := SampleRamdomBit(alice_device)
	alice_basis := SampleRamdomBit(alice_device)
	info[0] = alice_message
	info[1] = alice_basis

	bob_basis := SampleRamdomBit(bob_device)
	info[2] = bob_basis

	q, err := alice_device.Using_qubit()
	if err != nil {
		log.Println(err)
	}
	PrepareMessageQubit(alice_message, alice_basis, q)

	// Qubit sending

	bob_result := MeasureMessageQubit(bob_basis, q)
	info[3] = bob_result

	return info
}

func createSingleBitWithBB84() ([2]int, Qubit, error) {
	var result [2]int
	var my_device QuantumDevice
	my_message := SampleRamdomBit(my_device)
	my_basis := SampleRamdomBit(my_device)
	result[0] = my_message
	result[1] = my_basis

	q, err := my_device.Using_qubit()
	if err != nil {
		log.Println(err)
	}
	PrepareMessageQubit(my_message, my_basis, q)


	return result, q, err
}

func chooseBasisBobside(q Qubit) [2]int {
	var receiver_device QuantumDevice
	var receiver_info [2]int
	receiver_basis := SampleRamdomBit(receiver_device)
	receiver_result := MeasureMessageQubit(receiver_basis, q)
	receiver_info[0] = receiver_basis
	receiver_info[1] = receiver_result
	return receiver_info
}


func SimulateBB84(n_bit int) []int {
	var alice_device QuantumDevice
	var bob_device QuantumDevice

	var key []int
	round := 0
	for {
		if len(key) >= n_bit {
			break
		}
		round += 1
		result := SendSingleBitWithBB84(alice_device, bob_device)
		alice_message := result[0]
		alice_basis := result[1]
		bob_basis := result[2]
		bob_result := result[3]

		if alice_basis == bob_basis {
			if alice_message == bob_result {
				key = append(key, alice_message)
			}
		}
	}
	fmt.Printf("Took %d rounds to generate a %d-bit key.\n", round, n_bit)

	return key
}

func ApplyOneTimePad(message []int, key []int) []int {
	encrypted_message := make([]int, len(message))
	for i := 0; i < len(message); i++ {
		encrypted_message[i] = message[i] ^ key[i]
	}
	return encrypted_message
}
