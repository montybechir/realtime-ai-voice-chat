package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"interviews-ai/internal/ai"

	"interviews-ai/internal/ai/types"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func handleWs(w http.ResponseWriter, r *http.Request, hub *ai.Hub, config *ai.Config) {

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
	aiClientConn, err := ai.CreateAIWebSocketConnection(config)
	if err != nil {
		log.Fatalf("Error establishing websocket connection with AI endpoint %v", err)
		return
	}

	// get a unique identifier
	clientId := generateConnectionID("CLI")
	aiClientId := generateConnectionID("AI")

	client := &ai.Client{
		ClientId:   clientId,
		AiClientId: aiClientId,
		Conn:       clientConn,
		Hub:        hub,
		Send:       make(chan types.Message, 1024),
	}

	aiClient := &ai.AIClient{
		ClientId:   clientId,
		AiClientId: aiClientId,
		Conn:       aiClientConn,
		Hub:        hub,
		Send:       make(chan types.Message, 1024),
	}

	hub.RegisterClient <- client
	hub.RegisterAIClient <- aiClient

	go client.ClientReadPump()
	go client.ClientWritePump()
	go aiClient.AiClientReadPump()
	go aiClient.AiClientWritePump()

}

func main() {

	config, configErr := ai.LoadConfig()
	if configErr != nil {
		log.Fatal("Error loading config: ", configErr)
	}
	hub := ai.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWs(w, r, hub, config)
	})

	log.Printf("Starting new socket server on port 5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatalln("Unexpected serve error: ", err)
	}
}

func generateConnectionID(prefix string) string {
	timestamp := time.Now().Format("20250104150405")
	uid := uuid.New().String()[:8]
	return fmt.Sprintf("%s_%s_%s", prefix, timestamp, uid)
}
