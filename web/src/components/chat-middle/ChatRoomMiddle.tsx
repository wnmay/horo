"use client";

import { Message } from "./Message";
import { NotificationMessage } from "./NotificationMessage";
import { useEffect, useMemo, useRef, useState } from "react";
import { ChatRoom } from "../LeftChat";
import { useWebSocket } from "@/lib/ws/useWebSocket";
import { ChatMessage } from "@/types/ws-message";
import api from "@/lib/api/api-client";
import MessagsInput from "./MessageInput";
import { MessageProps } from "@/types/common-type";

interface ChatMiddleProps {
    room: ChatRoom | null;
    orderStatus: string;
    userId: string;
}

export default function ChatRoomMiddle({
    room = null,
    orderStatus,
    userId,
}:ChatMiddleProps) 
{
    const { messages, connected, joinRoom, sendMessage } = useWebSocket();
    const [historyMessage, setHistoryMessage] = useState<ChatMessage[]>([]);
    const [loading, setLoading] = useState(false);
    const chatContainerRef = useRef<HTMLDivElement | null>(null);

    // join room
    useEffect(() => {
        if (!room?.ID || !connected) return;
        const t = setTimeout(() => joinRoom(room.ID), 200);
        return () => clearTimeout(t);
    }, [room?.ID, connected, joinRoom]);
    
    // fetch old message
    useEffect(() => {
        const fetchMessages = async () => {
            if (!room) return;
            setLoading(true);

            try {
                const res = await api.get(`/api/chat/${room.ID}/messages`);

                const raw = res.data.data ?? [];
                const normalized = raw.map((m: MessageProps) => ({
                    messageId: m.ID,
                    roomId: m.RoomID,
                    senderId: m.SenderID,
                    content: m.Content,
                    type: m.Type,
                    trigger: m.Trigger,
                    createdAt: m.CreatedAt
            }));

                setHistoryMessage(normalized);
                console.log(`[Chat] Loaded ${res.data.data?.length} old messages`);
            } catch (err) {
                console.error("[Chat] Failed to load messages:", err);
            } finally {
                setLoading(false);
            }
        };

        fetchMessages();
    }, [room?.ID]);

    // const allMessages = messages;
    const allMessages = useMemo(() => {
        if(!room) return [];
        const realTime = messages.filter((m) => m.roomId === room.ID);
        return [...historyMessage, ...realTime];
    }, [historyMessage, messages, room?.ID]);

    // Auto scroll 
    useEffect(() => {
        if (!chatContainerRef.current) return;
        chatContainerRef.current.scrollTo({
        top: chatContainerRef.current.scrollHeight,
        behavior: "smooth",
        });
    }, [allMessages.length]); 

    if(!room) 
        return (<div className="w-full h-full border-r-2 border-l-2 border-gray-300"/>);

    const username = room.CustomerID === userId? room.CustomerName: room.ProphetName;
    const rec_username = room.CustomerID === userId? room.ProphetName: room.CustomerName;
    const display_coursename = (!room.courseName? "Horoscope Session":
                                room.courseName?.length <= 24? room.courseName: room.courseName?.slice(0,24)+"...");
    
    return (
        <div className="flex flex-col w-full h-full border-r-2 border-l-2 border-gray-300">
            <header className="flex flex-col w-full h-[10%] p-3 text-2xl font-bold bg-slate-50">
                {/** TODO: fetch course name */}
                🔮 {display_coursename} #{room.ID.slice(0, 8)} 
                <span className="text-lg font-light ml-10">Status: {orderStatus.toLocaleLowerCase()}</span>
            </header>

            {/* Chat session */}
            <div 
                ref={chatContainerRef}
                className="flex-1 w-full p-3 overflow-y-auto bg-white">
                {allMessages.length === 0? (
                    <p className="text-gray-400 italic text-center">
                        Welcome to a new chat.
                    </p>
                ) : (
                allMessages.map((msg, index) => 
                    msg.type === "text" ? (
                        <Message
                            key={index}
                            userId={userId}
                            senderId={msg.senderId}
                            senderName={username}
                            receiverName={rec_username}
                            content={msg.content}
                            status={"sent"}
                            createdAt={msg.createdAt}
                        />):
                        (<NotificationMessage
                            key={index}
                            msg={msg}
                        />)
                )
                )}
                {room.IsDone && 
                <div className="w-full h-6 rounded-lg bg-gray-50">
                    Chat room ended.
                </div>}
            </div>
            
            <div className="w-full h-[10%] flex justify-center items-center">
                <MessagsInput 
                    connected={connected}
                    sendMessage={sendMessage}
                    roomId={room.ID}
                    roomIsDone={room.IsDone}
                />
            </div>
        </div>
    );
}