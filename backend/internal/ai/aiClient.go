package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"interviews-ai/internal/ai/templates"
	"interviews-ai/internal/ai/types"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

// ServerEvent represents the structure of events exchanged with the server.
type ServerEvent struct {
	Type     string                 `json:"type"`
	Response map[string]interface{} `json:"response,omitempty"`
	Session  map[string]interface{} `json:"session,omitempty"`
	Delta    string                 `json:"delta,omitempty"`
}

// SessionUpdateEvent represents the session.update event structure.
type SessionUpdateEvent struct {
	Type    string                 `json:"type"`
	Session map[string]interface{} `json:"session,omitempty"`
}

// ResponseCreateEvent represents the response.create event structure.
type ResponseCreateEvent struct {
	Type     string                 `json:"type"`
	Response map[string]interface{} `json:"response,omitempty"`
}

type InputAudioBufferEvent struct {
	Type  string                 `json:"type"`
	Audio map[string]interface{} `json:"response,omitempty"`
}

type InputAudioBufferAppend struct {
	EventID string `json:"event_id,omitempty"`
	Type    string `json:"type"`
	Audio   string `json:"audio"`
}

type AIClient struct {
	AiClientId string
	ClientId   string
	Conn       *websocket.Conn
	Send       chan types.Message
	Hub        *Hub
}

// Constants for message types.
const (
	MsgTypeSessionUpdate                = "session.update"
	MsgTypeAudioBufferAppend            = "input_audio_buffer.append"
	MsgTypeAudioBufferCommit            = "input_audio_buffer.commit"
	MsgTypeResponseCreate               = "response.create"
	MsgTypeResponseDone                 = "response.done"
	MsgTypeResponseError                = "error"
	MsgTypeResponseAudioDelta           = "response.audio.delta"
	MsgTypeResponseAudioTranscriptDelta = "response.audio_transcript.delta"
	MsgTypeAudioTranscriptDelta         = "response.audio_transcript.delta"
	MsgTypeResponseContentPartAdded     = "response.content_part.added"

	// WebSocket timing constants
	writeWait  = 10 * time.Second
	pingPeriod = (writeWait * 9) / 10
)

type Modality string

const (
	AudioModality Modality = "audio"
	TextModality  Modality = "text"
)

