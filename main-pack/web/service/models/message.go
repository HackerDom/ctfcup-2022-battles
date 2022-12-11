package models

import (
	"encoding/json"
	"github.com/golang/glog"
)

const (
	CreateUserMessage = 1
	WriteNoteMessage  = 2
)

type Message struct {
	Data []byte `json:"data"`
	Type int    `json:"type"`
}

func (m Message) ToString() string {
	data, err := json.Marshal(m)
	if err != nil {
		glog.Error(err)
	}

	return string(data)
}

type CreateUser struct {
	Credentials Credentials `json:"credentials"`
}

func (c CreateUser) ToMessage() Message {
	data, err := json.Marshal(c)
	if err != nil {
		glog.Error(err)
	}

	return Message{
		Type: CreateUserMessage,
		Data: data,
	}
}

type WriteNote struct {
	Note
	Dir string `json:"dir"`
}

func (wc WriteNote) ToMessage() Message {
	data, err := json.Marshal(wc)
	if err != nil {
		glog.Error(err)
	}

	return Message{
		Type: WriteNoteMessage,
		Data: data,
	}
}
