import { useRef, useCallback } from 'react';

interface AudioPlaybackQueueProps {
    initAudioContext: () => AudioContext;
}

interface AudioPlaybackQueueState {
    playAudioBuffer: (audioBuffer: Float32Array) => Promise<void>;
    processAudioQueue: () => Promise<void>;
    clearQueue: () => void;
    isPlaying: boolean;
    isPlayingRef: React.MutableRefObject<boolean>;
    addToQueue: (audioBuffer: Float32Array) => void;
}
export default function useAudioPlaybackQueue({
    initAudioContext
}: AudioPlaybackQueueProps): AudioPlaybackQueueState {
    const isPlayingRef = useRef<boolean>(false);
    const audioQueueRef = useRef<Float32Array[]>([]);
    const sourceRef = useRef<AudioBufferSourceNode | null>(null);

    const playAudioBuffer = useCallback(async (audioBuffer: Float32Array) => {
        const ctx = initAudioContext();
        const buffer = ctx.createBuffer(1, audioBuffer.length, 24000);
        buffer.copyToChannel(audioBuffer, 0);

        const source = ctx.createBufferSource();
        sourceRef.current = source;

        source.buffer = buffer;
        source.connect(ctx.destination);

        return new Promise<void>((resolve, reject) => {
            source.onended = () => {
                isPlayingRef.current = false;
                sourceRef.current = null;
                resolve();
            };
            // source.onerror = (error) => {
            //     cleanup();
            //     reject(error);
            // };
            source.start();
            isPlayingRef.current = true;
        });
    }, []);

    const processAudioQueue = useCallback(async () => {
        if (isPlayingRef.current) return;

        while (audioQueueRef.current.length > 0) {
            const nextBuffer = audioQueueRef.current.shift();
            if (nextBuffer) {
                isPlayingRef.current = true;
                try {
                    await playAudioBuffer(nextBuffer);
                } catch (error) {
                    console.error('Error playing audio buffer:', error);
                }
                isPlayingRef.current = false;
            }
        }
    }, [playAudioBuffer]);

    const addToQueue = useCallback(
        (audioBuffer: Float32Array) => {
            audioQueueRef.current.push(audioBuffer);
            processAudioQueue();
        },
        [processAudioQueue]
    );

    const clearQueue = useCallback(() => {
        audioQueueRef.current = [];
        if (sourceRef.current) {
            sourceRef.current.stop();
            sourceRef.current.disconnect();
            sourceRef.current = null;
        }
        isPlayingRef.current = false;
    }, []);
    return {
        playAudioBuffer,
        processAudioQueue,
        isPlaying: isPlayingRef.current,
        // audioQueueRef: audioQueueRef,
        clearQueue,
        isPlayingRef,
        addToQueue
    };
}
