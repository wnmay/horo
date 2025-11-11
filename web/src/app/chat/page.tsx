"use client";

import React, { useEffect, useState } from "react";
import { auth } from "@/firebase/firebase";
import ChatRoomList from "@/components/ChatRoomList";
import RightPanel from "@/components/RightChatPanel";
import { ChatRoom } from "@/components/LeftChat";

export default function Page() {
  const [role, setRole] = useState<"customer" | "prophet" | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedRoom, setSelectedRoom] = useState<ChatRoom | null>(null);

  useEffect(() => {
    const unsub = auth.onAuthStateChanged(async (user) => {
      if (user) {
        try {
          const tokenResult = await user.getIdTokenResult();
          const claims = tokenResult.claims;
          const userRole =
            (claims.role as "customer" | "prophet") ?? "customer";
          setRole(userRole);
        } catch (err) {
          console.error("getIdTokenResult error:", err);
          setRole("customer");
        }
      } else {
        setRole(null);
      }
      setLoading(false);
    });
    return () => unsub();
  }, []);

  const handleRoomSelect = (room: ChatRoom) => {
    setSelectedRoom(room);
  };

  if (loading) return <p className="p-4 text-gray-500">Loading user...</p>;
  if (!role) return <p className="p-4 text-red-500">No user signed in.</p>;

  return (
    <div className="flex w-full h-screen">
      <div className="w-[30%] border-r">
        <ChatRoomList onRoomSelect={handleRoomSelect} />
      </div>
      <div className="flex-1">
        {selectedRoom ? (
          <RightPanel
            roomId={selectedRoom.ID}
            role={role}
            courseId={selectedRoom.CourseID}
          />
        ) : (
          <div className="flex items-center justify-center h-full text-gray-500">
            Select a chat room to start messaging
          </div>
        )}
      </div>
    </div>
  );
}
