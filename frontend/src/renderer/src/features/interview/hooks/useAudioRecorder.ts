import { useRef, useCallback, useState, useEffect } from 'react';

interface UseAudioRecorder {
    startRecording: () => Promise<void>;
    stopRecording: () => void;
    playRecording: () => Promise<void>;
    isRecording: boolean;
    processCurrentRecording: () => Promise<Int16Array<ArrayBuffer>>;
}
export default function useAudioRecorder({
    audioContextRef,
    initAudioContext,
    isPlayingRef,
    onAudioProcess
}): UseAudioRecorder {
    const [isRecording, setIsRecording] = useState(false);
    const processorRef = useRef<ScriptProcessorNode | null>(null);
    const streamRef = useRef<MediaStream | null>(null);
    const recordedChunksRef = useRef<Float32Array[]>([]);

    const startRecording = useCallback(async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                audio: {
                    sampleRate: 24000,
                    channelCount: 1,
                    echoCancellation: true
                }
            });

            const audioContext = initAudioContext();
            const source = audioContext.createMediaStreamSource(stream);
            const processor = audioContext.createScriptProcessor(4096, 1, 1);

            processor.onaudioprocess = (event: AudioProcessingEvent) => {
                const audioData = event.inputBuffer.getChannelData(0);
                recordedChunksRef.current.push(new Float32Array(audioData));

                onAudioProcess(event);
            };

            source.connect(processor);
            processor.connect(audioContext.destination);
            streamRef.current = stream;
            setIsRecording(true);
        } catch (error) {
            console.error('Failed to start recording:', error);
        }
    }, [onAudioProcess, initAudioContext]);

    const playRecording = useCallback(async () => {
        try {
            const ctx = new AudioContext({ sampleRate: 24000 });
            await ctx.resume();

            // Concatenate recorded chunks
            const totalLength = recordedChunksRef.current.reduce(
                (acc, curr) => acc + curr.length,
                0
            );
            const combinedBuffer = new Float32Array(totalLength);

            let offset = 0;
            recordedChunksRef.current.forEach((buffer) => {
                combinedBuffer.set(buffer, offset);
                offset += buffer.length;
            });

            const audioBuffer = ctx.createBuffer(1, combinedBuffer.length, 24000);
            audioBuffer.copyToChannel(combinedBuffer, 0);

            const source = ctx.createBufferSource();
            source.buffer = audioBuffer;
            source.connect(ctx.destination);

            return new Promise<void>((resolve, reject) => {
                if (!isPlayingRef) {
                    reject(new Error('isPlayingRef is undefined'));
                    return;
                }

                source.onended = () => {
                    if (isPlayingRef) {
                        isPlayingRef.current = false;
                    }
                    resolve();
                };

                try {
                    source.start();
                    if (isPlayingRef) {
                        isPlayingRef.current = true;
                    }
                } catch (error) {
                    reject(error);
                }
            });
        } catch (error) {
            console.error('Error playing recording:', error);
            isPlayingRef.current = false;
        }
    }, [initAudioContext, isPlayingRef, recordedChunksRef]);

    const processCurrentRecording = useCallback(async (): Promise<Int16Array<ArrayBuffer>> => {
        if (!recordedChunksRef.current.length) {
            throw new Error('No recorded chunk to play');
        }

        try {
            // First, combine all Float32 chunks
            const totalLength = recordedChunksRef.current.reduce(
                (acc, curr) => acc + curr.length,
                0
            );
            const combinedFloat32 = new Float32Array(totalLength);

            let offset = 0;
            recordedChunksRef.current.forEach((buffer) => {
                combinedFloat32.set(buffer, offset);
                offset += buffer.length;
            });

            // Convert to PCM16
            const int16Buffer: Int16Array<ArrayBuffer> = new Int16Array(combinedFloat32.length);
            for (let i = 0; i < combinedFloat32.length; i++) {
                int16Buffer[i] = Math.min(1, combinedFloat32[i]) * 0x7fff;
            }

            return int16Buffer;
        } catch (error) {
            console.error('Error sending recording:', error);
            throw error;
        }
    }, []);

    const cleanup = useCallback(() => {
        try {
            if (audioContextRef.current?.state !== 'closed') {
                audioContextRef.current?.close().catch(console.error);
            }

            recordedChunksRef.current = [];

            setIsRecording(false);
        } catch (error) {
            console.error('Cleanup error:', error);
        }
    }, []);

    const stopRecording = useCallback(() => {
        try {
            if (processorRef.current) {
                processorRef.current.disconnect();
                processorRef.current = null;
            }

            if (streamRef.current) {
                streamRef.current.getTracks().forEach((track) => track.stop());
                streamRef.current = null;
            }
        } catch (error) {
            console.error('Error in stopRecording:', error);
        }
    }, [cleanup]);
    useEffect(() => {
        return cleanup;
    }, [cleanup]);

    return {
        startRecording,
        stopRecording,
        playRecording,
        isRecording,
        processCurrentRecording
    };
}