type IncomingMessage struct {
	Type     string                 `json:"type"`
	Audio    string                 `json:"audio,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
	Response struct {
		Modalities   []Modality `json:"modalities"`
		Instructions string     `json:"instructions"`
	} `json:"response,omitempty"`
}

type AIMessageType string

const (
	AudioBufferAppend   AIMessageType = "input_audio_buffer.append"
	AudioBufferComplete AIMessageType = "input_audio_buffer.complete"
	ResponseCreate      AIMessageType = "response.create"
)

type Config struct {
	APIKey   string
	Endpoint string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables.")
	}

	apiKey := os.Getenv("AZURE_OPENAI_API_KEY")
	endpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")

	if apiKey == "" || endpoint == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return &Config{
		APIKey:   apiKey,
		Endpoint: endpoint,
	}, nil
}

// createAIWebSocketConnection establishes a WebSocket connection to Azure OpenAI's Realtime API.
func CreateAIWebSocketConnection(config *Config) (*websocket.Conn, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Proceeding with environment variables.")
	}

	// Define WebSocket dialer
	dialer := websocket.DefaultDialer
	header := make(http.Header)

	header.Set("api-key", config.APIKey) // Use Authorization header

	// Connect to WebSocket server
	conn, _, err := dialer.Dial(config.Endpoint, header)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}

	// Update the initial session to our desired task
	sessionUpdate := SessionUpdateEvent{
		Type: "session.update",
		Session: map[string]interface{}{
			"modalities":          []string{"audio", "text"},
			"instructions":        templates.InterviewInstructions,
			"temperature":         0.8,
			"voice":               "alloy",
			"input_audio_format":  "pcm16",
			"output_audio_format": "pcm16",
			"turn_detection": map[string]interface{}{
				"type": "server_vad",
			},
		},
	}
	initialData, err := json.Marshal(sessionUpdate)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, initialData); err != nil {
		return nil, fmt.Errorf("write initial message error: %v", err)
	}

	return conn, nil
}

func SendSessionUpdate(c *AIClient) {
	// Send initial event (session.update)
	sessionUpdate := SessionUpdateEvent{
		Type: "session.update",
		Session: map[string]interface{}{
			"commit":          true,
			"cancel_previous": true,
			"instructions":    "Help me prepare for my upcoming Growth Engineering interview",
			"modalities":      []string{"audio", "text"},
		},
	}
	data, err := json.Marshal(sessionUpdate)
	if err != nil {
		log.Printf("Failed to marshal response.create event: %v", err)
		return
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to send response.create event: %v", err)
		return
	}

	log.Println("Sent response.create event to server.")
}

// sendResponseCreate sends a response.create event to the server.
func SendResponseCreate(c *AIClient) {
	responseCreate := ResponseCreateEvent{
		Type: "response.create",
		Response: map[string]interface{}{
			"commit":          true,
			"cancel_previous": true,
			"instructions":    "Help me prepare for my upcoming Growth Engineering interview",
			"modalities":      []string{"audio", "text"},
		},
	}
	data, err := json.Marshal(responseCreate)
	if err != nil {
		log.Printf("Failed to marshal response.create event: %v", err)
		return
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to send response.create event: %v", err)
		return
	}

	log.Println("Sent response.create event to server.")
}

// handleAIResponse processes incoming server events.
func handleAIResponse(c *AIClient, message []byte) {
	var event ServerEvent
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("handleAIResponse JSON Parse error: %v. Message: %s", err, string(message))
		return
	}

	eventType := event.Type
	log.Printf("******Client event %v", eventType)

	switch eventType {
	case "session.created":
		log.Println("Session created successfully.")
		SendSessionUpdate(c)
	case "session.updated":
		//handleAudioDelta(c, event)
	case MsgTypeResponseCreate:
		log.Println("Response creation initiated.")
	case MsgTypeResponseDone:
		handleAudioDone(c, event)
	case MsgTypeResponseError:
		errMsg, exists := event.Response["error"].(string)
		if exists {
			log.Printf("Received error from server: %v", errMsg)
		} else {
			log.Printf("Received error from server: <nil>")
		}
	case MsgTypeResponseAudioDelta:
		//
	case "response.created":
		log.Println("Response created successfully.")
	case "response.output_item.added":
		log.Println("Response output item added.")
	case "conversation.item.created":
		log.Println("Conversation item created.")
	case "response.audio_transcript.delta":
		log.Println("AI event: response.audio_transcript.delta")
	default:
		log.Printf("Unknown AI client event type: %s", eventType)
	}
}

// handleAudioDone handles the response.audio.done event.
func handleAudioDone(c *AIClient, event ServerEvent) {
	log.Println("Handling response.audio.done event")
}

// handleAudioDelta handles the response.audio.delta event.
func handleAudioDelta(c *AIClient, event ServerEvent) {
	log.Printf("Played audio chunk of size: %d bytes")
}

// aiClientReadPump listens for incoming messages from the AI WebSocket connection.
func (c *AIClient) AiClientReadPump() {
	defer func() {
		log.Println("AiClientReadPump: closing connection.")
		c.Hub.UnregisterAIClient <- c
		c.Conn.Close()
	}()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Println("Unexpected websocket close err: ", err)
			}
			break
		}

		switch messageType {
		case websocket.TextMessage:
			log.Println("Received text message from AI")
			handleAIResponse(c, message)
		case websocket.BinaryMessage:
			log.Println("Received binary (audio) message from AI")
			handleAIResponse(c, message)
		default:
			log.Printf("Unknown message type: %d", messageType)
			continue
		}

		log.Printf("aiClientReadPump messageType: %v", messageType)
		// write to the hub
		c.Hub.HandleAIClientWrite <- types.Message{
			SenderID:   c.AiClientId,
			Payload:    message,
			ReceiverID: c.ClientId,
			Type:       types.TextMessage,
		}
	}
}

// aiClientWritePump sends outgoing messages to the WebSocket connection.
func (c *AIClient) AiClientWritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("AiClientWritePump: closing connection.")
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Parse incoming message
			var incomingMsg IncomingMessage
			if err := json.Unmarshal(message.Payload, &incomingMsg); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			switch incomingMsg.Type {
			case "input_audio_buffer.append":
				// Forward audio to AI
				audioMessage := InputAudioBufferAppend{
					Type:  "input_audio_buffer.append",
					Audio: incomingMsg.Audio, // base64-encoded
				}
				jsonData, err := json.Marshal(audioMessage)
				if err != nil {
					log.Printf("Error marshalling audio message: %v", err)
					continue
				}
				if err := c.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("Error writing audio to AI websocket: %v", err)
					return
				}

			case "response.create":
				responseCreate := ResponseCreateEvent{
					Type: "response.create",
					Response: map[string]interface{}{
						"modalities":   []string{"audio", "text"},
						"instructions": incomingMsg.Response.Instructions,
						"commit":       true,
					},
				}
				jsonData, err := json.Marshal(responseCreate)
				if err != nil {
					log.Printf("Error marshalling text esponse create: %v", err)
					continue
				}
				if err := c.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("Error writing text to AI websocket: %v", err)
					return
				}
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
