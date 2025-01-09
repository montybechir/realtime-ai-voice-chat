export interface Message {
    text: string;
    timestamp: Date;
}
export enum ConnectionState {
    DISCONNECTED,
    CONNECTING,
    CONNECTED,
    RECONNECTING
}

export interface AudioDelta {
    type: 'response.audio.delta';
    delta: string; // base64 encoded PCM16 24kHz audio
}

export interface TextDelta {
    type: 'response.audio_transcript.delta';
    delta: string;
}

export type WSMessage = AudioDelta | TextDelta;
