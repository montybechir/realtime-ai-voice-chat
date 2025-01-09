package ai

import (
	"encoding/json"
	"log"
	"time"

	"interviews-ai/internal/ai/types"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId   string
	AiClientId string
	Conn       *websocket.Conn
	Send       chan types.Message
	Hub        *Hub
}

// Reads from the socket connection and sends the data to the hub's handleClientRead channel
func (c *Client) ClientReadPump() {
	defer func() {
		log.Println("ClientReadPump: closing client connection")
		c.Hub.UnregisterClient <- c
		c.Conn.Close()
	}()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				c.Conn.WriteMessage(websocket.CloseGoingAway, []byte{})
			}
			break
		}

		log.Println("Client read messageType: ", messageType)

		c.Hub.HandleClientWrite <- types.Message{SenderID: c.ClientId, Payload: message, ReceiverID: c.AiClientId, Type: types.MessageType(messageType)}
	}

}

func (c *Client) ClientWritePump() {
	var ticker *time.Ticker = time.NewTicker(pingPeriod)
	defer func() {
		log.Println("ClientWritePump: closing connection")
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		// a message is sent via this specific client's send channel
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseGoingAway, []byte{})
				return
			}

			// if it wasn't a message of type text, don't write it to the user
			if message.Type != types.TextMessage {
				log.Printf("\n*******Not sending to client because message type was %v", message.Type)
				continue
			}

			// the message can potentially contain binary data
			var serverEvent ServerEvent
			if err := json.Unmarshal(message.Payload, &serverEvent); err != nil {
				log.Println("JSON Unmarshal error:", err)
				return
			}

			c.Conn.WriteMessage(websocket.TextMessage, message.Payload)

		case <-ticker.C:
			// periodically ping the client to ensure the client is listening
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
