"use client";

import MessagsInput from "./MessageInput";
import { Message } from "./Message";
import { NotificationMessage } from "./NotificationMessage";
import { useWebSocket } from "@/lib/ws/useWebSocket";
import { useEffect, useState } from "react";

type orderStatusType = "PENDING" | "CONFIRMED" | "CANCELLED" | "COMPLETED";

type NotiType = "create_order" | "order_done";

export interface ChatRoomProps {
	ID: string; 
	ProphetID: string;
	CustomerID: string;
	CourseID: string; 
	CreatedAt: string;
	LastMessage: string; 
	IsDone: boolean;
	ProphetName: string;
	CustomerName: string;
}

interface ChatMiddleProps {
    room: ChatRoomProps | null;
    orderStatus: orderStatusType;
    userId: string;
}
export default function ChatRoomMiddle({
    room = null,
    orderStatus,
    userId,
}:ChatMiddleProps) 
{
    const { messages, send, connected } = useWebSocket();
    
    if(!room)
        return (<div className="w-full h-full border-r-2 border-l-2 border-gray-300"/>);

    const username = room.CustomerID === userId? room.CustomerName: room.ProphetName;

    // log for debugging
    useEffect(() => {
        if (messages.length > 0)
            console.log("💬 New WS message:", messages.at(-1));
    }, [messages]);
    
    const filtered = messages.filter((m) => m.roomId === room.ID);
    
    return (
        <div className="flex flex-col w-full h-full border-r-2 border-l-2 border-gray-300">
            <header className="flex flex-col w-full h-[10%] p-3 text-2xl font-bold bg-slate-50">
                {/** TODO: fetch course name */}
                🔮 Horoscope Session #{room.ID} 
                <span className="text-lg font-light ml-10">Status: {orderStatus.toLocaleLowerCase()}</span>
            </header>

            {/* Chat session */}
            <div className="flex-1 w-full p-3 overflow-y-auto bg-white">
                {filtered.length === 0? (
                    <p className="text-gray-400 italic text-center">
                        Welcome to a new chat.
                    </p>
                ) : (
                filtered.map((msg, index) => 
                    msg.type === "text" ? (
                        <Message
                            key={index}
                            userId={userId}
                            senderId={msg.senderId}
                            senderName={username}
                            content={msg.content}
                            status={"sent"}
                            createdAt={msg.createdAt}
                        />):
                        (<NotificationMessage
                            key={index}
                            msg={msg}
                        />)
                ))}
            </div>
            
            <div className="w-full h-[10%] flex justify-center items-center">
                <MessagsInput 
                    connected={connected}
                    send={send}
                    roomId={room.ID}
                    senderId={userId} 
                    username={username}
                />
            </div>
        </div>
    );
}