"use client";

import React, { useEffect, useState } from 'react'
import ChatRoomList from '@/components/ChatRoomList'
import ChatRoomMiddle, { ChatRoomProps } from '@/components/chat-middle/ChatRoomMiddle'
import { auth } from "@/firebase/firebase";
import { onAuthStateChanged } from 'firebase/auth';
import { ChatRoom } from "@/components/LeftChat";

function page() {
  const [currentChatRoom, setCurrentChatRoom] = useState<ChatRoomProps | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedRoom, setSelectedRoom] = useState<ChatRoom | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [userId, setUserId] = useState<string>("");

  // Wait for Firebase Auth to be ready
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      if (user) {
        setUserId(user.uid);
        setError(null);
      } else {
        setUserId("");
        setError("Please login to view chat rooms");
      }
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const handleRoomSelect = (room: ChatRoom) => {
    setSelectedRoom(room);
  };

  const handleRoomSelect = (room: ChatRoom) => {
    setSelectedRoom(room);
  };

  /**TODO: fetch order from right chat */

  if (loading)
    return (
      <div className="flex items-center justify-center h-screen text-gray-500">
        Loading...
      </div>
    );
  
  if (error)
    return (
      <div className="flex items-center justify-center h-screen text-red-500">
        {error}
      </div>
    );

  return (
    <div className='fixed inset-0 flex w-screen h-screen overflow-hidden'>
      <div className='w-[30%]'>
        <ChatRoomList setCurrentChatRoom={setCurrentChatRoom}/>
    <div className="flex w-full h-screen">
      <div className="w-[30%] border-r">
        <ChatRoomList onRoomSelect={handleRoomSelect} />
      </div>

      <div className="flex w-[50%] justify-center items-center h-full bg-gray-100">
        <ChatRoomMiddle room={currentChatRoom} userId={userId} orderStatus="PENDING"/>
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
