package douyin

import "github.com/xxjwxc/public/mywebsocket"

type DouyinMsg struct {
	Ttwid         string
	LiveRoomId    string
	LiveRoomTitle string
	UserUniqueId  string
	Wss           *mywebsocket.MyWebSocket
}
type MessageInfo struct {
	Method  string
	Payload []byte
}

type douyinResp struct {
	App douyinApp `json:"app"`
}

type douyinApp struct {
	InitialState initialState `json:"initialState"`
	Odin         odin         `json:"odin"`
}

type initialState struct {
	RoomStore roomStore `json:"roomStore"`
}
type roomStore struct {
	RoomInfo roomInfo `json:"roomInfo"`
}

type roomInfo struct {
	RoomId string `json:"roomId"`
	Room   room   `json:"room"`
}

type room struct {
	Title string `json:"title"`
}

type odin struct {
	UserUniqueId string `json:"user_unique_id"`
}
