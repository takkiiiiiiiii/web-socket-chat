package main

import (
	"fmt"
	"github.com/takkiiiiiiiii/chat/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"log"
	"net/http"
)


// チャネルはバッファとして使える
type room struct {
	//forwardは他のクライアントに転送するためのメッセージを保持するためのチャネル
	forward chan *message
	//joinはチャットルームに参加しようとしているクライアントのためのチャネル (そのクライアントを保持)
	join chan *client
	// 量子情報を保持するためのチャネル
	quantum chan *message
	//leaveはチャットルームから退室しようとしているクライアントのためのチャネル (そのクライアントを保持)
	leave chan *client
	//roomに在室しているすべてのクライアントが保持されている
	clients_exist map[*client]bool
	// この部屋にいるクライアントの記録
	clients []*client
	//tracerはチャットルーム上で行われた操作のログを受け取ります
	tracer trace.Tracer //traceパッケージのTrace型(interface)
	avatar Avatar       //アバター情報の取得
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join: //<-channel 構文  チャネルから値を受信
			//参加
			r.clients_exist[client] = true
			r.clients = append(r.clients, client)
			fmt.Println("新しいクライアントが参加しました")
			if len(r.clients) == 2 {
				key := client.SimulateBB84(96)
				for _, c := range r.clients {
					c.shareKey = key
				}
			}
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			//退室
			r.removeClient(client) // クライアントをスライスから削除
			delete(r.clients_exist, client) //map rooo型のclientsからclientを削除
			close(client.send)
			r.tracer.Trace("クライアントは退室しました")
		case msg := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", msg.Message)
			//すべてのクライアントにメッセージを送信
			for client := range r.clients_exist {
				select {
				case client.send <- msg:
					//メッセージを送信
					r.tracer.Trace(" -- クライアントに送信しました")
				default:
					//送信に失敗
					r.removeClient(client) // クライアントをスライスから削除
					delete(r.clients_exist, client)
					close(client.send)
					r.tracer.Trace(" -- 送信失敗しました。クライアントをクリーンアップします。")
				}
			}
		}
	}
}

func (r *room) removeClient(clientToRemove *client) {
    var newClients []*client
    for _, client := range r.clients {
        if client != clientToRemove {
            newClients = append(newClients, client)
        }
    }
    r.clients = newClients
}

//すぐに利用できるチャットルームを生成して返す
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client, 2),
		leave:   make(chan *client),
		clients_exist: make(map[*client]bool, 2),
		clients: make([]*client, 0),
		tracer:  trace.Off(), // trace構造体とともに定義   trace.Off() 戻り値 *niltrace  newRoom生成したら　traceパッケージのOffメソッドも実行される
		avatar:  avatar,
	}
}

const (
	socketBufferSize  = 2048
	messageBufferSize = 1024
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil) // UpgraderのUpgradeはHTTP通信からWebSocket通信に更新してくれる
	if err != nil {
		log.Fatal("ServeHttp:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました: ", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value), // MustFromBase64の戻り値 map[string]interface{}  エンコードされたクッキーの値をマップのオブジェクトへ復元	
	}

	
	r.join <- client
	
	defer func() { r.leave <- client }()
	go client.write() // 他のクライアントにメッセージを送信
	client.read() // クライアントからのメッセージを待機
}

