package qkd

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"math"
)

var ket_0 = [2]float64{1, 0}

var q Qubit
var quantumDevice QuantumDevice


func init() {
	q.state = ket_0
	quantumDevice.available_qubits = append(quantumDevice.available_qubits, q)
}

var H = [][]float64{{1 / math.Sqrt(2), 1 / math.Sqrt(2)},
	{1 / math.Sqrt(2), -1 / math.Sqrt(2)}}

var X = [][]float64{{0 / math.Sqrt(2), 1 / math.Sqrt(2)},
	{1 / math.Sqrt(2), 0 / math.Sqrt(2)}}

func NewQubit() *Qubit {
	return &Qubit{}
}

func (q *Qubit) Hadamard(state [2]float64) {
	q.state[0] = q.state[0]*H[0][0] + q.state[1]*H[1][0]
	q.state[1] = q.state[0]*H[0][1] + q.state[1]*H[1][1]
}

func (q *Qubit) X(state [2]float64) {
	q.state[0] = q.state[0]*H[0][0] + q.state[1]*H[1][0]
	q.state[1] = q.state[0]*H[0][1] + q.state[1]*H[1][1]
}

func (q *Qubit) Measure() int {
	pr0 := math.Pow(math.Abs(q.state[0]), 2)
	num := randomFloat()
	if pr0 >= num {
		return 0
	}
	return 1
}

func (q *Qubit) Reset() {
	q.state = ket_0
}

func randomFloat() float64 {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		log.Fatalln("Fatal", err)
	}
	return float64(binary.LittleEndian.Uint64(buf[:])) / (1 << 64)
}

func (qd *QuantumDevice) allocate_qubit() (Qubit, error) {
	if len(quantumDevice.available_qubits) != 0 {
		q := quantumDevice.available_qubits[len(quantumDevice.available_qubits)-1]
		q.state = ket_0
		qd.available_qubits = quantumDevice.available_qubits[:len(quantumDevice.available_qubits)-1]
		return q, nil
	}
	return Qubit{}, errors.New("No available qubits")
}

func (qd *QuantumDevice) deallocate_qubit(q Qubit) {
	qd.available_qubits = append(qd.available_qubits, q)
}

func (qd *QuantumDevice) using_qubit() Qubit {
	q, err := qd.allocate_qubit()
	if err != nil {
		errors.New("No available qubits")
	}
	defer q.Reset()
	defer qd.deallocate_qubit(q)
	return q
}
