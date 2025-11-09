"use client";
import { useState } from "react";

export default function ChatBox({ onClose }: { onClose: () => void }) {
  const [message, setMessage] = useState("");
  const [chat, setChat] = useState<string[]>([]);

  const handleSend = () => {
    if (message.trim()) {
      setChat([...chat, message]);
      setMessage("");
    }
  };

  return (
    <div className="fixed bottom-6 right-6 bg-white w-80 h-96 shadow-lg rounded-lg flex flex-col z-50 border border-gray-300">
      <div className="flex justify-between items-center bg-blue-500 text-white p-3 rounded-t-lg">
        <span>Chat</span>
        <button onClick={onClose} className="hover:text-gray-200">✕</button>
      </div>

      <div className="flex-1 overflow-y-auto p-3 space-y-2">
        {chat.map((msg, i) => (
          <div
            key={i}
            className="self-end bg-blue-100 px-3 py-1 rounded-lg text-sm"
          >
            {msg}
          </div>
        ))}
      </div>

      <div className="p-3 flex items-center space-x-2 border-t">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          className="flex-1 border rounded px-2 py-1 text-sm"
          placeholder="Type a message..."
        />
        <button
          onClick={handleSend}
          className="bg-blue-500 text-white px-3 py-1 rounded"
        >
          Send
        </button>
      </div>
    </div>
  );
}
