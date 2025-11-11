// mock-messages.ts
type Status = "sent" | "delivered" | "read" | "failed";
type messageType = "text" | "notification";

export const mockMessages: Array<{
  id: string;
  senderId: string;
  roomId: string;
  senderName: string;
  receiverName: string;
  content: string;
  status: Status;
  type: messageType;
  createdAt: string;
}> = [
  {
    id: crypto.randomUUID(),
    senderId: "customer-1",            // ← ฝั่งเรา
    roomId: "room1",
    senderName: "Prim",
    receiverName: "Tungmay",
    content: "สวัสดีครับ ✨",
    status: "read",
    type: "text",
    createdAt: new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  },
  {
    id: crypto.randomUUID(),
    senderId: "prophet-9",             // ← ฝั่งเค้า
    roomId: "room1",
    senderName: "Tungmay",
    receiverName: "Prim",
    content: "สวัสดีค่ะ พร้อมเริ่มเลยค่ะ",
    status: "read",
    type: "text",
    createdAt: new Date(Date.now() + 60_000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  },
  {
    id: crypto.randomUUID(),
    senderId: "customer-1",
    roomId: "room1",
    senderName: "Prim",
    receiverName: "Tungmay",
    content: "โอเคครับ เดี๋ยวผมส่งวันเกิดให้นะ",
    status: "delivered",
    type: "text",
    createdAt: new Date(Date.now() + 120_000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  },
  {
    id: crypto.randomUUID(),
    senderId: "prophet-9",
    roomId: "room1",
    senderName: "Tungmay",
    receiverName: "Prim",
    content: "รับทราบค่ะ 🙌",
    status: "sent",
    type: "text",
    createdAt: new Date(Date.now() + 180_000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  },
  {
    id: crypto.randomUUID(),
    senderId: "system",
    roomId: "room1",
    senderName: "",
    receiverName: "",
    content: "create_order",
    status: "sent",
    type: "notification",
    createdAt: new Date(Date.now() + 200_000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  },
  {
    id: crypto.randomUUID(),
    senderId: "system",
    roomId: "room1",
    senderName: "",
    receiverName: "",
    content: "order_done",
    status: "sent",
    type: "notification",
    createdAt: new Date(Date.now() + 240_000).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", hour12: false }),
  }
];
