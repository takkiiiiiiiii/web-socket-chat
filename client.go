package main

import (
	"github.com/takkiiiiiiiii/chat/qkd"
	"github.com/gorilla/websocket"
	"time"
	"log"
	"fmt"
)


type client struct {
	//socketはwebクライアントのためのWebsocket  //WebSocketとは、WebサーバとWebブラウザの間で双方向通信できるようにする技術
	socket *websocket.Conn
	//sendはメッセージが送られるチャネル
	send chan *message
	// Qubitを送るためのチャネル
	quantumChannel chan qkd.Qubit
	// basisを報告するためのチャネル
	classicalChannel chan int
	//roomはこのクライアントが参加しているチャットルーム
	room *room
	//userdata  ユーザーに関するデータを保持
	userData map[string]interface{}
	// 量子鍵配送用の鍵
	shareKey []int
}

// WriteMessage and ReadMessage methods to send and receive messages as a slice of bytes

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil { // ReadJSON func(v interface{}) error   message.goのmessage型をデコード
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			c.room.forward <- msg //チャネルへ値を 送信
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (client *client) SimulateBB84(nBit int) {
	round := 0

	for {
		if len(client.shareKey) >= nBit {
			break
		}
		round++

		senderInfo, senderQubit, err := qkd.CreateSingleBitWithBB84()
		if err != nil {
			log.Println(err)
		}

		client.quantumChannel <- senderQubit // Qubitを送信

		client.classicalChannel <- senderInfo[1] // 受信者に basis を送信

		receiveQubit := <- client.quantumChannel


		receiverBasis := <-client.classicalChannel // 受信者から basis を受信
		receiverResult := qkd.MeasureMessageQubit(receiverBasis, receiveQubit)

		if senderInfo[1] == receiverBasis {
			if senderInfo[0] == receiverResult {
				client.shareKey = append(client.shareKey, senderInfo[0]) // 鍵を生成
			}
		}
	}

	fmt.Printf("Took %d rounds to generate a %d-bit key.\n", round, nBit)
}

func (c *client) shareKeyWithOthers() {
    // Qubitと基底情報を生成
    senderInfo, senderQubit , _ := qkd.CreateSingleBitWithBB84()

    // QuantumInfo構造体に情報を格納
    senderQuantumInfo := &QuantumInfo{
        Qubit: &senderQubit,
        Basis: senderInfo[1],
    }

	senderBit := senderInfo[0]

    // QuantumInfoを他のクライアントに送信
    c.sendQuantumInfoToOthers(senderQuantumInfo)

    // 他のクライアントからの情報を受信
    receiverQuantumInfo := c.receiveQuantumInfoFromOthers()

    // 測定を行い、鍵を生成
    if senderQuantumInfo.Basis == receiverQuantumInfo.Basis { // 基底を比較
        measuredBit := qkd.MeasureMessageQubit(receiverQuantumInfo.Basis, *receiverQuantumInfo.Qubit) // 受信者の量子情報を元にbit値を取得
        if senderBit == measuredBit { // bitが等しいかどうか
            // 鍵共有成功、鍵を生成
            c.shareKey = append(c.shareKey, senderBit) // 鍵の一部としてappend
        }
    }
}

// QuantumInfoを他のクライアントに送信する関数
func (c *client) sendQuantumInfoToOthers(info *QuantumInfo) {
    // QuantumInfoを他のクライアントに送信するためのロジックをここに追加
}

// 他のクライアントからQuantumInfoを受信する関数
func (c *client) receiveQuantumInfoFromOthers() *QuantumInfo {
    // 他のクライアントからQuantumInfoを受信するためのロジックをここに追加
    // 受信したQuantumInfoを返す
    return receivedInfo
}