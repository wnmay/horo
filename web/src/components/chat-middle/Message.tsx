"use client";

type statusType = "sent" | "delivered" | "read" | "failed";

interface messageProps {
    userId: string;
    senderId: string;
    senderName: string;
    content: string;
    status: statusType;
    createdAt: string;
}

export function Message({
    userId,
    senderId,
    senderName,
    content,
    status,
    createdAt
}: messageProps)
{
    const isMine = senderId === userId;
    const time = new Date(createdAt).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    });

    return (
        <div className={`w-full h-auto flex flex-col ${isMine? "items-end":"items-start"} gap-2 p-3`}>
            <div className="flex gap-2">
                <p className={`font-bold text-lg ${isMine && "order-2"}`}>{senderName}</p>
            </div>

            <div className={`flex gap-2`}>
                <div className={`h-full p-3 rounded-full text-base ${isMine? "order-2 bg-slate-200":"bg-blue-500 text-white"}`}>{content}</div>
                <div className={`flex flex-col self-end text-xs text-gray-400 ${isMine? "order-1 items-end" : "items-start"}`}>
                    <p>{status}</p>
                    <p>{time}</p>
                </div>
            </div>
        </div>
    );
}