"use client";
import { use } from "react";
import { useState, useRef, useEffect } from "react";

interface Message {
  sender: "me" | "prophet";
  text: string;
}

interface Chat {
  prophet: string;
  messages: Message[];
}

export default function ChatPage({ params }: { params: Promise<{ prophet: string }> }) {
  const { prophet } = use(params); // unwrap promise

  const [chats, setChats] = useState<Chat[]>([
    { prophet: "master-flook", messages: [] },
    { prophet: "prophet-lyra", messages: [] },
  ]);

  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chats, prophet]);

  const activeChat = chats.find(
    (chat) => chat.prophet.toLowerCase().replace(/\s+/g, "-") === prophet
  );

  const sendMessage = () => {
    if (!input.trim() || !activeChat) return;

    // Add my message
    setChats((prev) =>
      prev.map((chat) =>
        chat.prophet === activeChat.prophet
          ? { ...chat, messages: [...chat.messages, { sender: "me", text: input }] }
          : chat
      )
    );
    setInput("");

    // Simulate prophet reply
    setTimeout(() => {
      setChats((prev) =>
        prev.map((chat) =>
          chat.prophet === activeChat.prophet
            ? { ...chat, messages: [...chat.messages, { sender: "prophet", text: "🤖 Auto-reply" }] }
            : chat
        )
      );
    }, 1000);
  };

  return (
    <div className="flex h-screen">
      {/* Sidebar */}
      <div className="w-1/4 border-r p-2">
        {chats.map((chat) => {
          const isActive = chat.prophet.toLowerCase().replace(/\s+/g, "-") === prophet;
          return (
            <a
              key={chat.prophet}
              href={`/chat/${chat.prophet.toLowerCase().replace(/\s+/g, "-")}`}
              className={`block p-3 rounded-lg mb-2 ${
                isActive ? "bg-indigo-100 font-semibold" : "hover:bg-gray-100"
              }`}
            >
              {chat.prophet}
            </a>
          );
        })}
      </div>

      {/* Chat window */}
      <div className="flex-1 flex flex-col">
        {/* Chat header */}
        <div className="border-b p-4 font-semibold text-lg bg-gray-50">
          {activeChat?.prophet}
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 flex flex-col space-y-4">
          {activeChat?.messages.map((msg, idx) => (
            <div
              key={idx}
              className={`flex flex-col max-w-xs ${
                msg.sender === "prophet" ? "self-start" : "self-end"
              }`}
            >
              {/* Sender name */}
              <span className="text-xs font-semibold mb-1 text-gray-500">
                {msg.sender === "me" ? "You" : activeChat.prophet}
              </span>

              {/* Message bubble */}
              <div
                className={`p-2 rounded-lg ${
                  msg.sender === "prophet"
                    ? "bg-gray-200 text-black"
                    : "bg-blue-500 text-white"
                }`}
              >
                {msg.text}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef}></div>
        </div>

        {/* Input box */}
        <div className="flex border-t p-2">
          <input
            type="text"
            className="flex-1 border rounded-l-lg p-2"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && sendMessage()}
            placeholder="Type a message..."
          />
          <button
            onClick={sendMessage}
            className="bg-blue-500 text-white p-2 rounded-r-lg"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}
