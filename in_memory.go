package websocket_channels

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type group struct {
	messageChannel chan interface{}
	groupName      string
	clientMap      map[*websocket.Conn]bool
}

func (grp *group) listenChannel() {
	for {
		select {
		case message := <-grp.messageChannel:
			for client := range grp.clientMap {
				err := client.WriteJSON(message)

				if err != nil {
					fmt.Println("Error in sending message", err)
				}
			}
		}
	}
}

func (grp *group) Add(conn *websocket.Conn) {
	grp.clientMap[conn] = true
	logger.Debug().Str("Group", grp.groupName).Int("New Count", len(grp.clientMap)).Msg("New Client Added")
}

func (grp *group) Remove(conn *websocket.Conn) bool {
	delete(grp.clientMap, conn)
	logger.Debug().Str("Group", grp.groupName).Int("New Count", len(grp.clientMap)).Msg("New Client Removed")
	return len(grp.clientMap) == 0
}

func newGroup(name string) *group {
	grp := &group{
		groupName:      name,
		clientMap:      make(map[*websocket.Conn]bool),
		messageChannel: make(chan interface{}),
	}
	go grp.listenChannel()
	return grp
}

type inMemoryLayer struct {
	groupMap map[string]*group
}

func (layer *inMemoryLayer) GroupSend(groupName string, data interface{}) {
	if group := layer.groupMap[groupName]; group != nil {
		group.messageChannel <- data
	}
}

func (layer *inMemoryLayer) LeaveGroup(conn *websocket.Conn, groupName string) {
	if group, ok := layer.groupMap[groupName]; ok {
		if empty := group.Remove(conn); empty {
			delete(layer.groupMap, groupName)
		}
	}
}

func (layer *inMemoryLayer) GroupAdd(conn *websocket.Conn, groupName string) {
	if group, ok := layer.groupMap[groupName]; ok {
		group.Add(conn)
	} else {
		group = newGroup(groupName)
		layer.groupMap[groupName] = group
		group.Add(conn)
		logger.Debug().Str("Group", groupName).Int("New Group Count", len(layer.groupMap)).Msg("New Group Created")
	}
}

type MemoryLayerConfig struct {
	DebugLevel bool
}

func NewMemoryLayer() ChannelLayerI {
	var layer = new(inMemoryLayer)
	layer.groupMap = make(map[string]*group)
	return layer
}
