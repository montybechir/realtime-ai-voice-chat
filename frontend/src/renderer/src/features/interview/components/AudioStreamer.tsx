import { useState } from "react"

const AudioStreamer: React.FC = () => {
    const [mediaChunk, setMediaChunks] = useState<Float32Array[]>([])
    const {isRecording, startRecording, stopRecording} = useMicroPhone(
        
    )
    return (

    )
}