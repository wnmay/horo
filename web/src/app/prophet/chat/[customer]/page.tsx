"use client";

import { useEffect, useState } from "react";

interface Message {
  sender: string;
  text: string;
  time: string;
}

interface Chat {
  customer: string;
  messages: Message[];
}

export default function ChatPage({ params }: { params: { customer: string } }) {
  const [chats, setChats] = useState<Chat[]>([]);
  const customer = params.customer;

  useEffect(() => {
    async function fetchChats() {
      // Replace with your API call
      const res = await fetch(`/api/chats/${customer}`);
      const data = await res.json();
      setChats(data);
    }
    fetchChats();
  }, [customer]);

  return (
    <div className="p-4">
      <h1 className="text-xl font-bold mb-4">Chat with {customer}</h1>
      <div className="flex flex-col gap-3">
        {chats.map((chat, idx) =>
          chat.messages.map((msg, i) => (
            <div
              key={i}
              className={`p-2 rounded ${
                msg.sender === "customer" ? "bg-blue-100 self-start" : "bg-green-100 self-end"
              }`}
            >
              <strong>{msg.sender}:</strong> {msg.text}
              <div className="text-xs text-gray-500">{msg.time}</div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
