package main

import (
	"github.com/gorilla/websocket"
	"github.com/takkiiiiiiiii/chat/qkd"
	"time"
	// "log"
	"fmt"
)

type client struct {
	//socketはwebクライアントのためのWebsocket  //WebSocketとは、WebサーバとWebブラウザの間で双方向通信できるようにする技術
	socket *websocket.Conn
	//sendはメッセージが送られるチャネル
	send chan *message
	// Qubitを送るためのチャネル
	quantumChannel chan QuantumInfo
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
		var padded_message_bit []int
		if err := c.socket.ReadJSON(&msg); err == nil { // ReadJSON func(v interface{}) error   message.goのmessage型をデコード
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			message_bit := generate_message_bit(msg.Message)
			msg.MessageBit = message_bit
			fmt.Println(c.userData["name"].(string) + "'s key :" , c.shareKey)
			if len(message_bit) < len(c.shareKey) {
				padded_message_bit = generate_padded_message_bit(message_bit, len(c.shareKey))	
				msg.PaddedMessageBit = padded_message_bit
			}
			msg.EncryptedMessage = qkd.ApplyOneTimePad(padded_message_bit, c.shareKey)
			fmt.Println("encrepted_message :" , msg.EncryptedMessage)
			fmt.Println("sending message...")
			fmt.Println("----------------------------------------")
			c.room.forward <- msg //チャネルへ値を 送信
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		decrypted_message_bit := qkd.ApplyOneTimePad(msg.EncryptedMessage, c.shareKey)
		padded_len := len(msg.PaddedMessageBit) - len(msg.MessageBit)
		decrypted_message := decryption_message_bit(decrypted_message_bit[padded_len:])
		fmt.Println(c.userData["name"].(string) + " " + decrypted_message)
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

// func (client *client) SimulateBB84(nBit int) {
// 	round := 0
// 	for {
// 		if len(client.shareKey) >= nBit {
// 			break
// 		}
// 		round++
// 		client.shareKeyWithOthers()
// 	}
// 	close(client.room.quantum)
// 	fmt.Printf("Took %d rounds to generate a %d-bit key.\n", round, nBit)
// }

// func (c *client) shareKeyWithOthers() {
// 	// Qubitと基底情報を生成
// 	senderInfo, senderQubit, _ := qkd.CreateSingleBitWithBB84()

// 	// QuantumInfo構造体に情報を格納
// 	senderQuantumInfo := &QuantumInfo{
// 		Qubit: &senderQubit,
// 		Basis: senderInfo[1],
// 	}

// 	senderBit := senderInfo[0]

// 	// Qubitを他のクライアントに送信
// 	go c.sendQuantumInfoToOthers(*senderQuantumInfo)

//     // 他のクライアントからの情報を受信
// 	receiverQuantumInfo := <- c.receiveQuantumInfoFromOthers()
// 	fmt.Println(receiverQuantumInfo)

// 	// 測定を行い、鍵を生成
// 	if senderQuantumInfo.Basis == receiverQuantumInfo.Basis { // 基底を比較
// 		measuredBit := qkd.MeasureMessageQubit(receiverQuantumInfo.Basis, *receiverQuantumInfo.Qubit) // 受信者の量子情報を元にbit値を取得
// 		if senderBit == measuredBit {                                                                 // bitが等しいかどうか
// 			// 鍵共有成功、鍵を生成
// 			c.shareKey = append(c.shareKey, senderBit) // 鍵の一部としてappend
// 		}
// 	}
// }

// // QuantumInfoを他のクライアントに送信する関数
// func (c *client) sendQuantumInfoToOthers(info QuantumInfo) {
//     c.quantumChannel <- info
// }

// // 他のクライアントからQuantumInfoを非同期に受信
// func (c *client) receiveQuantumInfoFromOthers() <-chan *QuantumInfo {
//     ch := make(chan *QuantumInfo) // チャネルを作成

//     go func() {
//         defer close(ch) // ゴルーチンが終了したらチャネルをクローズ

//         for receivedMessage := range c.room.quantum {
//             // 受信した量子情報をチャネルに送信
//             ch <- &receivedMessage.QuantumInfo
//         }
//     }()

//     return ch
// }

func (c *client) SimulateBB84(n_bit int) []int {
	var alice_device qkd.QuantumDevice
	var bob_device qkd.QuantumDevice

	var key []int
	round := 0
	for {
		if len(key) >= n_bit {
			break
		}
		round += 1
		result := qkd.SendSingleBitWithBB84(alice_device, bob_device)
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