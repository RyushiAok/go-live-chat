package json

import (
	"RealTime_Group_Chat/domain/entity"
	"encoding/json"
	opt "github.com/repeale/fp-go/option"
)

type Message struct {
	JsonType string         `json:"type"`
	Data     entity.Message `json:"data"`
}

func CreateMessage(message *entity.Message) *Message {
	return &Message{
		JsonType: "message",
		Data:     *message,
	}
}

func TryParseMessage(input []byte) opt.Option[Message] {
	var message Message
	// https://qiita.com/nayuneko/items/2ec20ba69804e8bf7ca3
	if err := json.Unmarshal(input, &message); err != nil {
		return opt.None[Message]()
	} else {
		return opt.Some[Message](message)
	}
}
