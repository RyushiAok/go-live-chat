package ws

import (
	"RealTime_Group_Chat/domain/entity"
	ws_json "RealTime_Group_Chat/presen/ws/json"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

//	func tryParseJson[T any](input []byte) (result T, err error) {
//		var tmp T
//		if err := json.Unmarshal(input, &tmp); err != nil {
//			return result, fmt.Errorf("tryParseJson")
//		} else {
//			return tmp, nil
//		}
//	}

func tryParseMessage(input []byte) (result *ws_json.Message, err error) {
	var tmp ws_json.Message
	if err := json.Unmarshal(input, &tmp); err != nil { //tmp.Dataの型はみられてないぽい
		return result, fmt.Errorf("tryParseJson")
	} else if tmp.JsonType != "message" {
		return result, fmt.Errorf("tryParseJson")
	} else {
		return &tmp, nil
	}
}

func tryParseJoinNewMember(input []byte) (result *ws_json.JoinNewMember, err error) {
	var tmp ws_json.JoinNewMember
	if err := json.Unmarshal(input, &tmp); err != nil { //tmp.Dataの型はみられてないぽい
		return result, fmt.Errorf("tryParseJson")
	} else if tmp.JsonType != "join_new_member" {
		return result, fmt.Errorf("tryParseJson")
	} else {
		return &tmp, nil
	}
}

// Goでカリー化するには，いかつい書き方をする必要があるので，明確なメリットがない場合は用いない方が良さそう．
// カリー化することで，意図せずハンドラに関係のない値を投げることがなくなる.（同じ型ならコンパイルが通って気づかない）
// また，引数のバケツリレーで呼び出し側のネストが深くなることを避けられる．

func receiveMessage(handler func(msg *ws_json.Message) (err error)) func(input []byte) (err error) {
	return func(input []byte) (err error) {
		msg, err := tryParseMessage(input) // tryParseJson[json2.Message](input)
		if err != nil {
			return fmt.Errorf("receiveMessage (tryParseJson)")
		}
		return handler(msg)
	}
}

func receiveJoinNewMember(handler func(msg *ws_json.JoinNewMember) (err error)) func(input []byte) (err error) {
	return func(input []byte) (err error) {
		msg, err := tryParseJoinNewMember(input) // tryParseJson[json2.JoinNewMember](input)
		if err != nil {
			return fmt.Errorf("receiveJoinNewMember (tryParseJson)")
		}
		return handler(msg)
	}
}

func HandleMessage(h *Hub, s *Subscription) func(msg *ws_json.Message) (err error) {
	return func(msg *ws_json.Message) (err error) {
		resp := *ws_json.CreateMessage(&entity.Message{Name: msg.Data.Name, Message: "receiveMessage : " + msg.Data.Message})
		data, err := json.Marshal(resp)
		if err == nil {
			h.Broadcast <- Message{Room: s.Room, Data: data} //ws.message{ }
			return nil
		}
		return fmt.Errorf("receiveMessage (handler)")
	}
}

func HandleJoinNewMember(h *Hub, s *Subscription) func(msg *ws_json.JoinNewMember) (err error) {
	return func(msg *ws_json.JoinNewMember) (err error) {
		resp := *ws_json.CreateJoinNewMember(&entity.Member{Name: "receiveJoinNewMember : " + msg.Data.Name})
		data, err := json.Marshal(resp)
		if err == nil {
			h.Broadcast <- Message{Room: s.Room, Data: data} //ws.message{ }
			return nil
		}
		return fmt.Errorf("receiveMessage (handler)")
	}
}

// 受信したメッセージを基にハンドラを再定義するので（常にレスポンスが同じではない），部分適用したハンドラを使いまわすことはない．
// msgのアドレスが常に同じなら， (* msg, handlers) => () => error という定義はなくはないかもだけどナシ（きもい）
// カリー化による旨味はない． (handlers) => (msg) => option
func handleWSMessage(input []byte, fs []func(input []byte) (err error)) (err error) {
	for i := range fs {
		if err := fs[i](input); err == nil {
			return nil
		}
	}
	return fmt.Errorf("handleWSMessage (No Match)")
}

// 受け取ったメッセージに応じてイベント(ルームにメッセージをブロードキャストする)を発火する
func (s *Subscription) readPump(h *Hub) {
	c := *s.conn
	defer func() {
		h.unregister <- *s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		if err := handleWSMessage(msg, []func(input []byte) (err error){
			// ここにまとめて全部定義するのはいかつすぎ
			receiveMessage(HandleMessage(h, s)),
			receiveJoinNewMember(HandleJoinNewMember(h, s)),
		}); err != nil {
			h.Broadcast <- Message{s.Room, []byte("Unexpected Message")}
			fmt.Println(err)
		}
	}
}

func (c *connection) write(messageType int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, payload)
}

func (s *Subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServeWs(h *Hub) func(w http.ResponseWriter, r *http.Request, roomId string) {
	// コネクションをHubに登録する
	return func(w http.ResponseWriter, r *http.Request, roomId string) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		c := &connection{send: make(chan []byte, 256), ws: ws}
		s := &Subscription{c, roomId}
		fmt.Print(roomId)
		h.register <- *s
		// タスクを登録する
		go s.writePump()
		go s.readPump(h)
	}
}
