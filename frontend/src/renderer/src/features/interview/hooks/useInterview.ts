import { useState, useRef, useEffect, useCallback } from 'react';
import useAudio from './useAudio';
import useAudioRecorder from './useAudioRecorder';
import { useWebSocket } from './useWebSocket';
import { Message } from '../types';

interface UseInterview {
    playRecording: () => Promise<void>;
    stopRecording: () => void;
    messages: Message[];
    connected: boolean;
    handleSubmit: (event: React.FormEvent) => void;
    onBegin: () => void;
    messageLogRef: React.RefObject<HTMLDivElement>;
    inputValue: string;
    setInputValue: React.Dispatch<React.SetStateAction<string>>;
    isPlayingRef: React.MutableRefObject<boolean>;
}
export default function useInterview(): UseInterview {
    const [inputValue, setInputValue] = useState('');
    const messageLogRef = useRef<HTMLDivElement>(null);

    const { handleAudioData, audioContextRef, initAudioContext, isPlayingRef } = useAudio();
    const { ws, connect, connected, messages, sendMessage } = useWebSocket({
        handleAudioData
    });

    const onAudioProcess = useCallback(
        (event: AudioProcessingEvent) => {
            const audioData = event.inputBuffer.getChannelData(0);
            const int16Buffer = new Int16Array(audioData.length);

            // Convert Float32 to Int16
            for (let i = 0; i < audioData.length; i++) {
                int16Buffer[i] = audioData[i] * 0x7fff;
            }

            // Encode and send
            const chunkBase64 = btoa(
                String.fromCharCode.apply(null, Array.from(new Uint8Array(int16Buffer.buffer)))
            );

            if (ws.current?.readyState === WebSocket.OPEN) {
                ws.current.send(
                    JSON.stringify({
                        type: 'input_audio_buffer.append',
                        audio: chunkBase64
                    })
                );
            }
        },
        [ws]
    );
    const { startRecording, playRecording, stopRecording, isRecording } = useAudioRecorder({
        audioContextRef,
        initAudioContext,
        isPlayingRef,
        onAudioProcess
    });

    const handleSubmit = (event: React.FormEvent) => {
        event.preventDefault();
        if (!inputValue.trim() || !connected) return;

        sendMessage(inputValue.trim());
        setInputValue('');
    };

    const onBegin = useCallback(async () => {
        try {
            // 1. Connect WebSocket
            await startRecording();
            //connect();
        } catch (error) {
            console.error('Failed to start interview:', error);
        }
    }, [startRecording]);

    useEffect(() => {
        if (!connected && isRecording) {
            connect();
        }
    }, [connected, isRecording]);

    useEffect(() => {
        if (messageLogRef.current) {
            messageLogRef.current.scrollTop = messageLogRef.current.scrollHeight;
        }
    }, [messages]);

    useEffect(() => {
        return () => {
            stopRecording();
        };
    }, [stopRecording]);

    return {
        playRecording,
        stopRecording,
        messages,
        connected,
        handleSubmit,
        onBegin,
        messageLogRef,
        inputValue,
        setInputValue,
        isPlayingRef
    };
}
