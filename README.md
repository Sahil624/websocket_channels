# Websocket Channels

Websocket Channel is a library inspired by [Django Channels](https://github.com/django/channels) which sends messages over groups.
[Gorilla Connections](https://github.com/gorilla/websocket) can be added to or removed from groups and messages can be broadcasted over a WebSocket.


## Usage

```golang
package main

import (
      "github.com/Sahil624/websocket_channels"
      "github.com/gorilla/websocket"
)

channel := websocket_channels.New(websocket_channels.ChannelsConfig{
	Debug: true,
})

http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	groupName := "GROUP NAME"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Error
	}
	channel.GroupAdd(conn, groupName)
	channel.GroupSend(groupName, "Someone Entered This Group")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// ...
			}
			break
		}
		channel.GroupSend(groupName, string(message))
	}
})
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
