package handlers

import (
	"interviews-ai/internal/ai"
	"interviews-ai/internal/ai/types"
	"interviews-ai/internal/common/config"
	"interviews-ai/internal/common/utils"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	Hub    *ai.Hub
	Config *config.Config
}

func NewWsHandler(hub *ai.Hub, config *config.Config) *WsHandler {
	return &WsHandler{
		Hub:    hub,
		Config: config,
	}
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandleWs(w, r)
}

func (h *WsHandler) HandleWs(w http.ResponseWriter, r *http.Request) {

	log.Println("Incoming websocket connection")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // bad
		},
	}

	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading client's http request to a websocket connection: %v", err)
		return
	}

	// establish a websocket connection with the AI endpoint
	aiClientConn, err := ai.CreateAIWebSocketConnection(h.Config)
	if err != nil {
		log.Fatalf("Error establishing websocket connection with AI endpoint %v", err)
		return
	}

	// get a unique identifier
	clientId := utils.GenerateID("CLI")
	aiClientId := utils.GenerateID("AI")

	client := &ai.Client{
		ClientId:   clientId,
		AiClientId: aiClientId,
		Conn:       clientConn,
		Hub:        h.Hub,
		Send:       make(chan types.Message, 1024),
	}

	aiClient := &ai.AIClient{
		ClientId:   clientId,
		AiClientId: aiClientId,
		Conn:       aiClientConn,
		Hub:        h.Hub,
		Send:       make(chan types.Message, 1024),
	}

	h.Hub.RegisterClient <- client
	h.Hub.RegisterAIClient <- aiClient

	go client.ClientReadPump()
	go client.ClientWritePump()
	go aiClient.AiClientReadPump()
	go aiClient.AiClientWritePump()

}
