package govuegui

import "github.com/gorilla/websocket"

type hub struct {
	connections map[*websocket.Conn]*websocket.Conn
}

func newWebsocketHub() *hub {
	return &hub{
		connections: make(map[*websocket.Conn]*websocket.Conn),
	}
}

func (h *hub) addConnection(c *websocket.Conn) {
	h.connections[c] = c
}

func (h *hub) removeConnection(c *websocket.Conn) {
	delete(h.connections, c)
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
