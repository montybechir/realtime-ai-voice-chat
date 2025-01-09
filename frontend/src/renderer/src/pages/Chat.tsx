import useInterview from '@renderer/features/interview/hooks/useInterview';
import React from 'react';

export const Chat: React.FC = () => {
    const {
        messageLogRef,
        messages,
        handleSubmit,
        inputValue,
        connected,
        setInputValue,
        onBegin,
        playRecording,
        stopRecording
    } = useInterview();

    return (
        <div className="h-screen bg-gray-100 dark:bg-gray-800 flex flex-col">
            <div
                ref={messageLogRef}
                className="flex-1 bg-white dark:bg-gray-700 m-2 p-2 rounded-lg overflow-y-auto"
            >
                {messages.map((msg, idx) => (
                    <div
                        key={`${msg.timestamp.getTime()}-${idx}`}
                        className="p-2 mb-2 bg-gray-50 dark:bg-gray-600 rounded"
                    >
                        {msg.text}
                    </div>
                ))}
            </div>
            <form onSubmit={handleSubmit} className="px-2 mb-4 flex gap-2">
                <input
                    type="text"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    placeholder="Type a message..."
                    disabled={!connected}
                    className="flex-1 p-2 rounded border border-gray-300 dark:border-gray-600 
                         dark:bg-gray-700 dark:text-white disabled:opacity-50"
                />
                <button
                    type="submit"
                    disabled={!connected}
                    className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 
                         disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    Send
                </button>
            </form>

            <button onClick={onBegin}>Enable Mic & Record</button>
            <button onClick={stopRecording}>Stop Recording</button>
            {/* <button onClick={onConnect}>Connect</button>
            <button onClick={onRecord}>Start recordding</button>
           
            <button onClick={onSendRecording}>Send Recording</button>
           */}
            <button onClick={playRecording}>Play Recording</button>
        </div>
    );
};
