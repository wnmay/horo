"use client";
import { use, useEffect, useRef, useState } from "react";
import { Paperclip, Image as ImageIcon, Send } from "lucide-react";

interface Message {
  sender: "customer" | "prophet";
  text?: string;
  file?: string; // file name
  image?: string; // base64 for preview
}

interface Chat {
  customer: string;
  messages: Message[];
}

export default function ProphetChatPage({
  params,
}: {
  params: Promise<{ customer: string }>;
}) {
  const { customer } = use(params);
  const [chats, setChats] = useState<Chat[]>([
    { customer: "Alice", messages: [] },
    { customer: "Bob", messages: [] },
  ]);
  const [input, setInput] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const activeChat = chats.find(
    (chat) => chat.customer.toLowerCase().replace(/\s+/g, "-") === customer
  );

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [chats, customer]);

  const sendMessage = () => {
    if ((!input.trim() && !selectedFile && !selectedImage) || !activeChat) return;

    const newMessage: Message = { sender: "prophet" };
    if (input) newMessage.text = input;
    if (selectedFile) newMessage.file = selectedFile.name;
    if (selectedImage) newMessage.image = selectedImage;

    setChats((prev) =>
      prev.map((chat) =>
        chat.customer === activeChat.customer
          ? { ...chat, messages: [...chat.messages, newMessage] }
          : chat
      )
    );

    // reset
    setInput("");
    setSelectedFile(null);
    setSelectedImage(null);

    // Simulate auto-reply
    setTimeout(() => {
      setChats((prev) =>
        prev.map((chat) =>
          chat.customer === activeChat.customer
            ? {
                ...chat,
                messages: [
                  ...chat.messages,
                  {
                    sender: "customer",
                    text: "🙏 Got your message!",
                  },
                ],
              }
            : chat
        )
      );
    }, 1000);
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (file.type.startsWith("image/")) {
      const reader = new FileReader();
      reader.onload = (e) => setSelectedImage(e.target?.result as string);
      reader.readAsDataURL(file);
    } else {
      setSelectedFile(file);
    }
  };

  return (
    <div className="flex h-screen">
      {/* Sidebar */}
      <div className="w-1/4 border-r p-2">
        <h2 className="text-lg font-semibold mb-3 px-2">Your Customers</h2>
        {chats.map((chat) => {
          const isActive =
            chat.customer.toLowerCase().replace(/\s+/g, "-") === customer;
          return (
            <a
              key={chat.customer}
              href={`/prophet/chat/${chat.customer
                .toLowerCase()
                .replace(/\s+/g, "-")}`}
              className={`block p-3 rounded-lg mb-2 ${
                isActive ? "bg-green-100 font-semibold" : "hover:bg-gray-100"
              }`}
            >
              {chat.customer}
            </a>
          );
        })}
      </div>

      {/* Chat window */}
      <div className="flex-1 flex flex-col">
        {/* Header */}
        <div className="border-b p-4 font-semibold text-lg bg-gray-50">
          Chat with {activeChat?.customer || "Unknown"}
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto p-4 flex flex-col space-y-4">
          {activeChat?.messages.map((msg, i) => (
            <div
              key={i}
              className={`flex flex-col max-w-xs ${
                msg.sender === "customer" ? "self-start" : "self-end"
              }`}
            >
              <span className="text-xs text-gray-500 mb-1">
                {msg.sender === "prophet" ? "You" : activeChat.customer}
              </span>

              {msg.text && (
                <div
                  className={`p-2 rounded-lg mb-1 ${
                    msg.sender === "customer"
                      ? "bg-gray-200 text-black"
                      : "bg-green-500 text-white"
                  }`}
                >
                  {msg.text}
                </div>
              )}

              {msg.image && (
                <img
                  src={msg.image}
                  alt="sent"
                  className="max-w-[200px] rounded-lg border mb-1"
                />
              )}

              {msg.file && (
                <a
                  href="#"
                  className="text-sm text-blue-600 underline"
                  onClick={(e) => e.preventDefault()}
                >
                  📎 {msg.file}
                </a>
              )}
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>

        {/* Input */}
        <div className="flex items-center gap-2 border-t p-2 bg-white">
          {/* Attach button */}
          <label className="cursor-pointer text-gray-600 hover:text-blue-500">
            <input
              type="file"
              hidden
              onChange={handleFileSelect}
              accept="image/*,.pdf,.doc,.docx,.txt"
            />
            <Paperclip className="w-5 h-5" />
          </label>

          {/* Image button */}
          <label className="cursor-pointer text-gray-600 hover:text-blue-500">
            <input
              type="file"
              hidden
              onChange={handleFileSelect}
              accept="image/*"
            />
            <ImageIcon className="w-5 h-5" />
          </label>

          {/* Input box */}
          <input
            type="text"
            className="flex-1 border rounded-lg p-2"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && sendMessage()}
            placeholder="Type a message..."
          />

          {/* Send button */}
          <button
            onClick={sendMessage}
            className="bg-green-500 text-white p-2 rounded-lg hover:bg-green-600"
          >
            <Send className="w-5 h-5" />
          </button>
        </div>

        {/* Preview area */}
        {(selectedFile || selectedImage) && (
          <div className="border-t p-3 bg-gray-50 flex items-center gap-3">
            {selectedImage && (
              <img
                src={selectedImage}
                alt="preview"
                className="w-16 h-16 object-cover rounded border"
              />
            )}
            {selectedFile && (
              <p className="text-sm text-gray-700">📎 {selectedFile.name}</p>
            )}
            <button
              className="ml-auto text-red-500 text-sm"
              onClick={() => {
                setSelectedFile(null);
                setSelectedImage(null);
              }}
            >
              ✕ Cancel
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
