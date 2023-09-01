package qkd

import (
	"fmt"
	"strconv"
	"strings"
)

func sampleRamdomBit(device QuantumDevice) int {
	q = device.using_qubit()
	q.Hadamard(q.state)
	result := q.Measure()
	q.Reset()
	return result
}

func prepareMessageQubit(message int, basis int, q Qubit) {
	if message == 1 {
		q.Hadamard(q.state)
	}
	if basis == 1 {
		q.X(q.state)
	}
}

func measureMessageQubit(basis int, q Qubit) int {
	if basis == 1 {
		q.Hadamard(q.state)
	}
	result := q.Measure()
	q.Reset()
	return result
}

func convertToHex(bits []int) string {
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

func generateHex(bits []int) string {
	var hexStr string
	var hexChunk []string
	for i := 0; i < len(bits); i += 4 {
		end := i + 4
		fourBit := bits[i:end]
		hexStr = convertToHex(fourBit)
		hexChunk = append(hexChunk, hexStr)
	}
	finalHex := "0x"
	finalHex += strings.Join(hexChunk, "")
	return finalHex
}

// BB84 protocol for sending a classical bit
func sendSingleBitWithBB84(alice_device QuantumDevice, bob_device QuantumDevice) [4]int {
	var info [4]int
	alice_message := sampleRamdomBit(alice_device)
	alice_basis := sampleRamdomBit(alice_device)
	info[0] = alice_message
	info[1] = alice_basis

	bob_basis := sampleRamdomBit(bob_device)
	info[2] = bob_basis

	q := alice_device.using_qubit()
	prepareMessageQubit(alice_message, alice_basis, q)

	// Qubit sending

	bob_result := measureMessageQubit(bob_basis, q)
	info[3] = bob_result

	return info
}

func simulateBB84(n_bit int) []int {
	var alice_device QuantumDevice
	var bob_device QuantumDevice

	var key []int
	round := 0
	for {
		if len(key) >= n_bit {
			break
		}
		round += 1
		result := sendSingleBitWithBB84(alice_device, bob_device)
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

func applyOneTimePad(message []int, key []int) []int {
	encrypted_message := make([]int, len(message))
	for i := 0; i < len(message); i++ {
		encrypted_message[i] = message[i] ^ key[i]
	}
	return encrypted_message
}
