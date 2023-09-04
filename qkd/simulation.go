package qkd

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"math"
)

var ket_0 = [2]float64{1, 0}

var quantumDevice QuantumDevice


func init() {
	var q Qubit
	q.State = ket_0
	quantumDevice.available_qubits = append(quantumDevice.available_qubits, q)
}

var H = [][]float64{{1 / math.Sqrt(2), 1 / math.Sqrt(2)},
	{1 / math.Sqrt(2), -1 / math.Sqrt(2)}}

var X = [][]float64{{0 / math.Sqrt(2), 1 / math.Sqrt(2)},
	{1 / math.Sqrt(2), 0 / math.Sqrt(2)}}

func NewQubit() *Qubit {
	return &Qubit{}
}

func (q *Qubit) Hadamard(State [2]float64) {
	q.State[0] = q.State[0]*H[0][0] + q.State[1]*H[1][0]
	q.State[1] = q.State[0]*H[0][1] + q.State[1]*H[1][1]
}

func (q *Qubit) X(State [2]float64) {
	q.State[0] = q.State[0]*H[0][0] + q.State[1]*H[1][0]
	q.State[1] = q.State[0]*H[0][1] + q.State[1]*H[1][1]
}

func (q *Qubit) Measure() int {
	pr0 := math.Pow(math.Abs(q.State[0]), 2)
	num := randomFloat()
	if pr0 >= num {
		return 0
	}
	return 1
}

func (q *Qubit) Reset() {
	q.State = ket_0
}

func randomFloat() float64 {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		log.Fatalln("Fatal", err)
	}
	return float64(binary.LittleEndian.Uint64(buf[:])) / (1 << 64)
}

func (qd *QuantumDevice) Allocate_qubit() (Qubit, error) {
	if len(quantumDevice.available_qubits) != 0 {
		q := quantumDevice.available_qubits[len(quantumDevice.available_qubits)-1]
		q.State = ket_0
		qd.available_qubits = quantumDevice.available_qubits[:len(quantumDevice.available_qubits)-1]
		return q, nil
	}
	return Qubit{}, errors.New("No available qubits")
}

func (qd *QuantumDevice) Deallocate_qubit(q Qubit) {
	qd.available_qubits = append(qd.available_qubits, q)
}

func (qd *QuantumDevice) Using_qubit() (Qubit, error) {
	q, err := qd.Allocate_qubit()
	if err != nil {
		return q, errors.New("No available qubits")
	}
	defer q.Reset()
	defer qd.Deallocate_qubit(q)
	return q, nil
}
