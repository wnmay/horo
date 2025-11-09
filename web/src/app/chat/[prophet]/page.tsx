"use client";
import { use } from "react";
import { useState, useRef, useEffect } from "react";
import { Paperclip, Image as ImageIcon, Send } from "lucide-react";

interface Message {
  sender: "me" | "prophet";
  text: string;
  attachment?: string | null; // ✅ Allow null or undefined
}

interface Chat {
  prophet: string;
  messages: Message[];
}

export default function ChatPage({ params }: { params: Promise<{ prophet: string }> }) {
  const { prophet } = use(params);
  const [chats, setChats] = useState<Chat[]>([
    { prophet: "master-flook", messages: [] },
    { prophet: "prophet-lyra", messages: [] },
  ]);

  const [input, setInput] = useState("");
  const [attachment, setAttachment] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chats, prophet]);

  const activeChat = chats.find(
    (chat) => chat.prophet.toLowerCase().replace(/\s+/g, "-") === prophet
  );

  const sendMessage = () => {
    if ((!input.trim() && !attachment) || !activeChat) return;

    setChats((prev) =>
      prev.map((chat) =>
        chat.prophet === activeChat.prophet
          ? {
              ...chat,
              messages: [
                ...chat.messages,
                { sender: "me", text: input || "(attachment)", attachment },
              ],
            }
          : chat
      )
    );

    setInput("");
    setAttachment(null);

    // Simulate prophet reply
    setTimeout(() => {
      setChats((prev) =>
        prev.map((chat) =>
          chat.prophet === activeChat.prophet
            ? {
                ...chat,
                messages: [
                  ...chat.messages,
                  { sender: "prophet", text: "🤖 Got your message!" },
                ],
              }
            : chat
        )
      );
    }, 1000);
  };

  const handleAttachment = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) setAttachment(file.name); // For mock display only
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
        <div className="border-b p-4 font-semibold text-lg bg-gray-50">
          {activeChat?.prophet}
        </div>

        <div className="flex-1 overflow-y-auto p-4 flex flex-col space-y-4">
          {activeChat?.messages.map((msg, idx) => (
            <div
              key={idx}
              className={`flex flex-col max-w-xs ${
                msg.sender === "prophet" ? "self-start" : "self-end"
              }`}
            >
              <span className="text-xs font-semibold mb-1 text-gray-500">
                {msg.sender === "me" ? "You" : activeChat.prophet}
              </span>

              <div
                className={`p-2 rounded-lg ${
                  msg.sender === "prophet"
                    ? "bg-gray-200 text-black"
                    : "bg-blue-500 text-white"
                }`}
              >
                {msg.text}
                {msg.attachment && (
                  <div className="mt-1 text-sm italic underline">
                    📎 {msg.attachment}
                  </div>
                )}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef}></div>
        </div>

        {/* Input bar with icons */}
        <div className="flex items-center border-t p-2 space-x-2">
          <label className="cursor-pointer">
            <ImageIcon className="w-5 h-5 text-gray-500 hover:text-blue-500" />
            <input type="file" accept="image/*" className="hidden" onChange={handleAttachment} />
          </label>
          <label className="cursor-pointer">
            <Paperclip className="w-5 h-5 text-gray-500 hover:text-blue-500" />
            <input type="file" className="hidden" onChange={handleAttachment} />
          </label>

          <input
            type="text"
            className="flex-1 border rounded-lg p-2"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && sendMessage()}
            placeholder="Type a message..."
          />

          <button
            onClick={sendMessage}
            className="bg-blue-500 text-white p-2 rounded-lg hover:bg-blue-600"
          >
            <Send className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}
