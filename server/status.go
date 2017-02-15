package server

import (
	"github.com/PurpurProject/elytra/jsonutil"
)

type ServerStatus struct {
	Version     StatusVersion       `json:"version"`
	Players     StatusPlayers       `json:"players"`
	Description jsonutil.ChatObject `json:"description"`
	Favicon     string              `json:"favicon,omitempty"`
}

type StatusVersion struct {
	Name     string `json:"name"`
	Protocol int32  `json:"protocol"`
}

type StatusPlayers struct {
	MaxPlayers    int32 `json:"max"`
	OnlinePlayers int32 `json:"online"`
}

func CreateStatusObject() *ServerStatus {
	s := new(ServerStatus)

	s.Version = StatusVersion{Name: minecraftVersionName, Protocol: minecraftProtocolVersion}
	s.Players = StatusPlayers{MaxPlayers: 20, OnlinePlayers: 0}
	exampleExtra := make([]jsonutil.ChatObject, 1)
	exampleExtra[0] = jsonutil.ChatObject{Text: "And goodbye...", Bold: true, Italic: true, Color: "red"}
	s.Description = jsonutil.ChatObject{Text: "Hello, World! ", Bold: true, Color: "blue", Extra: exampleExtra}
	return s
}
