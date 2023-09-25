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
		var padded_message_bit []int
		if err := c.socket.ReadJSON(&msg); err == nil { // ReadJSON func(v interface{}) error   message.goのmessage型をデコード
			msg.When = time.Now().Format(time.DateTime)
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			message_bit := generate_message_bit(msg.Message)
			msg.MessageBit = message_bit
			fmt.Println(c.userData["name"].(string) + "'s key :" , c.shareKey)
			if len(message_bit) < len(c.shareKey) {
				padded_message_bit = generate_padded_message_bit(message_bit, len(c.shareKey))	// 鍵の長さ以下の場合、メッセージの長さを鍵と同じ長さにする
				msg.PaddedMessageBit = padded_message_bit // パディングしたビット値（スライス）
				msg.EncryptedMessage = qkd.ApplyOneTimePad(padded_message_bit, c.shareKey, 1)
			} else {
				msg.EncryptedMessage = qkd.ApplyOneTimePad(message_bit, c.shareKey, 1)
			}
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
		decrypted_message_bit := qkd.ApplyOneTimePad(msg.EncryptedMessage, c.shareKey, 0)
		if len(msg.MessageBit) < len(c.shareKey) {
			padded_len := len(msg.PaddedMessageBit) - len(msg.MessageBit)
			decrypted_message := decryption_message_bit(decrypted_message_bit[padded_len:])
			fmt.Println(c.userData["name"].(string) + " " + decrypted_message + "\n")
		} else {
			decrypted_message := decryption_message_bit(decrypted_message_bit)
			fmt.Println(c.userData["name"].(string) + " " + decrypted_message + "\n")
		}
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

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
