'use client';

import React, { useEffect, useState } from 'react';
import LeftChat, { ChatRoom } from './LeftChat';
import api from "@/lib/api/api-client";
import { auth } from '@/firebase/firebase';
import { onAuthStateChanged } from 'firebase/auth';

interface ChatRoomListProps {
  onRoomSelect?: (room: ChatRoom) => void;
}

const ChatRoomList: React.FC<ChatRoomListProps> = ({ onRoomSelect }) => {
  const [chatRooms, setChatRooms] = useState<ChatRoom[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [authReady, setAuthReady] = useState(false);
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null);

  // Wait for Firebase Auth to be ready
  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setAuthReady(true);
      if (!user) {
        setLoading(false);
        setError('Please login to view chat rooms');
      }
    });
    return () => unsubscribe();
  }, []);

  const fetchChatRooms = async () => {
    try {
      setLoading(true);
      
      const response = await api.get('/api/chat/user/rooms');
      console.log('API Response:', response.data); // Debug log
      
      const rooms = response.data?.data || [];
      
      const roomsWithCourseNames = await Promise.all(
        rooms.map(async (room: ChatRoom) => {
          try {
            const courseResponse = await api.get(`/api/courses/${room.CourseID}`);
            return {
              ...room,
              courseName: courseResponse.data.data.coursename || 'Unknown Course'
            };
          } catch (error) {
            console.error(`Error fetching course name for ${room.CourseID}:`, error);
            return {
              ...room,
              courseName: 'Unknown Course'
            };
          }
        })
      );
      
      setChatRooms(roomsWithCourseNames);
    } catch (err: any) {
      // Check if it's a 401 Unauthorized error
      if (err?.response?.status === 401 || err?.status === 401) {
        setError('Please login to view chat rooms');
      } else {
        setError(err instanceof Error ? err.message : 'An error occurred');
      }
      console.error('Error fetching chat rooms:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (authReady && auth.currentUser) {
      fetchChatRooms();
    }
  }, [authReady]);

  const handleRoomClick = (room: ChatRoom) => {
    setSelectedRoomId(room.ID);
    if (onRoomSelect) {
      onRoomSelect(room);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <p className="text-gray-500">Loading chat rooms...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <p className="text-red-500">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full overflow-y-auto">
      {chatRooms.length === 0 ? (
        <p className="text-gray-500 text-center p-4">No chat rooms found</p>
      ) : (
        chatRooms.map((room) => (
          <div 
            key={room.ID}
            className={selectedRoomId === room.ID ? 'bg-blue-50' : ''}
          >
            <LeftChat
              room={room}
              onClick={() => handleRoomClick(room)}
            />
          </div>
        ))
      )}
    </div>
  );
};

export default ChatRoomList;
