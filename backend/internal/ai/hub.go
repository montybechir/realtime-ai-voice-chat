package ai

import (
	"interviews-ai/internal/ai/types"
	"log"
)

type Hub struct {
	Clients             map[string]*Client
	AiClients           map[string]*AIClient
	HandleClientWrite   chan types.Message
	HandleAIClientWrite chan types.Message
	RegisterClient      chan *Client
	UnregisterClient    chan *Client
	RegisterAIClient    chan *AIClient
	UnregisterAIClient  chan *AIClient
}

func NewHub() *Hub {
	return &Hub{
		Clients:             make(map[string]*Client),
		AiClients:           make(map[string]*AIClient),
		HandleClientWrite:   make(chan types.Message),
		HandleAIClientWrite: make(chan types.Message),
		RegisterClient:      make(chan *Client),
		UnregisterClient:    make(chan *Client),
		RegisterAIClient:    make(chan *AIClient),
		UnregisterAIClient:  make(chan *AIClient),
	}
}

// this handles communication between the aiClient and the serverClient
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.RegisterClient:
			hub.Clients[client.ClientId] = client
		case aiClient := <-hub.RegisterAIClient:
			hub.AiClients[aiClient.AiClientId] = aiClient
		case client := <-hub.UnregisterClient:
			if client != nil {
				_, ok := hub.Clients[client.ClientId]
				if ok {
					delete(hub.Clients, client.ClientId)
					close(client.Send)
				}
				aiClient, ok := hub.AiClients[client.AiClientId]
				if ok && aiClient != nil {
					delete(hub.AiClients, aiClient.AiClientId)
					close(aiClient.Send)
				}
			}
		case aiClient := <-hub.UnregisterAIClient:
			if aiClient != nil {
				_, ok := hub.AiClients[aiClient.AiClientId]
				if ok {
					delete(hub.AiClients, aiClient.AiClientId)
					close(aiClient.Send)
				}
				client, ok := hub.Clients[aiClient.ClientId]
				if ok && client != nil {
					delete(hub.Clients, client.ClientId)
					close(client.Send)
				}
			}
		case message := <-hub.HandleClientWrite:
			aiClient, ok := hub.AiClients[message.ReceiverID]
			if !ok {
				log.Printf("Warning: Unknown AI receiver: %v", message.ReceiverID)
				continue
			}
			log.Println("Got message from client")
			select {
			case aiClient.Send <- message:
				log.Printf("Message sent from client %s to AI %s", message.SenderID, message.ReceiverID)
			default:
				log.Printf("Failed to send message to AI %s, channel full", message.ReceiverID)
				hub.UnregisterAIClient <- aiClient
			}

		case message := <-hub.HandleAIClientWrite:
			client, ok := hub.Clients[message.ReceiverID]
			if !ok {
				log.Printf("Warning: Unknown client receiver: %v", message.ReceiverID)
				continue
			}
			select {
			case client.Send <- message:
				log.Printf("Message sent from AI %s to client %s", message.SenderID, message.ReceiverID)
			default:
				log.Printf("Failed to send message to client %s, channel full", message.ReceiverID)
				hub.UnregisterClient <- client
			}
		}

	}
}
