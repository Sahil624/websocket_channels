package websocket_channels

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var logger zerolog.Logger

type ChannelLayerI interface {
	GroupSend(groupName string, data interface{})
	LeaveGroup(conn *websocket.Conn, groupName string)
	GroupAdd(conn *websocket.Conn, groupName string)
}

type ChannelsConfig struct {
	channelLayer ChannelLayerI
	Debug        bool
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

func New(config ChannelsConfig) *Channel {
	var channelLayer ChannelLayerI
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Timestamp().Str("App", "Go Channel").Logger()
	if config.channelLayer != nil {
		//channelLayer =
	} else {
		channelLayer = NewMemoryLayer()
	}
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Warn().Msg("Debug Logger Used")
	}
	return &Channel{
		channelLayer: channelLayer,
	}
}
