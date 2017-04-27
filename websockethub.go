package govuegui

import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var wshub = newHub()

type hub struct {
	connections map[*websocket.Conn]*websocket.Conn
}

func newHub() *hub {
	return &hub{
		connections: make(map[*websocket.Conn]*websocket.Conn),
	}
}

func (h *hub) addConnection(c *websocket.Conn) {
	h.connections[c] = c
}

func (h *hub) writeJSON(v interface{}) error {
	for c := range h.connections {
		err := c.WriteJSON(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *hub) writeMessage(messageType int, data []byte) error {
	for c := range h.connections {
		err := c.WriteMessage(messageType, data)
		if err != nil {
			return err
		}
	}
	return nil
}
