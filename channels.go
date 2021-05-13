package websocket_channels

import "github.com/gorilla/websocket"

type ChannelLayerI interface {
	GroupSend(groupName string, data interface{})
	LeaveGroup(conn *websocket.Conn, groupName string)
	GroupAdd(conn *websocket.Conn, groupName string)
}

type ChannelsConfig struct {
	channelLayer ChannelLayerI
}

type Channel struct {
	channelLayer ChannelLayerI
}

func (channel *Channel) GroupAdd(conn *websocket.Conn, groupName string) {
	channel.channelLayer.GroupAdd(conn, groupName)
}

func (channel *Channel) LeaveGroup(conn *websocket.Conn, groupName string) {
	channel.channelLayer.LeaveGroup(conn, groupName)
}
func (channel *Channel) GroupSend(groupName string, data interface{}) {
	channel.channelLayer.GroupSend(groupName, data)
}

func New(config ...ChannelsConfig) *Channel {
	return &Channel{}
}
