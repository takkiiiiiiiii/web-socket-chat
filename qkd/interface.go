package qkd

import (
	// "fmt"
)

type SimulatedQubit interface {
	Hadamard() 
	X()
	Measure() int
	Reset() 
}

type Qubit struct {
	Name 	string // ユーザー名
	State 	[2]float64
}

type SingleQubitSimulator interface {
	allocate_qubit() Qubit
	deallocate_qubit(Qubit)
	using_qubit() SimulatedQubit
}

type QuantumDevice struct {
	available_qubits []Qubit
}

