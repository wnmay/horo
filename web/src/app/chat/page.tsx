"use client";

import { useEffect, useState } from 'react'
import ChatRoomList from '@/components/ChatRoomList'
import ChatRoomMiddle from '@/components/chat-middle/ChatRoomMiddle'
import { auth } from "@/firebase/firebase";
import { ChatRoom } from "@/components/LeftChat";
import RightPanel from '@/components/RightChatPanel';

function page() {
  const [role, setRole] = useState<"customer" | "prophet" | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedRoom, setSelectedRoom] = useState<ChatRoom | null>(null);
  const [userId, setUserId] = useState<string | null>(null);

  useEffect(() => {
    const unsub = auth.onAuthStateChanged(async (user) => {
      if (user) {
        try {
          const tokenResult = await user.getIdTokenResult();
          const claims = tokenResult.claims;
          const userRole =
            (claims.role as "customer" | "prophet") ?? "customer";
          setRole(userRole);
          // const uid = claims.id as string;
          const uid = user.uid;
          setUserId(uid);
        } catch (err) {
          console.error("getIdTokenResult error:", err);
          setRole("customer");
          setUserId(null);
        }
      } else {
        setRole(null);
        setUserId(null);
      }
      setLoading(false);
    });
    return () => unsub();
  }, []);

  const handleRoomSelect = (room: ChatRoom) => {
    setSelectedRoom(room);
  };

  if (loading) return <p className="p-4 text-gray-500">Loading user...</p>;
  if (!role || !userId) return <p className="p-4 text-red-500">No user signed in.</p>;

  return (
    <div className="flex h-screen self-end w-full">
      <div className="w-[30%] border-r">
        <ChatRoomList onRoomSelect={handleRoomSelect} />
      </div>

      <div className="flex w-[50%] justify-center items-center h-full bg-gray-100">
        <ChatRoomMiddle room={selectedRoom} userId={userId} orderStatus="PENDING"/>
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

export default page