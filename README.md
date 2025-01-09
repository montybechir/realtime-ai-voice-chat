# Go-React-Electron Realtime AI Audio Chat

This is an application built using a GoLang backend and Electron + React + Vite + front-end. It's designed to help users prepare for upcoming interviews using a real-time multi-modal AI. Engage in realtime audio conversations or type messages to the AI, which responds with both transcripts and audio to facilitate seamless interaction.

## Features

- **Realtime Audio Conversations:** Speak with the AI, receive immediate audio and text responses.
- **Text Messaging:** Type questions, receive immediate answers.
- **Multi-Platform Support:** Available as an app on macOS, Windows, and Linux, and browsers.
- **AI-Powered:** Expert guidance for technical interview preparation

## Technologies

- **Backend:** Go (Golang)
- **Frontend:** Electron, React, TypeScript, Vite
- **AI:** Azure OpenAI GPT-4o-Realtime-Preview (2024-10-01)

## Prerequisites

- Node.js v20.x+
- Go 1.22+
- npm 10.x+

## Quick Start

### Clone the Repository

```bash
git clone https://github.com/montybechir/realtime-ai-voice-chat.git
cd realtime-ai-voice-chat
```

### Frontend Setup

```bash
cd frontend
npm install
```

### Configure Backend

```bash
cd backend

# install go dependencies
go mod download

# copy environment config
cp .env.example .env
```

## Development

### Environment Setup

```bash
#/backend/.env

AZURE_OPENAI_API_KEY=yourkey
AZURE_OPENAI_ENDPOINT="wss://<your openai azure deployment's endpoint>.openai.azure.com/openai/realtime?api-version=2024-10-01-preview&deployment=gpt-4o-realtime-preview"
```

If you don't have access to Azure OpenAI and prefer to use the standard OpenAI endpoints, replace the above with your OpenAI keys and endpoint. Ensure you update your AI WebSocket connection headers in aiClient.go:

```bash

                Authorization: `Bearer <your api key>`,
                "OpenAI-Beta": "realtime=v1"
```

## Running the Application

### Start the Frontend:

```bash
cd frontend
npm run dev
```

### Start the AI service:

In a separate terminal, run the following:

```bash
cd backend
go run cmd/ai-service/main.go
```

## Using the Application

Desktop: Launches automatically with npm run dev
Browser: Visit http://localhost:5173

## Running Tests

To be implemented.

## Roadmap

<input disabled="" type="checkbox"> User Authentication
<input disabled="" type="checkbox"> Testing
<input disabled="" type="checkbox"> Database Integration
<input disabled="" type="checkbox"> CI/CD Pipeline
<input disabled="" type="checkbox"> Customizable Enhanced AI Prompts
<input disabled="" type="checkbox"> Session Management
