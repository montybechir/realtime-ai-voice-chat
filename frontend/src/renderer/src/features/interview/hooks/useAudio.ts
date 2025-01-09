import { useRef, useCallback, useEffect } from 'react';
import useAudioPlaybackQueue from './useAudioPlaybackQueue';
import { convertInt16ToFloat32 } from '@renderer/utils/audioUtils';

interface UseAudio {
    handleAudioData: (base64Data: string) => Promise<void>;
    initAudioContext: () => AudioContext;
    audioContextRef: React.MutableRefObject<AudioContext | null>;
    isPlaying: boolean;
    isReady: boolean;
    error: Error | null;
    isPlayingRef: React.MutableRefObject<boolean>;
}

interface AudioProcessError extends Error {
    code: 'DECODE_ERROR' | 'CONTEXT_ERROR' | 'QUEUE_ERROR';
}
export default function useAudio(): UseAudio {
    const audioContextRef = useRef<AudioContext | null>(null);
    const errorRef = useRef<AudioProcessError | null>(null);

    const initAudioContext = useCallback(() => {
        try {
            if (!audioContextRef.current) {
                audioContextRef.current = new AudioContext({ sampleRate: 24000 });
            }
            return audioContextRef.current;
        } catch (error) {
            const processError = new Error(
                'Failed to initialize AudioContext'
            ) as AudioProcessError;
            processError.code = 'CONTEXT_ERROR';
            errorRef.current = processError;
            throw processError;
        }
    }, []);

    const { processAudioQueue, addToQueue, isPlaying, clearQueue, isPlayingRef } =
        useAudioPlaybackQueue({
            initAudioContext
        });

    const handleAudioData = useCallback(
        async (base64Data: string) => {
            const binaryString = atob(base64Data);
            const bytes = new Uint8Array(binaryString.length);

            for (let i = 0; i < binaryString.length; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }

            const int16Data = new Int16Array(bytes.buffer);
            const float32Data = convertInt16ToFloat32(int16Data);
            for (let i = 0; i < int16Data.length; i++) {
                float32Data[i] = int16Data[i] / 0x8000; // Convert back to Float32
            }

            addToQueue(float32Data);
        },
        [processAudioQueue]
    );

    useEffect(() => {
        return () => {
            clearQueue();
            if (audioContextRef.current?.state !== 'closed') {
                audioContextRef.current?.close().catch(console.error);
            }
        };
    }, [clearQueue]);

    return {
        handleAudioData,
        initAudioContext,
        isPlaying,
        isReady: audioContextRef.current?.state === 'running',
        error: errorRef.current,
        audioContextRef,
        isPlayingRef
    };
}
