package main

import (
	"github.com/takkiiiiiiiii/chat/qkd"
	"github.com/gorilla/websocket"
	"time"
	"log"
	"fmt"
)

// var segment qkd.Segment

type client struct {
	//socketはwebクライアントのためのWebsocket  //WebSocketとは、WebサーバとWebブラウザの間で双方向通信できるようにする技術
	socket *websocket.Conn
	//sendはメッセージが送られるチャネル
	send chan *message
	// Qubitを送るためのチャネル
	quantumChannel chan qkd.Qubit
	// classical channel
	classicalChannel chan []int
	//roomはこのクライアントが参加しているチャットルーム
	room *room
	//userdata  ユーザーに関するデータを保持
	userData map[string]interface{}
	// 量子鍵配送用の鍵
	shareKey []int
}

//WriteMessage and ReadMessage methods to send and receive messages as a slice of bytes

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil { //ReadJSON func(v interface{}) error   message.goのmessage型をデコード
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

func (c *client) SimulateBB84(n_bit int) []int {

	var key []int
	round := 0
	for {
		if len(key) >= n_bit {
			break
		}
		round += 1

		sender_info, sender_qubit, err := qkd.CreateSingleBitWithBB84()
		if err != nil {
			log.Println(err)
		}
		c.quantumChannel <- sender_qubit

		receiver_qubit := <-c.quantumChannel
		receiver_info := qkd.ChooseBasisBobside(receiver_qubit)
		if sender_info[1] == receiver_info[0] {
			if sender_info[0] == receiver_info[1] {
				
				key = append(key, sender_info[0])
			}
		}
	}
	c.shareKey = key
	fmt.Printf("Took %d rounds to generate a %d-bit key.\n", round, n_bit)

	return key
}