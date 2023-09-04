package main

import (
	"github.com/takkiiiiiiiii/chat/qkd"
	"time"
)

//messageは1つのメッセージを表す
type message struct {
	Name      string    //ユーザー名
	Message   string    //contents
	When      time.Time //送信された時刻
	QuantumInfo QuantumInfo
	AvatarURL string
}


type QuantumInfo struct {
	Qubit *qkd.Qubit
	Basis int
}