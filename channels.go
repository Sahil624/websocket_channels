package websocket_channels

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var logger zerolog.Logger

// ChannelLayerI is interface which layers have to implement to integrate with Channel library
type ChannelLayerI interface {
	GroupSend(groupName string, data interface{})
	LeaveGroup(conn *websocket.Conn, groupName string)
	GroupAdd(conn *websocket.Conn, groupName string)
}

// ChannelsConfig contains config which user can add for using third party layers
// Can set as blank config, default params are used in that case
type ChannelsConfig struct {
	channelLayer ChannelLayerI
	Debug        bool
}

// Channel denotes the channel application.
type Channel struct {
	channelLayer ChannelLayerI
}

// GroupAdd adds a connection to a group.
func (channel *Channel) GroupAdd(conn *websocket.Conn, groupName string) {
	channel.channelLayer.GroupAdd(conn, groupName)
}

// LeaveGroup removes a connection from a group.
func (channel *Channel) LeaveGroup(conn *websocket.Conn, groupName string) {
	channel.channelLayer.LeaveGroup(conn, groupName)
}

// GroupSend sends message to while group.
func (channel *Channel) GroupSend(groupName string, data interface{}) {
	channel.channelLayer.GroupSend(groupName, data)
}

// New creates a new Channels instance.
// channel := New(ChannelsConfig{})
// Need to pass channelConfig to this function.
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
