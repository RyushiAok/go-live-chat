package handler

import (
	"RealTime_Group_Chat/domain/entity"
	"RealTime_Group_Chat/presen/ws"
	json2 "RealTime_Group_Chat/presen/ws/json"
	"encoding/json"
	"fmt"
)

func HandleMessage(h ws.Hub, s ws.Subscription) func(msg json2.Message) (err error) {
	return func(msg json2.Message) (err error) {
		resp := *json2.CreateMessage(&entity.Message{Name: msg.Data.Name, Message: "receiveMessage : " + msg.Data.Message})
		data, err := json.Marshal(resp)
		if err == nil {
			h.Broadcast <- ws.Message{Room: s.Room, Data: data} //ws.message{ }
			return nil
		}
		return fmt.Errorf("receiveMessage (handler)")
	}
}
