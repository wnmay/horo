"use client";

import { useState } from "react";
import { ChatTextMessage } from "@/types/ws_message";

interface Props {
    connected: boolean;
    send: (data: {
        type: string;
        roomId: string;
        senderId: string;
        content: string;
    }) => void;
    roomId: string;
    senderId: string;
    username: string;
}

export default function MessagsInput ({
  connected,
  send,
  roomId,
  senderId,
  username,
}: Props) {
    const [text, setText] = useState("");

    // Handle sending a message
    const handleSendMessage = () => {
        if (!text.trim()) return;
        send({
            type: "text",
            roomId: roomId,
            senderId: senderId,
            content: text,
        });
        setText("");
    };

    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter") {
            event.preventDefault();
            handleSendMessage();
        }
    };

    return (
        <div className="flex w-full p-3 items-center gap-2 ">
            <input
                type="text"
                placeholder={connected? "type your message...":"Connecting..."}
                value={text}
                onChange={(e) => setText(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={!connected}
                className="flex-1 border rounded px-3 py-2 focus:ring-2 focus:ring-blue-400"
            />
            <button 
                onClick={handleSendMessage}
                disabled={!connected}
                className={`px-4 py-2 rounded text-white ${
                    connected ? "bg-blue-500 hover:bg-blue-600" : "bg-gray-400"
                    }`}>
                Send
            </button>
        </div>
    );
}