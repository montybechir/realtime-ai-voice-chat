import { useRef, useState, useCallback } from 'react';
import { ConnectionState, Message, WSMessage } from '../types';

export const useWebSocket = ({ handleAudioData }) => {
    const [messages, setMessages] = useState<Message[]>([]);
    const [connected, setConnected] = useState(false);

    const ws = useRef<WebSocket | null>(null);
    const connectionState = useRef<ConnectionState>(ConnectionState.DISCONNECTED);
    const reconnectTimeout = useRef<NodeJS.Timeout>();
    const reconnectAttempts = useRef(0);
    const MAX_RECONNECT_ATTEMPTS = 5;
    const RECONNECT_INTERVAL = 1000;
    const isCleaningUp = useRef(false);

    const WEBSOCKET_URL = import.meta.env.VITE_WEBSOCKET_URL;
    if (!WEBSOCKET_URL) {
        throw new Error('No websocket');
    }

    const handleMessage = useCallback(
        (evt: MessageEvent) => {
            try {
                const data = JSON.parse(evt.data);
                const messages = Array.isArray(data) ? data : [data];

                messages.forEach((message: WSMessage) => {
                    switch (message.type) {
                        case 'response.audio.delta':
                            handleAudioData(message.delta);
                            break;
                        case 'response.audio_transcript.delta':
                            setMessages((prev) => [
                                ...prev,
                                {
                                    text: message.delta,
                                    timestamp: new Date()
                                }
                            ]);
                            break;
                    }
                });
            } catch (error) {
                console.error('Error processing message:', error);
            }
        },
        [handleAudioData]
    );

    const connect = () => {
        if (
            connectionState.current === ConnectionState.CONNECTING ||
            connectionState.current === ConnectionState.CONNECTED //||
            //isCleaningUp.current
        ) {
            return;
        }

        connectionState.current = ConnectionState.CONNECTING;
        ws.current = new WebSocket(WEBSOCKET_URL);

        ws.current.onopen = () => {
            connectionState.current = ConnectionState.CONNECTED;
            reconnectAttempts.current = 0;
            setConnected(true);
        };

        ws.current.onclose = () => {
            if (isCleaningUp.current) return;
            connectionState.current = ConnectionState.RECONNECTING;
            setConnected(false);

            if (reconnectAttempts.current < MAX_RECONNECT_ATTEMPTS) {
                reconnectTimeout.current = setTimeout(
                    () => {
                        reconnectAttempts.current++;
                        connect();
                    },
                    RECONNECT_INTERVAL * Math.pow(2, reconnectAttempts.current)
                );
            }
        };

        ws.current.onerror = (error) => {
            console.error('WebSocket error:', error);
            setConnected(false);
        };

        ws.current.onmessage = (evt) => {
            handleMessage(evt);
        };
    };

    const sendMessage = (message: string) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            const textmessage = {
                type: 'response.create',
                response: {
                    modalities: ['audio', 'text'],
                    instructions: message
                }
            };
            ws.current.send(JSON.stringify(textmessage));
        } else {
            console.error("Couldn't send message");
        }
    };

    const sendRecording = useCallback(async (int16Buffer: Int16Array<ArrayBuffer>) => {
        try {
            // Split into chunks (e.g., 16KB chunks)
            const CHUNK_SIZE = 16 * 1024;
            for (let i = 0; i < int16Buffer.length; i += CHUNK_SIZE) {
                const chunk = int16Buffer.slice(i, i + CHUNK_SIZE);
                const chunkBase64 = btoa(
                    String.fromCharCode.apply(null, Array.from(new Uint8Array(chunk.buffer)))
                );

                const audioAppend = {
                    type: 'input_audio_buffer.append',
                    audio: chunkBase64
                };

                if (ws.current?.readyState === WebSocket.OPEN) {
                    ws.current.send(JSON.stringify(audioAppend));
                } else {
                    console.error('Could not send audio data');
                }
            }

            // Send end marker
            if (ws.current?.readyState === WebSocket.OPEN) {
                ws.current.send(
                    JSON.stringify({
                        type: 'audio.input.complete'
                    })
                );
            }
        } catch (error) {
            console.error('Error sending recording:', error);
        }
    }, []);

    return {
        ws,
        connect,
        connected,
        messages,
        sendMessage,
        sendRecording
    };
};
